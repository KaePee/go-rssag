package main

import (
	db "github.com/KaePee/go-rssag/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        pgtype.UUID      `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Name      string           `json:"name"`
	ApiKey    string           `json:"api_key"`
}

func databaseUserToUser(dbUser db.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}
