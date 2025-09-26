package services

import (
	"errors"
	"net/mail"
	"project/internal/config"
	"project/internal/models"
	"project/internal/repositories"
	"project/pkg/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo repositories.UserRepositoryInterface
}

func NewAuthService(userRepo repositories.UserRepositoryInterface) *AuthService {
    return &AuthService{userRepo: userRepo}
}

// verifyPassword must match the hashing used in UserService
func verifyPasswordAuth(stored string, plain string) bool {
    if strings.HasPrefix(stored, "bcrypt$") {
        return bcrypt.CompareHashAndPassword([]byte(stored[len("bcrypt$"):]), []byte(plain)) == nil
    }
    return false
}

func (s *AuthService) Login(email, password string) (string,*models.UserResponse, error) {
    email = strings.TrimSpace(strings.ToLower(email))
    if email == "" || password == "" {
        return "",nil, errors.New("invalid email or password")
    }
    if _, err := mail.ParseAddress(email); err != nil {
        return "",nil, errors.New("invalid email or password")
    }
    user, err := s.userRepo.FindByEmail(email)
    if err != nil || user == nil {
        return "",nil, errors.New("invalid email or password")
    }
    if !verifyPasswordAuth(user.Password, password) {
        return "",nil, errors.New("invalid email or password")
    }
    cfg := config.LoadConfig()
    token, err := utils.GenerateJWT(user.ID.Hex(), cfg.JWT.Secret, cfg.JWT.Expiry)
    if err != nil {
        return "",nil, errors.New("failed to generate token")
    }
    return  token, user.ToResponse(), nil
}


