package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KaePee/go-rssag/internal/auth"
	db "github.com/KaePee/go-rssag/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type parameters struct {
	Name string `json:"name"`
}

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding params: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
		CreatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't create user: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusForbidden, fmt.Sprintf("Authentication error: %v", err))
	}

	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not get user: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
