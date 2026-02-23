package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	db "github.com/KaePee/go-rssag/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
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
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal("❌Cannot connect to database", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	apiCfg := apiConfig{
		DB: queries,
	}

	//go routine for scraper
	go startScraper(
		queries, 10, time.Minute,
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	v1Router := chi.NewRouter()
	v1Router.Get("/health", handleReady)
	v1Router.Get("/err", handleErr)

	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handleGetPostsForUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollow))

	r.Mount("/v1", v1Router)

	fmt.Printf("starting http server on port: %v", portString)
	err = http.ListenAndServe(":"+portString, r)
	if err != nil {
		log.Fatalf("Error starting http server: %v", err)
	}
	//OR use the deference of http.Server
	// srv := &http.Server{
	// 	Addr:    ":" + portString,
	// 	Handler: r,
	// }
	// log.Printf("Listening to http on port: %v", portString)
	// srv.ListenAndServe()

}
