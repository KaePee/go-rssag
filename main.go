package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func welcome() {
	fmt.Println("Welcome")
}
func main() {
	welcome()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load godotenv")
	}
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not set in .env")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleReady)
	v1Router.Get("/err", handleErr)

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
