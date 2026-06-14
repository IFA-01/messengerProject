package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/IFA-01/messenger/internal/config" // ← добавь
	"github.com/IFA-01/messenger/internal/handlers"
	"github.com/IFA-01/messenger/internal/repository/queries"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}
	defer dbpool.Close()

	q := queries.New(dbpool)

	fmt.Println("Server is Up:")
	fmt.Printf("Server is on port: %s\n", cfg.ServerPort)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
	}))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	//routes
	r.Route("/v1", func(r chi.Router) {
		r.Post("/users", handlers.HandlerCreateUser(q))
	})

	http.ListenAndServe(":"+cfg.ServerPort, r)

}
