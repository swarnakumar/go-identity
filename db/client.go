package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/swarnakumar/go-identity/db/conn"
	"github.com/swarnakumar/go-identity/db/users"
)

type Client struct {
	db    *pgxpool.Pool
	Users *users.Users
}

func (client *Client) Close() {
	client.db.Close()
}

func New(ctx context.Context) *Client {
	// setup database
	db := conn.NewPool(ctx)

	client := Client{db: db, Users: users.New(db)}
	return &client
}
