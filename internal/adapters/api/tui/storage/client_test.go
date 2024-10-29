package storage

import (
	"context"
	"testing"
	"time"

	"github.com/k0st1a/gophkeeper/internal/adapters/storage/inmemory"
	"github.com/k0st1a/gophkeeper/internal/pkg/client/model"
	pclient "github.com/k0st1a/gophkeeper/internal/ports/client"
	mockinmemory "github.com/k0st1a/gophkeeper/mock/storage/inmemory"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name    string
		storage *inmemory.Storage
		body    any
		meta    Meta
	}{
		{
			name:    "Check CreateItem Password",
			storage: inmemory.New(),
			body: &Password{
				Resource: "Resource",
				UserName: "Username",
				Password: "Password",
			},
			meta: map[string]string{
				model.MetaKeyDescription:           "Password description",
				model.MetaKeyAdditionalInformation: "Password additional information",
			},
		},
		{
			name:    "Check CreateItem Card",
			storage: inmemory.New(),
			body: &Card{
				Number:  "Number",
				Expires: "Expires",
				Holder:  "Holder",
			},
			meta: map[string]string{
				model.MetaKeyDescription:           "Card description",
				model.MetaKeyAdditionalInformation: "Card additional information",
			},
		},
		{
			name:    "Check CreateItem Note",
			storage: inmemory.New(),
			body: &Note{
				Name: "Name",
				Body: "Body",
			},
			meta: map[string]string{
				model.MetaKeyDescription:           "Note description",
				model.MetaKeyAdditionalInformation: "Note additional information",
			},
		},
		{
			name:    "Check CreateItem File",
			storage: inmemory.New(),
			body: &File{
				Name: "Name",
				Body: []byte("body"),
			},
			meta: map[string]string{
				model.MetaKeyDescription:           "File description",
				model.MetaKeyAdditionalInformation: "File additional information",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := New(test.storage)
			ctx := context.Background()
			id, err := c.CreateItem(ctx, test.body, test.meta)
			require.NoError(t, err)
			item, err := c.GetItem(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id, item.ID)
			require.Equal(t, test.body, item.Body)
			require.Equal(t, test.meta, item.Meta)
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name    string
		storage *inmemory.Storage
		body    any
		meta    map[string]string
	}{
		{
			name:    "Check DeleteItem Password",
			storage: inmemory.New(),
			body: &Password{
				Resource: "Resource",
				UserName: "Username",
				Password: "Password",
			},
		},
		{
			name:    "Check DeleteItem Card",
			storage: inmemory.New(),
			body: &Card{
				Number:  "Number",
				Expires: "Expires",
				Holder:  "Holder",
			},
		},
		{
			name:    "Check DeleteItem Note",
			storage: inmemory.New(),
			body: &Note{
				Name: "Name",
				Body: "Body",
			},
		},
		{
			name:    "Check DeleteItem File",
			storage: inmemory.New(),
			body: &File{
				Name: "Name",
				Body: []byte("body"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := New(test.storage)
			ctx := context.Background()
			id, err := c.CreateItem(ctx, test.body, test.meta)
			require.NoError(t, err)
			err = c.DeleteItem(ctx, id)
			require.NoError(t, err)
			_, err = c.GetItem(ctx, id)
			require.Error(t, err)
		})
	}
}

func TestUpdateItem(t *testing.T) {
	tests := []struct {
		name         string
		storage      *inmemory.Storage
		createBodies []any
		createMetas  []Meta
		updateBodies []any
		updateMetas  []Meta
	}{
		{
			name:    "Check UpdateItem",
			storage: inmemory.New(),
			createBodies: []any{
				&Password{
					Resource: "Resource",
					UserName: "Username",
					Password: "Password",
				},
				&Card{
					Number:  "Number",
					Expires: "Expires",
					Holder:  "Holder",
				},
				&Note{
					Name: "Name",
					Body: "Body",
				},
				&File{
					Name: "Name",
					Body: []byte("body"),
				},
			},
			createMetas: []Meta{
				Meta{
					model.MetaKeyDescription:           "Password description",
					model.MetaKeyAdditionalInformation: "Password additional information",
				},
				Meta{
					model.MetaKeyDescription:           "Card description",
					model.MetaKeyAdditionalInformation: "Card additional information",
				},
				Meta{
					model.MetaKeyDescription:           "Note description",
					model.MetaKeyAdditionalInformation: "Note additional information",
				},
				Meta{
					model.MetaKeyDescription:           "File description",
					model.MetaKeyAdditionalInformation: "File additional information",
				},
			},
			updateBodies: []any{
				&Password{
					Resource: "Resource updated",
					UserName: "Username updated",
					Password: "Password updated",
				},
				&Card{
					Number:  "Number updated",
					Expires: "Expires updated",
					Holder:  "Holder updated",
				},
				&Note{
					Name: "Name updated",
					Body: "Body updated",
				},
				&File{
					Name: "Name updated",
					Body: []byte("body updated"),
				},
			},
			updateMetas: []Meta{
				Meta{
					model.MetaKeyDescription:           "Updated password description",
					model.MetaKeyAdditionalInformation: "Updated password additional information",
				},
				Meta{
					model.MetaKeyDescription:           "Updated card description",
					model.MetaKeyAdditionalInformation: "Updated card additional information",
				},
				Meta{
					model.MetaKeyDescription:           "Updated note description",
					model.MetaKeyAdditionalInformation: "Updated note additional information",
				},
				Meta{
					model.MetaKeyDescription:           "Updated file description",
					model.MetaKeyAdditionalInformation: "Updated file additional information",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := New(test.storage)
			ctx := context.Background()
			ids := make([]string, len(test.createBodies))
			for p, b := range test.createBodies {
				id, err := c.CreateItem(ctx, b, test.createMetas[p])
				require.NoError(t, err)
				ids[p] = id
			}
			createdItems := make([]*Item, len(test.createBodies))
			for p, id := range ids {
				i, err := c.GetItem(ctx, id)
				require.NoError(t, err)
				createdItems[p] = i
			}
			for p, b := range test.updateBodies {
				createdItems[p].Body = b
				createdItems[p].Meta = test.updateMetas[p]
				err := c.UpdateItem(ctx, createdItems[p])
				require.NoError(t, err)
			}
			for _, ci := range createdItems {
				i, err := c.GetItem(ctx, ci.ID)
				require.NoError(t, err)
				require.Equal(t, ci.Body, i.Body)
				require.Equal(t, ci.Meta, i.Meta)
				require.Equal(t, ci.CreateTime, i.CreateTime)
			}
		})
	}
}

func TestListItems(t *testing.T) {
	tests := []struct {
		name          string
		inmemoryItems []pclient.Item
		tuiItems      []Item
	}{
		{
			name: "Check ListItems",
			inmemoryItems: []pclient.Item{
				pclient.Item{
					Body: model.Item{
						Card: &model.Card{
							Number:  "Number",
							Expires: "Expires",
							Holder:  "Holder",
						},
						Meta: model.Meta{
							"description":            "card description",
							"additional information": "card additional information",
						},
					},
					CreateTime: time.Date(2024, time.May, 5, 8, 10, 0, 0, time.UTC),
					UpdateTime: time.Date(2024, time.May, 5, 8, 10, 0, 0, time.UTC),
					ID:         "ID",
					RemoteID:   1,
					DeleteMark: false,
				},
				pclient.Item{
					Body: model.Item{
						Card: &model.Card{
							Number:  "Number 2",
							Expires: "Expires 2",
							Holder:  "Holder 2",
						},
						Meta: model.Meta{
							"description":            "mark deleted card description",
							"additional information": "mark deleted card additional information",
						},
					},
					CreateTime: time.Date(2024, time.May, 5, 8, 10, 0, 2, time.UTC),
					UpdateTime: time.Date(2024, time.May, 5, 8, 10, 0, 2, time.UTC),
					ID:         "ID 2",
					RemoteID:   2,
					DeleteMark: true,
				},
			},
			tuiItems: []Item{
				Item{
					CreateTime: time.Date(2024, time.May, 5, 8, 10, 0, 0, time.UTC),
					UpdateTime: time.Date(2024, time.May, 5, 8, 10, 0, 0, time.UTC),
					Body: &Card{
						Number:  "Number",
						Expires: "Expires",
						Holder:  "Holder",
					},
					Meta: Meta{
						"description":            "card description",
						"additional information": "card additional information",
					},
					ID: "ID",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()
			ms := mockinmemory.NewMockItemStorage(ctrl)
			ms.
				EXPECT().
				ListItems(ctx).
				Return(test.inmemoryItems, nil)

			c := New(ms)
			items, err := c.ListItems(ctx)
			require.NoError(t, err)
			require.ElementsMatch(t, test.tuiItems, items)
		})
	}
}

func TestClear(t *testing.T) {
	tests := []struct {
		name    string
		storage *inmemory.Storage
		bodies  []any
		metas   []Meta
	}{
		{
			name:    "Check Clear storage",
			storage: inmemory.New(),
			bodies: []any{
				&Password{
					Resource: "Resource",
					UserName: "Username",
					Password: "Password",
				},
				&Card{
					Number:  "Number",
					Expires: "Expires",
					Holder:  "Holder",
				},
				&Note{
					Name: "Name",
					Body: "Body",
				},
				&File{
					Name: "Name",
					Body: []byte("body"),
				},
			},
			metas: []Meta{
				Meta{
					model.MetaKeyDescription:           "Password description",
					model.MetaKeyAdditionalInformation: "Password additional information",
				},
				Meta{
					model.MetaKeyDescription:           "Card description",
					model.MetaKeyAdditionalInformation: "Card additional information",
				},
				Meta{
					model.MetaKeyDescription:           "Note description",
					model.MetaKeyAdditionalInformation: "Note additional information",
				},
				Meta{
					model.MetaKeyDescription:           "File description",
					model.MetaKeyAdditionalInformation: "File additional information",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := New(test.storage)
			ctx := context.Background()
			for p, b := range test.bodies {
				_, err := c.CreateItem(ctx, b, test.metas[p])
				require.NoError(t, err)
			}
			items, err := c.ListItems(ctx)
			require.NoError(t, err)
			require.Len(t, items, len(test.bodies))
			c.Clear(ctx)
			items, err = c.ListItems(ctx)
			require.NoError(t, err)
			require.Len(t, items, 0)
		})
	}
}
