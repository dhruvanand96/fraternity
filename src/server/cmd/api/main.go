package main

import (
	"log"
	"net/http"

	"app/server/internal/auth"
	"app/server/internal/config"
	"app/server/internal/httpapi"
	"app/server/internal/user"
)

func main() {
	cfg := config.Load()
	store, err := user.NewFileStore(cfg.DataFile)
	if err != nil {
		log.Fatalf("initialize user store: %v", err)
	}

	handler := httpapi.NewHandler(store, auth.NewService(cfg.JWTSecret), cfg.AllowedOrigin)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	log.Printf("fraternity API listening on http://localhost:%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("serve API: %v", err)
	}
}
