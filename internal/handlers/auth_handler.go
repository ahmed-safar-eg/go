package handlers

import (
	"encoding/json"
	"net/http"
	"project/internal/models"
	"project/internal/repositories"
	"project/internal/services"
	"project/pkg/utils"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
    userRepo := repositories.NewUserRepositoryMongo()
    authService := services.NewAuthService(userRepo)
    return &AuthHandler{authService: authService}
}

type loginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req loginRequest
    if r.Header.Get("Content-Type") != "application/json" {
        sendErrorResponse(w, http.StatusBadRequest, "Content-Type must be application/json")
        return
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
        return
    }
    if errs, err := utils.ValidateStructDetailed(req); err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Validation error: "+err.Error())
        return
    } else if len(errs) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": false,
            "message": "Validation failed",
            "errors":  errs,
        })
        return
    }
    token, user, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        sendErrorResponse(w, http.StatusUnauthorized, err.Error())
        return
    }
    sendSuccessResponse(w, http.StatusOK, "Login successful", models.LoginResponse{Token: token, User: user})
    // return
}


