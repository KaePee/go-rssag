package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	db "github.com/KaePee/go-rssag/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB *db.Queries
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load godotenv")
	}
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not set in .env")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL not set in .env")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatal("❌Cannot connect to database", err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	apiCfg := apiConfig{
		DB: queries,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleReady)
	v1Router.Get("/err", handleErr)
	v1Router.Post("/users", apiCfg.handleCreateUser)

	r.Mount("/v1", v1Router)

	fmt.Printf("starting http server on port: %v", portString)
	log.Fatal(http.ListenAndServe(":"+portString, r))
	//OR use the deference of http.Server
	// srv := &http.Server{
	// 	Addr:    ":" + portString,
	// 	Handler: r,
	// }
	// log.Printf("Listening to http on port: %v", portString)
	// srv.ListenAndServe()

}
