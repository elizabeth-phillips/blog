package main

import (
	"fmt"
	"os"

	"github.com/elizabeth-phillips/blog/internal/database/psql"
	"github.com/elizabeth-phillips/blog/internal/post"
)

func main() {
	var postDataStore post.PostDataStore

	// --------- Credentials ---------
	host := os.Getenv("BLOG_HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_PASSWORD")
	dbname := os.Getenv("BLOG_NAME")

	// --------- Create Connection ---------
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, "disable")
	db, err := psql.NewPostgresConnection(connection)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	postDataStore = psql.NewPostgresPostDataStore(db)

	postService := post.NewService(postDataStore)
	fmt.Println(postService)
}
