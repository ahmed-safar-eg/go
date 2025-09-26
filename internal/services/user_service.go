package services

import (
	"errors"
	"net/mail"
	"project/internal/models"
	"project/internal/repositories"
	"strconv"
	"strings"

	"project/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAllUsers() ([]models.UserResponse, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) GetUserByID(idStr string) (*models.UserResponse, error) {
	if idStr == "" {
		return nil, errors.New("user ID is required")
	}
	
	return s.userRepo.FindByID(idStr)
}

// sanitizeUserInputs trims whitespace and normalizes fields like email.
func sanitizeUserInputs(user *models.User) {
    user.Name = strings.TrimSpace(user.Name)
    user.Email = strings.TrimSpace(strings.ToLower(user.Email))
}

func isValidEmail(email string) bool {
    if email == "" {
        return false
    }
    _, err := mail.ParseAddress(email)
    return err == nil
}

func hashPassword(plain string) (string, error) {
    // Use bcrypt for new hashes
    hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return "bcrypt$" + string(hashed), nil
}

func verifyPassword(stored string, plain string) bool {
    if strings.HasPrefix(stored, "bcrypt$") {
        return bcrypt.CompareHashAndPassword([]byte(stored[len("bcrypt$"):]), []byte(plain)) == nil
    }
    return false
}

// Login moved to AuthService; intentionally removed from UserService.

func (s *UserService) CreateUser(user models.User) (*models.UserResponse, error) {
    // Normalize
    sanitizeUserInputs(&user)

    // التحقق من البيانات عبر validator
    if err := utils.ValidateStruct(user); err != nil {
        return nil, errors.New("invalid user payload")
    }
	
	// التحقق من عدم وجود email مكرر
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

    // Hash password before storing
    hashed, err := hashPassword(user.Password)
    if err != nil {
        return nil, errors.New("failed to hash password")
    }
    user.Password = hashed
	
	return s.userRepo.Create(user)
}

func (s *UserService) UpdateUser(idStr string, user models.User) (*models.UserResponse, error) {
	if idStr == "" {
		return nil, errors.New("user ID is required")
	}
	
	// التحقق من وجود المستخدم أولاً
	_, err := s.userRepo.FindByID(idStr)
	if err != nil {
		return nil, errors.New("user not found")
	}
	
    // Normalize inputs for update as well
    sanitizeUserInputs(&user)

	// إذا كان هناك email جديد، التحقق من عدم التكرار
	if user.Email != "" {
        if !isValidEmail(user.Email) {
            return nil, errors.New("invalid email format")
        }
        existingUser, err := s.userRepo.FindByEmail(user.Email)
		if err == nil && existingUser != nil {
			// تحقق أن Email الجديد لا ينتمي لمستخدم آخر
			currentID, _ := strconv.ParseUint(idStr, 10, 32)
			_ = currentID // for Mongo, repository will ensure uniqueness; keep placeholder
			return nil, errors.New("email already exists for another user")
		}
	}
	
    // If password provided, validate and hash
    if user.Password != "" {
        if len(user.Password) < 8 {
            return nil, errors.New("password must be at least 8 characters")
        }
        hashed, err := hashPassword(user.Password)
        if err != nil {
            return nil, errors.New("failed to hash password")
        }
        user.Password = hashed
    }

	return s.userRepo.Update(idStr, user)
}

func (s *UserService) DeleteUser(idStr string) error {
	if idStr == "" {
		return errors.New("user ID is required")
	}
	
	// التحقق من وجود المستخدم أولاً
	_, err := s.userRepo.FindByID(idStr)
	if err != nil {
		return errors.New("user not found")
	}
	
	return s.userRepo.Delete(idStr)
}