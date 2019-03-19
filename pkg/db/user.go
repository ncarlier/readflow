package db

import "github.com/ncarlier/reader/pkg/model"

// UserRepository is the repository interface to manage Users
type UserRepository interface {
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	CreateOrUpdateUser(user model.User) (*model.User, error)
	DeleteUser(user model.User) error
}
