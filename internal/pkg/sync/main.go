package sync

import (
	"context"
	"fmt"

	gclient "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/client"
	pclient "github.com/k0st1a/gophkeeper/internal/ports/client"

	"github.com/rs/zerolog/log"
)

type Doer interface {
	Do(ctx context.Context) error
}

type sync struct {
	local  pclient.ItemStorage
	remote gclient.ItemManager
}

// New - возвращает новый экземпляр синхронизации предметов между локальным и удаленным хранилищем.
func New(s pclient.ItemStorage, c gclient.ItemManager) *sync {
	return &sync{
		local:  s,
		remote: c,
	}
}

// Do - запуск синхронизации предметов между локальным и удаленным хранилищем.
func (s *sync) Do(ctx context.Context) error {
	log.Ctx(ctx).Printf("Do sync items")

	// local itmes
	litems, err := s.local.ListItems(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error of list local items")
		return fmt.Errorf("error of list local items:%w", err)
	}
	log.Printf("litems size:%v", len(litems))

	// local items to check for update
	citems := make([]*pclient.Item, 0, len(litems))

	// local items to upload later
	uitems := make([]*pclient.Item, 0)

	for _, i := range litems {
		if i.RemoteID == 0 && i.DeleteMark {
			log.Ctx(ctx).Printf("Delete local item(%v)", i.ID)
			err := s.local.DeleteItem(ctx, i.ID)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msgf("error of delete local item(%v)", i.ID)
			}
			log.Ctx(ctx).Printf("Local item(%v) deleted", i.ID)
			continue
		}

		if i.RemoteID == 0 {
			log.Ctx(ctx).Printf("Add item(%v) to uitems", i.ID)
			uitems = append(uitems, &i)
			continue
		}

		log.Ctx(ctx).Printf("Add item(%v) to citems", i.ID)
		citems = append(citems, &i)
	}
	log.Printf("uitems size:%v", len(uitems))
	log.Printf("citems size:%v", len(citems))

	// remote items
	ritems, err := s.remote.ListItems(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error of list remote items")
		return fmt.Errorf("error of list remote items:%w", err)
	}
	log.Printf("ritems size:%v", len(ritems))

	mcitems := pclient.List2MapWithRemoteID(citems)

	for _, ri := range ritems {
		log.Ctx(ctx).Printf("Remote item id(%v), find local item", ri.ID)
		li, ok := mcitems[ri.ID]
		if !ok {
			log.Ctx(ctx).Printf("Not found local item(%v) => need download item", ri.ID)
			err := s.downloadItem(ctx, &ri)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("")
			}
			continue
		}
		log.Ctx(ctx).Printf("Found local item(%v) => compare items", li.ID)

		cmp := compare(li, &ri)
		if cmp == 0 {
			log.Ctx(ctx).Printf("local item(%v) equal remote item(%v) => skip", li.ID, ri.ID)
			continue
		}

		if cmp == 1 {
			log.Printf("Need update local item(%v)", li.ID)
			err := s.updateLocalItem(ctx, li, &ri)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("")
			}
			continue
		}

		if cmp == 2 {
			log.Printf("Need update remote item(%v)", ri.ID)
			err := s.updateRemoteItem(ctx, &ri, li)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("")
			}
			continue
		}

		if cmp == 3 {
			log.Printf("Need delete remote item(%v) and remote item(%v)", ri.ID, li.ID)
			err := s.deleteBothItems(ctx, &ri, li)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("")
			}
			continue
		}

		log.Ctx(ctx).Error().Msgf("Compare return unknown code:%v", cmp)
	}

	s.uploadItems(ctx, uitems)

	return nil
}

// uploadItems - загрузить локальные предмет на удаленное хранилище.
func (s *sync) uploadItems(ctx context.Context, items []*pclient.Item) {
	log.Ctx(ctx).Printf("Upload items(%v)", len(items))

	for _, i := range items {
		err := s.uploadItem(ctx, i)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("")
		}
	}
}

