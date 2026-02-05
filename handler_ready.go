package main

import "net/http"

func handleReady(w http.ResponseWriter, r *http.Request) {
	status := struct{
		Status string `json:"status"`
	} {
		Status: "Are you feeling good?",
	}
	respondWithJSON(w, http.StatusOK, status)
}
