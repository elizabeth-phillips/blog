package psql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/elizabeth-phillips/blog/internal/post"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq" //gives access to the PSQL driver
)

type postRepository struct {
	db *sql.DB
}

//NewPostgresRepository creates a connection with the postgres post data store
func NewPostgresRepository(db *sql.DB) *postRepository {
	return &postRepository{
		db,
	}
}

func (t postRepository) Create(post *post.Post) (err error) {
	sqlStatement := `INSERT INTO post(creator, title, content, created, updated) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	row := t.db.QueryRow(sqlStatement, post.Creator, post.Title, post.Content, time.Now(), time.Now())
	err = row.Scan(&post.ID)
	if err == nil {
		logrus.Info("DB: New record ID is:", post.ID)
	} else {
		logrus.Error("DB: ", err)
	}

	return
}

func (t postRepository) FindByID(id string) (foundPost *post.Post, err error) {
	foundPost = &post.Post{}
	sqlStatement := `SELECT * FROM post WHERE id=$1`
	row := t.db.QueryRow(sqlStatement, id)
	err = row.Scan(&foundPost.ID, &foundPost.Creator, &foundPost.Title, &foundPost.Content, &foundPost.Created, &foundPost.Updated)

	if err != nil {
		logrus.Error(fmt.Sprintf("DB: No post with the id %s", id), err)
		return foundPost, fmt.Errorf("No post with the id %s", id)
	}
	return
}

func (t postRepository) FindByDateRange(start time.Time, end time.Time) (allPosts []*post.Post, err error) {
	sqlStatement := `SELECT * FROM post WHERE created BETWEEN $1 AND $2`
	rows, err := t.db.Query(sqlStatement, start, end)
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		logrus.Error("DB: No rows were returned", err)
	} else {
		for rows.Next() {
			currentPost := &post.Post{}
			err = rows.Scan(&currentPost.ID, &currentPost.Creator, &currentPost.Title, &currentPost.Content, &currentPost.Created, &currentPost.Updated)
			if err == nil {
				allPosts = append(allPosts, currentPost) //append only if there was not an error
			}
		}
	}

	return
}

func (t postRepository) FindAll() (allPosts []*post.Post, err error) {
	sqlStatement := `SELECT * FROM post`
	rows, err := t.db.Query(sqlStatement)
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		logrus.Error("DB: No rows were returned", err)
	} else {
		for rows.Next() {
			currentPost := &post.Post{}
			err = rows.Scan(&currentPost.ID, &currentPost.Creator, &currentPost.Title, &currentPost.Content, &currentPost.Created, &currentPost.Updated)
			if err == nil {
				allPosts = append(allPosts, currentPost) //append only if there was not an error
			}
		}
	}

	return
}

func (t postRepository) Update(tool *post.Post) (err error) {
	found, err := t.FindByID(tool.ID) //Checks to see if ID is provided/valid
	if err != nil {
		return fmt.Errorf("update: id not found")
	}
	sqlStatement := `UPDATE post
					SET creator = $2, title = $3, content = $4, updated = $5
					WHERE id = $1;`
	result, err := t.db.Exec(sqlStatement, found.ID, tool.Creator, tool.Title, tool.Content, time.Now())
	if err == nil {
		numRows, _ := result.RowsAffected()
		logrus.Infof("%v row(s) affected\n", numRows)
	} else {
		logrus.Error("DB: ", err)
	}
	return err
}

func (t postRepository) Delete(id string) (err error) {
	_, err = t.FindByID(id)
	sqlStatement := `DELETE FROM post
					WHERE id = $1;`
	result, err := t.db.Exec(sqlStatement, id)
	numRows, _ := result.RowsAffected()
	if err == nil {
		logrus.Info("%v rows deleted\n", numRows)
	} else {
		logrus.Error("DB: ", err)
	}
	return
}
