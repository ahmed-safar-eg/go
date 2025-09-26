package main

import (
	"log"
	"net/http"
	"project/internal/config"
	"project/internal/database"

	// "project/internal/handlers"
	// "project/internal/middleware"
	"project/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
    // Initialize MongoDB
    if err := database.InitializeMongo(cfg.Mongo); err != nil {
        log.Fatal("Mongo connection failed: ", err)
    }
    if err := database.EnsureIndexes(); err != nil {
        log.Fatal("Mongo index ensure failed: ", err)
    }
    defer database.CloseMongo()
	// Create router
	router := mux.NewRouter()
	routes.RegisterAPIRoutes(router) 
	// Global middleware
	// router.Use(middleware.Logging)
	// router.Use(middleware.CORS)
	
	// API routes
	// api := router.PathPrefix("/api/v1").Subrouter()
	// api.Use(middleware.Auth)
	
	// Register handlers
	// userHandler := handlers.NewUserHandler(db)
	// productHandler := handlers.NewProductHandler(db)
	
	// User routes
	// api.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	// api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	// api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	
	// Product routes
	// api.HandleFunc("/products", productHandler.GetProducts).Methods("GET")
	// api.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	
	// Health check
	// router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("OK"))
	// }).Methods("GET")
	
	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))
}