package entity

type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"unique"`
    Password string
    Role     string // Admin, User, etc.
}
