package entity

import "time"

type TokenBlacklist struct {
    ID        uint      `gorm:"primaryKey"`
    Token     string    `gorm:"type:text;not null"`
    ExpiresAt time.Time `gorm:"not null"`
}
