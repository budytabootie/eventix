package entity

import "time"

type Ticket struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventID   uint      `json:"event_id"`
	UserID    uint      `json:"user_id"`
	Quantity  int       `json:"quantity"` // Tambahkan field Quantity
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
