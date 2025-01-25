package service

import (
    "eventix/repository"
    "time"
)

type AuthService interface {
    Logout(token string, expiresAt time.Time) error
}

type authService struct {
    blacklistRepo repository.TokenBlacklistRepository
}

func NewAuthService(blacklistRepo repository.TokenBlacklistRepository) AuthService {
    return &authService{blacklistRepo: blacklistRepo}
}

func (s *authService) Logout(token string, expiresAt time.Time) error {
    return s.blacklistRepo.AddToBlacklist(token, expiresAt)
}
