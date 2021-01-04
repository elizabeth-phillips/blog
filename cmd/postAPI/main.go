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

//DefaultPostgresURL creates a default Postgres URL if one is not provided
const DefaultPostgresURL = "postgresql://postgres@localhost/post?sslmode=disable"

func main() {
	var dbURL string
	var postRepo post.Repository

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
		pconn, err = postgresConnection(DefaultPostgresURL)
	}

	defer pconn.Close()
	postRepo = psql.NewPostgresRepository(pconn)

	postService := post.NewService(postRepo)
	postHandler := post.NewHandler(postService)

	router := mux.NewRouter()
	rootRouter := router.PathPrefix("/api").Subrouter()
	rootRouter.HandleFunc("/", rootHandler).Methods("GET")

	postRouter := rootRouter.PathPrefix("/posts").Subrouter()
	postRouter.HandleFunc("/", postHandler.Create).Methods("POST")
	postRouter.HandleFunc("/", postHandler.Get).Methods("GET")
	postRouter.HandleFunc("/{id}", postHandler.GetByID).Methods("GET")
	postRouter.HandleFunc("/{id}", postHandler.Update).Methods("PUT", "PATCH")
	postRouter.HandleFunc("/{id}", postHandler.Delete).Methods("DELETE")

	errs := make(chan error, 2)
	go func() {
		logrus.Info("Listening server mode on port: 3001")
		errs <- http.ListenAndServe(":3001", accessControl(router))
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome to %s from OM!", r.URL.Path)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
