package store

import (
	"redcellpartners.com/users-posts-api/model"
)

type PostStore interface {
	ListPosts() ([]*model.Post, error)
	CreatePost(*model.Post) (*model.Post, error)
	GetPost(id int) (*model.Post, error)
	UpdatePost(*model.Post) (*model.Post, error)
	DeletePost(id int) error
}
