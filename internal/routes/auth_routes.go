package routes

import (
	"project/internal/handlers"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler) {
    // Public auth routes - no Auth middleware
    authRouter := router.PathPrefix("/auth").Subrouter()
    authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
}


