package repositories

import (
	"golang/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUsers() ([]models.User, error)
	GetUser(ID int) (models.User, error)
}

type repository struct {
	db *gorm.DB
}

func RepositoryUser(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindUsers() ([]models.User, error) {
	var users []models.User
	// Using "Find" method here ...

	return users, err
}

func (r *repository) GetUser(ID int) (models.User, error) {
	var user models.User
	// Using "First" method here ...

	return user, err
}
