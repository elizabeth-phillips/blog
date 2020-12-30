package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/elizabeth-phillips/blog/internal/database/psql"
	"github.com/elizabeth-phillips/blog/internal/post"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const DefaultPostgresUrl = "postgresql://postgres@localhost/post?sslmode=disable"

func main() {
	var dbURL string
	var postRepo post.PostDataStore

	// --------- Credentials ---------
	host := os.Getenv("BLOG_HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_PASSWORD")
	dbname := os.Getenv("BLOG_NAME")

	// --------- Create Connection ---------
	dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, "disable")

	pconn, err := postgresConnection(dbURL)

	if err != nil {
		pconn, err = postgresConnection(DefaultPostgresUrl)
	}

	defer pconn.Close()
	postRepo = psql.NewPostgresPostDataStore(pconn)

	postService := post.NewService(postRepo)
	postHandler := post.NewHandler(postService)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/posts", postHandler.Get).Methods("GET")
	router.HandleFunc("/posts/{id}", postHandler.GetByID).Methods("GET")
	router.HandleFunc("/posts", postHandler.Create).Methods("POST")

	errs := make(chan error, 2)
	go func() {
		logrus.Info("Listening server mode on port :3001")
		errs <- http.ListenAndServe(":3001", nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logrus.Errorf("terminated %s", <-errs)

}
func postgresConnection(database string) (db *sql.DB, err error) {
	logrus.Info("Connecting to PostgreSQL DB")
	db, err = sql.Open("postgres", database)
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
	return
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
