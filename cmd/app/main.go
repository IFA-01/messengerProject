package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/IFA-01/messenger/internal/config"
	"github.com/IFA-01/messenger/internal/handlers"
	"github.com/IFA-01/messenger/internal/middleware"
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
		AllowedOrigins: []string{"https://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
	}))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Get("/err", handlers.HandleErr)

	//routes
	r.Route("/v1", func(r chi.Router) {
		r.Post("/users", handlers.HandlerCreateUser(q))
		r.Post("/login", handlers.HandleLogin(q, cfg.JWTSecret))

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleWare(cfg.JWTSecret))

			r.Get("/users/me", handlers.HandleGetUser(q))
			r.Post("/chats", handlers.HandleCreateChat(q))
		})
	})
	http.ListenAndServe(":"+cfg.ServerPort, r)
}
