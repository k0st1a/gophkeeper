package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/k0st1a/gophkeeper/internal/adapters/storage/inmemory"
)

func TestCreateItem(t *testing.T) {
	//nolint:dupl // not need here
	tests := []struct {
		name    string
		storage *inmemory.Storage
		body    any
	}{
		{
			name:    "Check CreateItem Password",
			storage: inmemory.New(),
			body: &Password{
				Resource: "Resource",
				UserName: "Username",
				Password: "Password",
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
		},
		{
			name:    "Check CreateItem Note",
			storage: inmemory.New(),
			body: &Note{
				Name: "Name",
				Body: "Body",
			},
		},
		{
			name:    "Check CreateItem File",
			storage: inmemory.New(),
			body: &File{
				Name:        "Name",
				Description: "Description",
				Body:        []byte("body"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := New(test.storage)
			ctx := context.Background()
			id, err := c.CreateItem(ctx, test.body)
			require.NoError(t, err)
			item, err := c.GetItem(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id, item.ID)
			require.Equal(t, test.body, item.Body)
		})
	}
}

func TestDeleteItem(t *testing.T) {
	//nolint:dupl // not need here
	tests := []struct {
		name    string
		storage *inmemory.Storage
		body    any
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
				Name:        "Name",
				Description: "Description",
				Body:        []byte("body"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := New(test.storage)
			ctx := context.Background()
			id, err := c.CreateItem(ctx, test.body)
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
		updateBodies []any
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
					Name:        "Name",
					Description: "Description",
					Body:        []byte("body"),
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
					Name:        "Name updated",
					Description: "Description updated",
					Body:        []byte("body updated"),
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
				id, err := c.CreateItem(ctx, b)
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
				err := c.UpdateItem(ctx, createdItems[p])
				require.NoError(t, err)
			}
			for _, ci := range createdItems {
				i, err := c.GetItem(ctx, ci.ID)
				require.NoError(t, err)
				require.Equal(t, ci.Body, i.Body)
				require.Equal(t, ci.CreateTime, i.CreateTime)
			}
		})
	}
}

func TestListItems(t *testing.T) {
	tests := []struct {
		name    string
		storage *inmemory.Storage
		bodies  []any
	}{
		{
			name:    "Check ListItems",
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
					Name:        "Name",
					Description: "Description",
					Body:        []byte("body"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := New(test.storage)
			ctx := context.Background()
			for _, b := range test.bodies {
				_, err := c.CreateItem(ctx, b)
				require.NoError(t, err)
			}
			items, err := c.ListItems(ctx)
			require.NoError(t, err)
			require.Len(t, items, len(test.bodies))
			bodies := make([]any, len(test.bodies))
			for p, i := range items {
				bodies[p] = i.Body
			}
			require.Equal(t, len(test.bodies), len(bodies))
			require.ElementsMatch(t, test.bodies, bodies)
		})
	}
}

func TestClear(t *testing.T) {
	tests := []struct {
		name    string
		storage *inmemory.Storage
		bodies  []any
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
					Name:        "Name",
					Description: "Description",
					Body:        []byte("body"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := New(test.storage)
			ctx := context.Background()
			for _, b := range test.bodies {
				_, err := c.CreateItem(ctx, b)
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
