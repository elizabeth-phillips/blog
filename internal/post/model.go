package post

import "time"

//Post creates the model that describes the attributes of a post.
type Post struct {
	ID      string    `json:"id" db:"id"`
	Creator string    `json:"creator" db:"creator"`
	Title   string    `json:"title" db:"title"`
	Content string    `json:"content" db:"content"`
	Created time.Time `json:"created" db:"created"`
	Updated time.Time `json:"updated" db:"updated"`
}