// uploadItem - загрузить локальный предмет на удаленное хранилище.
func (s *sync) uploadItem(ctx context.Context, l *pclient.Item) error {
	log.Ctx(ctx).Printf("Upload local item(%v)", l.ID)

	r := makeRemoteItem(l)

	id, err := s.remote.CreateItem(ctx, r)
	if err != nil {
		return fmt.Errorf("error of upload local item(%v):%w", l.ID, err)
	}
	log.Ctx(ctx).Printf("Local item(%v) uploaded, remote id:%v", l.ID, id)

	ui := &pclient.UpdateItem{
		ID:       l.ID,
		RemoteID: &id,
	}

	err = s.local.UpdateItem(ctx, ui)
	if err != nil {
		return fmt.Errorf("error of update local item(%v):%w", l.ID, err)
	}
	log.Ctx(ctx).Printf("Local item(%v) updated", l.ID)

	return nil
}

// downloadItem - загрузить удаленный предмет в локальное хранилище.
func (s *sync) downloadItem(ctx context.Context, r *gclient.Item) error {
	log.Ctx(ctx).Printf("Download remote item(%v)", r.ID)

	l := makeLocalItem(r)

	id, err := s.local.CreateItem(ctx, l)
	if err != nil {
		return fmt.Errorf("error of create local item(remote id:%v):%w", r.ID, err)
	}
	log.Ctx(ctx).Printf("Remote item(%v) downloaded, local id:%v", r.ID, id)

	return nil
}

func (s *sync) updateLocalItem(ctx context.Context, l *pclient.Item, r *gclient.Item) error {
	log.Printf("Update local item(%v)", l.ID)

	ui := &pclient.UpdateItem{
		ID:         l.ID,
		RemoteID:   &r.ID,
		Body:       &r.Body,
		UpdateTime: &r.UpdateTime,
	}

	err := s.local.UpdateItem(ctx, ui)
	if err != nil {
		return fmt.Errorf("error of update local item(%v):%w", l.ID, err)
	}
	log.Ctx(ctx).Printf("Local item(%v) updated", l.ID)

	return nil
}

func (s *sync) updateRemoteItem(ctx context.Context, r *gclient.Item, l *pclient.Item) error {
	log.Printf("Update remote item(%v)", r.ID)

	r = makeRemoteItem(l)

	err := s.remote.UpdateItem(ctx, r)
	if err != nil {
		return fmt.Errorf("error of update remote item(%v)", r.ID)
	}
	log.Ctx(ctx).Printf("Remote item(%v) updated", r.ID)

	return nil
}

func (s *sync) deleteBothItems(ctx context.Context, r *gclient.Item, l *pclient.Item) error {
	log.Printf("Delete remote item(%v) and remote item(%v)", r.ID, l.ID)

	err := s.remote.DeleteItem(ctx, r.ID)
	if err != nil {
		return fmt.Errorf("error of delete remote item(%v) => skip delete local item(%v)", r.ID, l.ID)
	}
	log.Printf("Remote item(%v) deleted", r.ID)

	err = s.local.DeleteItem(ctx, l.ID)
	if err != nil {
		return fmt.Errorf("error of delete local(%v) => wiil be delete in next sync", l.ID)
	}
	log.Printf("Local item(%v) deleted", l.ID)

	return nil
}

// compare - сравить локальный и удаленный предметы.
//
//		Возвращает: 
//			0, если предметы одинаковы
//	        1, если нужно обновить локальный предмет
//		    2, если нужно обновить удаленный предмет
//	        3, если нужно удалить  локальный и удаленный предметы
func compare(l *pclient.Item, r *gclient.Item) int {
	log.Printf("Compare items, l.UpdateTime:%v, r.UpdateTime:%v", l.UpdateTime, r.UpdateTime)
	ld := l.UpdateTime.UnixMilli()
	rd := r.UpdateTime.UnixMilli()
	if ld > rd {
		log.Printf("l.DeleteMark:%v", l.DeleteMark)
		if l.DeleteMark {
			return 3
		}
		return 2
	}

	if ld < rd {
		return 1
	}

	return 0
}

// makeLocalItem - создать локальный прдемет на основе удаленного.
func makeLocalItem(r *gclient.Item) *pclient.Item {
	return &pclient.Item{
		RemoteID:   r.ID,
		Body:       r.Body,
		CreateTime: r.CreateTime,
		UpdateTime: r.UpdateTime,
		DeleteMark: false,
	}
}

// makeRemoteItem - создать удаленный прдемет на основе локального.
func makeRemoteItem(l *pclient.Item) *gclient.Item {
	return &gclient.Item{
		ID:         l.RemoteID,
		Body:       l.Body,
		CreateTime: l.CreateTime,
		UpdateTime: l.UpdateTime,
	}
}
