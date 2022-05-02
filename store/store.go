package store

import (
	"fmt"
	"go-proto-crud/post"
	"sync"
)

type Store interface {
	List() []post.Post
	Create(post post.Post) error
	Read(id string) (*post.Post, error)
}

type store struct {
	mutex sync.RWMutex
	posts []post.Post
}

func (s *store) List() []post.Post {
	return s.posts
}

func (s *store) Read(id string) (*post.Post, error) {
	for i := range s.posts {
		if s.posts[i].Id == id {
			return &s.posts[i], nil
		}
	}
	return nil, fmt.Errorf("post not found")
}

func (s *store) Create(post post.Post) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.posts = append(s.posts, post)
	return nil
}

func NewStore() *store {
	return &store{
		posts: []post.Post{},
	}
}
