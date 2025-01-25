package repository

import (
    "eventix/entity"
    "gorm.io/gorm"
)

type UserRepository interface {
    CreateUser(user entity.User) (entity.User, error)
    GetUserByID(id uint) (entity.User, error)
    GetUserByUsername(username string) (entity.User, error)
    UpdateUserRole(userID uint, role string) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user entity.User) (entity.User, error) {
    result := r.db.Create(&user)
    return user, result.Error
}

func (r *userRepository) GetUserByID(id uint) (entity.User, error) {
    var user entity.User
    result := r.db.First(&user, id)
    return user, result.Error
}

func (r *userRepository) GetUserByUsername(username string) (entity.User, error) {
    var user entity.User
    result := r.db.Where("username = ?", username).First(&user)
    return user, result.Error
}

func (r *userRepository) UpdateUserRole(userID uint, role string) error {
    result := r.db.Model(&entity.User{}).Where("id = ?", userID).Update("role", role)
    return result.Error
}

