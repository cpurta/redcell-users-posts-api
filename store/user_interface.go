package store

import (
	"redcellpartners.com/users-posts-api/model"
)

type UserStore interface {
	ListUsers() ([]*model.User, error)
	CreateUser(*model.User) (*model.User, error)
	GetUser(id int) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id int) error
}
