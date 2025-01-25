package service

import (
	"errors"
	"eventix/repository"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, password string) (string, error)
	Logout(token string, expiresAt time.Time) error
}

type authService struct {
	userRepo      repository.UserRepository
	blacklistRepo repository.TokenBlacklistRepository
}

func NewAuthService(userRepo repository.UserRepository, blacklistRepo repository.TokenBlacklistRepository) AuthService {
	return &authService{
		userRepo:      userRepo,
		blacklistRepo: blacklistRepo,
	}
}

// Login implements the authentication logic
func (s *authService) Login(username, password string) (string, error) {
	// Cari user berdasarkan username
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := "your_secret_key" // Ganti dengan key rahasia Anda
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Logout implements the logic for blacklisting a token
func (s *authService) Logout(token string, expiresAt time.Time) error {
	return s.blacklistRepo.AddToBlacklist(token, expiresAt)
}
