package entity

import "time"

type Event struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);index" json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Capacity    int       `json:"capacity"`
	Price       float64   `json:"price"` // Harga tiket untuk event
	Status      string    `gorm:"type:varchar(50);default:'active'" json:"status"` // Status event: active, ongoing, completed
	ImageURL    string    `gorm:"type:varchar(255)" json:"image_url"` // URL gambar untuk event
	CreatedAt   time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
