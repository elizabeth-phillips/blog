package psql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/elizabeth-phillips/blog/internal/post"

	_ "github.com/lib/pq" //gives access to the PSQL driver
)

type postDataStore struct {
	db *sql.DB
}

//NewPostgresPostDataStore creates a connection with the postgres post data store
func NewPostgresPostDataStore(db *sql.DB) *postDataStore {
	return &postDataStore{
		db,
	}
}

func (t postDataStore) Create(post *post.Post) (err error) {
	sqlStatement := `INSERT INTO post(name, description, created_at, updated_at)
					 VALUES ($1, $2, $3, $4) RETURNING id`
	row := t.db.QueryRow(sqlStatement, post.Title, post.Content, time.Now(), time.Now())
	err = row.Scan(&post.ID)
	if err == nil {
		fmt.Println("New record ID is:", post.ID)
	}

	return
}

func (t postDataStore) FindByID(id string) (post *post.Post, err error) {
	sqlStatement := `SELECT * FROM post WHERE id=$1`
	row := t.db.QueryRow(sqlStatement, id)
	err = row.Scan(&post.ID, &post.Title, &post.Content, &post.Created, &post.Updated)

	switch err {
	case sql.ErrNoRows:
		return post, fmt.Errorf("no rows were returned")
	case nil:
		return
	default:
		return
	}
}

func (t postDataStore) FindAll() (allPosts []*post.Post, err error) {
	sqlStatement := `SELECT * FROM post`
	rows, err := t.db.Query(sqlStatement)
	defer rows.Close()
	for rows.Next() {
		currentPost := &post.Post{}
		err = rows.Scan(&currentPost.ID, &currentPost.Title, &currentPost.Content, &currentPost.Created, &currentPost.Updated)
		if err == nil {
			allPosts = append(allPosts, currentPost) //append only if there was not an error
		}
	}

	err = rows.Err()
	return
}

func (t postDataStore) Update(tool *post.Post) (err error) {
	_, err = t.FindByID(tool.ID) //Checks to see if ID is provided/valid
	if err != nil {
		return
	}
	sqlStatement := `UPDATE post
					SET name = $2, description = $3, updated_at = $4
					WHERE id = $1;`
	result, err := t.db.Exec(sqlStatement, tool.ID, tool.Title, tool.Content, time.Now())
	fmt.Printf("%v row affected\n", result)
	return err
}

func (t postDataStore) Delete(id string) (err error) {
	_, err = t.FindByID(id)
	if err != nil {
		return fmt.Errorf("delete: id not found")
	}
	sqlStatement := `DELETE FROM post
					WHERE id = $1;`
	result, err := t.db.Exec(sqlStatement, id)
	numRows, _ := result.RowsAffected()
	fmt.Printf("%v rows deleted", numRows)
	return err
}
