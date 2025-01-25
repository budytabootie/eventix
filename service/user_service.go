package service

import (
	"errors"
	"eventix/entity"
	"eventix/repository"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserService interface {
	RegisterUser(user entity.User) (entity.User, error)
	GetUserByID(id uint) (entity.User, error)
    UpdateUserRole(userID uint, role string) error
	Login(username, password string) (string, error)
}

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(user entity.User) (entity.User, error) {
    // Hash password sebelum disimpan
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return entity.User{}, err
    }
    user.Password = string(hashedPassword)

    return s.repo.CreateUser(user)
}

func (s *userService) GetUserByID(id uint) (entity.User, error) {
    return s.repo.GetUserByID(id)
}

func (s *userService) UpdateUserRole(userID uint, role string) error {
    // Validasi role (opsional, pastikan hanya role valid yang diterima)
    validRoles := []string{"Admin", "User"}
    isValid := false
    for _, r := range validRoles {
        if r == role {
            isValid = true
            break
        }
    }
    if !isValid {
        return errors.New("invalid role")
    }

    // Update role melalui repository
    return s.repo.UpdateUserRole(userID, role)
}


func (s *userService) Login(username, password string) (string, error) {
    user, err := s.repo.GetUserByUsername(username)
    if err != nil {
        return "", err
    }

    // Verifikasi password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate JWT
	claims := JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token berlaku 24 jam
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := "your_secret_key" // Ganti dengan key rahasia Anda
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
