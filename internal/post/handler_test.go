package post_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elizabeth-phillips/blog/internal/mocks"
	"github.com/elizabeth-phillips/blog/internal/post"

	"github.com/gorilla/mux"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestPostHandlerSuite(t *testing.T) {
	suite.Run(t, new(PostHandlerTestSuite))
}

type PostHandlerTestSuite struct {
	suite.Suite
	postService *mocks.MockService
	underTest   post.Handler
}

func (suite *PostHandlerTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	suite.postService = mocks.NewMockService(mockCtrl)
	suite.underTest = post.NewHandler(suite.postService)
}

func (suite *PostHandlerTestSuite) TestCreate() {
	t := &post.Post{
		Creator: "Elizabeth",
	}
	suite.postService.EXPECT().CreatePost(gomock.Eq(t)).Return(nil)

	body, _ := json.Marshal(t)
	r, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	suite.underTest.Create(w, r)

	response := w.Result()
	suite.Equal("201 Created", response.Status)

	defer response.Body.Close()
	result := new(post.Post)
	json.NewDecoder(response.Body).Decode(result)

	suite.Equal("Elizabeth", result.Creator)
}

func (suite *PostHandlerTestSuite) TestFindPostByID() {
	t := &post.Post{
		Creator: "Elizabeth",
	}
	suite.postService.EXPECT().FindPostByID("test").Return(t, nil)

	vars := map[string]string{
		"id": "test",
	}

	r, _ := http.NewRequest("GET", "/posts/test", nil)
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()
	suite.underTest.GetByID(w, r)

	response := w.Result()
	suite.Equal("200 OK", response.Status)

	defer response.Body.Close()
	result := new(post.Post)
	json.NewDecoder(response.Body).Decode(result)

	suite.Equal("Elizabeth", result.Creator)
}

func (suite *PostHandlerTestSuite) TestFindAll() {
	ts := []*post.Post{
		{
			ID:      "test1",
			Creator: "Elizabeth",
		},
		{
			ID:      "test2",
			Creator: "Other",
		},
	}
	suite.postService.EXPECT().FindAllPosts().Return(ts, nil)

	r, _ := http.NewRequest("GET", "/posts", nil)

	w := httptest.NewRecorder()
	suite.underTest.Get(w, r)

	response := w.Result()
	suite.Equal("200 OK", response.Status)

	defer response.Body.Close()
	result := new([]post.Post)
	json.NewDecoder(response.Body).Decode(result)
	suite.Len(*result, 2, "Should get two results")
}
