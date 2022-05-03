package store

import (
	"fmt"
	"go-proto-crud/post"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Store interface {
	List() []post.Post
	Read(id string) (*post.Post, error)
	Create(post post.Post) error
	Update(post post.Post) error
	Remove(id string) error
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
func (s *store) Update(post post.Post) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range s.posts {
		if s.posts[i].Id == post.Id {
			s.posts[i].Title = post.Title
			s.posts[i].Rate = post.Rate
			s.posts[i].IsDone = post.IsDone
			s.posts[i].Updated = timestamppb.Now()
			return nil
		}
	}
	return fmt.Errorf("post not found")
}

func (s *store) Remove(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range s.posts {
		if s.posts[i].Id == id {
			s.posts = append(s.posts[:i], s.posts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("post not found")
}

func NewStore() *store {
	return &store{
		posts: []post.Post{
			{
				Id:      uuid.New().String(),
				Title:   "new title",
				Rate:    4,
				IsDone:  false,
				Updated: nil,
			},
			{
				Id:      uuid.New().String(),
				Title:   "new title 2",
				Rate:    3,
				IsDone:  true,
				Updated: timestamppb.Now(),
			},
		},
	}
}
