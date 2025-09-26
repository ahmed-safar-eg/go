package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/internal/models"
	"project/internal/repositories"
	"project/internal/services"
	"project/pkg/utils"
	"strings"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler() *UserHandler {
    // Switch to Mongo repository
    userRepo := repositories.NewUserRepositoryMongo()
	userService := services.NewUserService(userRepo)
	return &UserHandler{userService: userService}
}

// دالة مساعدة لإرجاع errors كـ JSON
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Success: false,
		Message: message,
	})
}

// دالة مساعدة لإرجاع success كـ JSON
func sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := models.SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get users: "+err.Error())
		return
	}
	
	sendSuccessResponse(w, http.StatusOK, "Users retrieved successfully", users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	user, err := h.userService.GetUserByID(idStr)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "User not found: "+err.Error())
		return
	}
	
	sendSuccessResponse(w, http.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	
	// التحقق من Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		sendErrorResponse(w, http.StatusBadRequest, "Content-Type must be application/json")
		return
	}
	
	// تحليل JSON body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}
	
	// ✅ التصحيح: التحقق من البيانات بشكل صحيح
	if strings.TrimSpace(user.Name) == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Name is required")
		return
	}
	
	if strings.TrimSpace(user.Email) == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Email is required")
		return
	}
	
	// ✅ التصحيح المهم: التحقق من أن Password ليست string فارغة
	// عند decoding JSON، إذا كانت password غير موجودة أو null، ستصبح string فارغة
	if user.Password == "" {
		fmt.Println(user.Password)
		sendErrorResponse(w, http.StatusBadRequest, "Password is requireds")
		return
	}
	
	// التحقق من صحة Email (تحقق بسيط)
	if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid email format")
		return
	}
	
    // detailed validation
    if errs, err := utils.ValidateStructDetailed(user); err != nil {
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

    createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		// إذا كان الخطأ بسبب email مكرر
		if strings.Contains(err.Error(), "already exists") {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		}
		return
	}
	
	sendSuccessResponse(w, http.StatusCreated, "User created successfully", createdUser)
}

// Login removed from User handler; handled by AuthHandler.

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}
	
	updatedUser, err := h.userService.UpdateUser(idStr, user)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendErrorResponse(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), "already exists") {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		}
		return
	}
	
	sendSuccessResponse(w, http.StatusOK, "User updated successfully", updatedUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	err := h.userService.DeleteUser(idStr)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sendErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		}
		return
	}
	
	sendSuccessResponse(w, http.StatusOK, "User deleted successfully", nil)
}