package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	db "github.com/KaePee/go-rssag/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (apiCfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding params: %v", err))
		return
	}
	if params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "name required")
		return
	}
	if params.Url == "" {
		respondWithError(w, http.StatusBadRequest, "url required")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), db.CreateFeedParams{
		ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
		CreatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		Name:      params.Name,
		Url:       pgtype.Text{String: params.Url, Valid: true},
		UserID:    user.ID,
	})

	respondWithJSON(w, http.StatusCreated, databaseFeedtoFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error getting feeds: %v", err))
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedstoFeeds(feeds))
}
