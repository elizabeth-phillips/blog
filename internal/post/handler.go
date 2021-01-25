package post

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Handler creates the endpoints for reading the posts
type Handler interface {
	Get(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type postHandler struct {
	postService Service
}

//NewHandler is the constructor for endpoints
func NewHandler(postService Service) Handler {
	return &postHandler{
		postService,
	}
}

func (h *postHandler) Get(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.FindAllPosts()

	if err != nil || len(posts) <= 0 {
		logrus.Info("Handler: No rows found", err, posts)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("No rows found"))
	} else {
		response, _ := json.Marshal(posts)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	}
}

func (h *postHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	post, err := h.postService.FindPostByID(id)
	if err != nil {
		logrus.Error("Handler: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("%s ID is not found", id)))
	} else {
		response, err := json.Marshal(post)
		if err != nil {
			logrus.Error("Handler: ", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(fmt.Sprintf("ID %s is not found", id)))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(response)
		}
	}
}

func (h *postHandler) Create(w http.ResponseWriter, r *http.Request) {

	var post Post
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&post)
	err := h.postService.CreatePost(&post)

	if err != nil {
		logrus.Error("Handler: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Unable to create post"))
	} else {
		response, _ := json.Marshal(post)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(response)
	}
}

func (h *postHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var post Post
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&post)
	post.ID = id
	_ = h.postService.UpdatePost(&post)

	response, _ := json.Marshal(post)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

}

func (h *postHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.postService.DeletePost(id)

	response, _ := json.Marshal(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

}
