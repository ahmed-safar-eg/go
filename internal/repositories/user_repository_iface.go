package repositories

import "project/internal/models"

type UserRepositoryInterface interface {
    FindAll() ([]models.UserResponse, error)
    FindByID(idStr string) (*models.UserResponse, error)
    FindByEmail(email string) (*models.User, error)
    Create(user models.User) (*models.UserResponse, error)
    Update(idStr string, user models.User) (*models.UserResponse, error)
    Delete(idStr string) error
}


