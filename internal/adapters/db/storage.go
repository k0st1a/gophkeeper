package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k0st1a/gophkeeper/internal/ports"
	"github.com/rs/zerolog/log"
)

type db struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, dsn string) (*db, error) {
	err := runMigrations(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}

	return &db{
		pool: pool,
	}, nil
}

func (d *db) Close() {
	d.pool.Close()
}

func (d *db) Ping(ctx context.Context) error {
	log.Ctx(ctx).Printf("Ping")

	err := d.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed ping:%w", err)
	}

	log.Ctx(ctx).Printf("Ping success")
	return nil
}

func (d *db) CreateUser(ctx context.Context, login, password string) (int64, error) {
	log.Ctx(ctx).Printf("CreateUser, Login:%s, password:%s", login, password)
	var id int64

	err := d.pool.QueryRow(ctx,
		"INSERT INTO users (login,password) VALUES($1,$2) "+
			"ON CONFLICT DO NOTHING "+
			"RETURNING id",
		login, password).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ports.ErrLoginAlreadyBusy
	}

	if err != nil {
		return id, fmt.Errorf("failed to create user:%w", err)
	}

	return id, nil
}

func (d *db) GetUserIDAndPassword(ctx context.Context, login string) (int64, string, error) {
	log.Ctx(ctx).Printf("GetUserIDAndPassword, Login:%s", login)
	var id int64
	var password string

	err := d.pool.QueryRow(ctx, "SELECT id, password FROM users WHERE login = $1", login).Scan(&id, &password)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, "", ports.ErrUserNotFound
	}

	if err != nil {
		return 0, "", fmt.Errorf("failed to get user id and password:%w", err)
	}

	return id, password, nil
}

func (d *db) CreateItem(ctx context.Context, userID int64, name, dataType string, data []byte) (int64, error) {
	log.Ctx(ctx).Printf("CreateItem, userID:%v, name:%v, dataType:%v", userID, name, dataType)
	var id int64

	err := d.pool.QueryRow(ctx,
		"INSERT INTO items (user_id, name, type, data, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
		userID, name, dataType, data, time.Now().UTC()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create item:%w", err)
	}

	log.Ctx(ctx).Printf("CreateItem success, id:%v", id)
	return id, nil
}

func (d *db) UpdateItem(ctx context.Context, userID, itemID int64, data []byte) error {
	log.Ctx(ctx).Printf("UpdateItem, userID:%v, itemID:%v", userID, itemID)
	var id int64

	err := d.pool.QueryRow(ctx,
		"UPDATE items SET data = $1, created_at = $2 WHERE id = $3, user_id = $4 RETURNING id",
		data, time.Now().UTC(), itemID, userID).Scan(&id)
	if err != nil {
		return fmt.Errorf("query error of update item:%w", err)
	}

	log.Ctx(ctx).Printf("UpdateItem success, id:%v", id)
	return nil
}

func (d *db) GetItem(ctx context.Context, userID, itemID int64) (*ports.ItemInfo, error) {
	log.Ctx(ctx).Printf("GetItem, userID:%v, itemID:%v", userID, itemID)
	var item ports.ItemInfo

	err := d.pool.QueryRow(ctx,
		"SELECT id, name, type, data, created_at FROM items WHERE user_id = $1, id = $2",
		userID, itemID).Scan(&item.ID, &item.Name, &item.Type, &item.Data, &item.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ports.ErrItemNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get item", err)
	}

	log.Ctx(ctx).Printf("GetItem success")
	return &item, nil
}

func (d *db) ListItem(ctx context.Context, userID int64) ([]ports.ItemInfo, error) {
	log.Ctx(ctx).Printf("ListItem, userID:%v", userID)
	var items []ports.ItemInfo

	rows, err := d.pool.Query(ctx,
		"SELECT id, name, type, data, created_at FROM items WHERE user_id = $1",
		userID)
	if err != nil {
		return items, fmt.Errorf("query error of list item:%w", err)
	}

	for rows.Next() {
		var item ports.ItemInfo
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Type,
			&item.Data,
			&item.CreatedAt,
		)
		if err != nil {
			return items, fmt.Errorf("scan error of list item:%w", err)
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return items, fmt.Errorf("error of list item", err)
	}

	log.Ctx(ctx).Printf("ListItem success")
	return items, nil
}
