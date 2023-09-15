package main

import (
	"net/http"
	"os"

	"fn/internal/handlers"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Post("/api/bot", handlers.Interaction)
	r.Get("/api/registercommands", handlers.RegisterCommands)

	http.ListenAndServe(getPort(), r)
}

func getPort() string {
	port := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		port = ":" + val
	}
	return port
}
