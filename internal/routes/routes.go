package routes

import (
	"net/http"
	"project/internal/handlers"
	"project/internal/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	// Initialize handlers
	userHandler := handlers.NewUserHandler()
    authHandler := handlers.NewAuthHandler()
	// productHandler := handlers.NewProductHandler(db)
	
    // Global middlewares (JSON only here); Auth applied on protected subrouter below
    router.Use(middleware.JSONMiddleware)

    // Register public auth routes BEFORE applying auth to protected subrouter
    RegisterAuthRoutes(router, authHandler)

    // Protected API subrouter with Auth middleware
    protected := router.PathPrefix("").Subrouter()
    // protected.Use(middleware.Auth)

    // Register all protected routes
    RegisterUserRoutes(protected, userHandler)
	// RegisterProductRoutes(router, productHandler)
	
    // Health check (public route, JSON)
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        jsonBody := []byte(`{"success":true,"message":"OK"}`)
        w.Write(jsonBody)
    }).Methods("GET")

    // 404/405 JSON responses
    router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"success":false,"message":"route not found"}`))
    })
    router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte(`{"success":false,"message":"method not allowed"}`))
    })
}

func RegisterAPIRoutes(router *mux.Router) {
	// API version 1 routes
	apiV1 := router.PathPrefix("/").Subrouter()
	// apiV1.Use(middleware.Auth) // Apply auth middleware to all API routes
	
	RegisterRoutes(apiV1)
}