package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/elizabeth-phillips/blog/internal/database/psql"
	"github.com/elizabeth-phillips/blog/internal/post"
	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("BLOG_HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_PASSWORD")
	dbname := os.Getenv("BLOG_NAME")

	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, "disable")
	db, err := psql.NewPostgresConnection(connection)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	postDataStore := psql.NewPostgresRepository(db)

	postService := post.NewService(postDataStore)

	allTools, err := postService.FindAllPosts()
	if err != nil {
		fmt.Println(err)
	}
	for _, tool := range allTools {
		fmt.Println(tool)
	}
	fmt.Println(allTools)
}

func postgresConnection(database string) (db *sql.DB, err error) {
	fmt.Println("Connecting to PostgreSQL DB")
	db, err = sql.Open("postgres", database)
	if err != nil {
		fmt.Println(err)
	}
	return
}
