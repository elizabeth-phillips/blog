package post

import (
	"time"

	"github.com/google/uuid"
)

//Service creates the business logic that is a concrete implementation of how a service will interact with the PostDataStore methods
type Service interface {
	CreatePost(post *Post) error
	FindPostByID(id string) (*Post, error)
	FindAllPosts() ([]*Post, error)
	UpdatePost(post *Post) error
	DeletePost(id string) error
}

type postService struct {
	dataStore PostDataStore
}

//NewService is a constructor for Service
func NewService(dataStore PostDataStore) Service {
	return &postService{
		dataStore,
	}
}

func (s *postService) CreatePost(post *Post) error {
	post.ID = uuid.New().String()
	post.Created = time.Now()
	post.Updated = time.Now()

	return s.dataStore.Create(post)
}

func (s *postService) FindPostByID(id string) (*Post, error) {
	post, err := s.dataStore.FindByID(id)
	return post, err
}

func (s *postService) FindAllPosts() ([]*Post, error) {
	posts, err := s.dataStore.FindAll()
	return posts, err
}

func (s *postService) UpdatePost(post *Post) error {
	post.Updated = time.Now()
	return s.dataStore.Update(post)
}

func (s *postService) DeletePost(id string) error {
	return s.dataStore.Delete(id)
}
