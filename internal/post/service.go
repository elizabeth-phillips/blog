package post

import (
	"time"

	"github.com/google/uuid"
)

//Service creates the business logic that is a concrete implementation of how a service will interact with the Repository methods
type Service interface {
	CreatePost(post *Post) error
	FindPostByID(id string) (*Post, error)
	FindAllPosts() ([]*Post, error)
	UpdatePost(post *Post) error
	DeletePost(id string) error
}

type postService struct {
	repo Repository
}

//NewService is a constructor for Service
func NewService(repo Repository) Service {
	return &postService{
		repo,
	}
}

func (s *postService) CreatePost(post *Post) error {
	post.ID = uuid.New().String()
	post.Created = time.Now()
	post.Updated = time.Now()

	return s.repo.Create(post)
}

func (s *postService) FindPostByID(id string) (*Post, error) {
	post, err := s.repo.FindByID(id)
	return post, err
}

func (s *postService) FindAllPosts() ([]*Post, error) {
	posts, err := s.repo.FindAll()
	return posts, err
}

func (s *postService) UpdatePost(post *Post) error {
	post.Updated = time.Now()
	return s.repo.Update(post)
}

func (s *postService) DeletePost(id string) error {
	return s.repo.Delete(id)
}
