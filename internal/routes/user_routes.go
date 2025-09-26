package routes

import (
	"project/internal/handlers"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
	// User routes
	userRouter := router.PathPrefix("/users").Subrouter()
	
	userRouter.HandleFunc("", userHandler.GetUsers).Methods("GET")
	userRouter.HandleFunc("", userHandler.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")
	userRouter.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")
	
    // Authentication routes moved to auth routes file
}