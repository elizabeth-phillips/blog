package post

type PostDataStore interface {
	Create(Post *Post) error
	FindByID(id string) (*Post, error)
	FindAll() ([]*Post, error)
	Update(Post *Post) error
	Delete(id string) error
}
