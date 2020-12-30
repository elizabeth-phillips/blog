package post

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//Handler creates the endpoints for reading the posts
type Handler interface {
	Get(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
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
	posts, _ := h.postService.FindAllPosts()

	response, _ := json.Marshal(posts)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *postHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	post, _ := h.postService.FindPostByID(id)

	response, _ := json.Marshal(post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *postHandler) Create(w http.ResponseWriter, r *http.Request) {

	var post Post
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&post)
	_ = h.postService.CreatePost(&post)

	response, _ := json.Marshal(post)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(response)

}
