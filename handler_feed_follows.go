package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/KaePee/go-rssag/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (apiCfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding params: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), db.CreateFeedFollowParams{
		ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
		CreatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		UserID:    user.ID,
		FeedID:    pgtype.UUID{Bytes: params.FeedId, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create feed follow: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user db.User) {

	feedFollows, err := apiCfg.DB.GetFeedFollow(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't get feed follows: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedsFollowstoFeedsFollows(feedFollows))
}

func (apiCfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user db.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Cannot parse UUID: %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), db.DeleteFeedFollowParams{
		ID: pgtype.UUID{Bytes: feedFollowID, Valid: true},
		UserID: user.ID,
	})
	log.Printf("details to delete; feedID: %v, userID: %v", feedFollowID, user.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to delete feed follow. ID: %s error: %v", feedFollowID, err))
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}