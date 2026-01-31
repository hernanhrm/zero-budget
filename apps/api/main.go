package main

import (
	"fmt"
	"net/http"

	"backend/infra/logger"

	"feature/user"
)

func main() {
	log := logger.NewProduction()

	log.Info("starting api server")

	// In production, you would inject a real database pool here
	// db := database.New(ctx, "postgres://...")
	// For now, we'll skip database initialization

	// Example of how to wire up the user feature:
	// repo := user.NewRepository(db, log)
	// svc := user.NewService(repo, log)
	// handler := user.NewHandler(svc, log)

	// Register routes
	// mux := http.NewServeMux()
	// mux.HandleFunc("GET /users", handler.GetAll)
	// mux.HandleFunc("GET /users/{id}", handler.GetByID)
	// mux.HandleFunc("POST /users", handler.Create)
	// mux.HandleFunc("PUT /users/{id}", handler.Update)
	// mux.HandleFunc("DELETE /users/{id}", handler.Delete)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	log.Info("server listening", "port", 8080)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Error("server failed", "error", err)
	}

	// Demonstrate user types are accessible
	_ = user.User{}
	_ = user.CreateUser{}
}
