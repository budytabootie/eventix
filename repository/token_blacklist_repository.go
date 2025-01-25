package repository

import (
	"eventix/entity"
	"time"

	"gorm.io/gorm"
)

type TokenBlacklistRepository interface {
    AddToBlacklist(token string, expiresAt time.Time) error
    IsTokenBlacklisted(token string) (bool, error)
}

type tokenBlacklistRepository struct {
    db *gorm.DB
}

func NewTokenBlacklistRepository(db *gorm.DB) TokenBlacklistRepository {
    return &tokenBlacklistRepository{db: db}
}

func (r *tokenBlacklistRepository) AddToBlacklist(token string, expiresAt time.Time) error {
    blacklist := entity.TokenBlacklist{
        Token:     token,
        ExpiresAt: expiresAt,
    }
    return r.db.Create(&blacklist).Error
}

func (r *tokenBlacklistRepository) IsTokenBlacklisted(token string) (bool, error) {
    var count int64
    err := r.db.Model(&entity.TokenBlacklist{}).Where("token = ?", token).Count(&count).Error
    return count > 0, err
}
