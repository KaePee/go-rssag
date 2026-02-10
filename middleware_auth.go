package main

import (
	"fmt"
	"net/http"

	"github.com/KaePee/go-rssag/internal/auth"
	db "github.com/KaePee/go-rssag/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, db.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("Authentication error: %v", err))
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
