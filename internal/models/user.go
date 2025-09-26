package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name      string             `json:"name" validate:"required,min=2" bson:"name"`
    Email     string             `json:"email" validate:"required,email" bson:"email"`
    Password  string             `json:"password" validate:"required,min=8" bson:"password"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
    DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type UserResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
// ErrorResponse هيكل لردود الأخطاء بشكل JSON
type ErrorResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}

// SuccessResponse هيكل لردود النجاح
type SuccessResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}
type LoginResponse struct {
    Token string `json:"token"`
    User  *UserResponse `json:"user"`
}
func (u *User) ToResponse() *UserResponse {
    return &UserResponse{
        ID:        u.ID.Hex(),
        Name:      u.Name,
        Email:     u.Email,
        CreatedAt: u.CreatedAt,
    }
}