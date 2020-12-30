package post_test

import (
	"testing"

	"github.com/elizabeth-phillips/blog/internal/mocks"
	"github.com/elizabeth-phillips/blog/internal/post"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestPostServiceSuite(t *testing.T) {
	suite.Run(t, new(PostServiceTestSuite))
}

type PostServiceTestSuite struct {
	suite.Suite
	postRepo  *mocks.MockRepository
	underTest post.Service
}

func (suite *PostServiceTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	suite.postRepo = mocks.NewMockRepository(mockCtrl)
	suite.underTest = post.NewService(suite.postRepo)
}

func (suite *PostServiceTestSuite) TestCreate() {
	//Arrange
	t := &post.Post{
		Creator: "Joel",
	}
	suite.postRepo.EXPECT().Create(gomock.Any()).Return(nil)

	//Act
	err := suite.underTest.CreatePost(t)

	//Assert
	suite.NoError(err, "Shouldn't error")
	suite.NotNil(t.ID, "should not be null")
	suite.NotNil(t.Created, "should not be null")
	suite.NotNil(t.Updated, "should not be null")

}

func (suite *PostServiceTestSuite) TestFindPostById() {
	t := &post.Post{
		ID:      "test",
		Creator: "Joel",
	}
	suite.postRepo.EXPECT().FindByID("test").Return(t, nil)

	result, err := suite.underTest.FindPostByID("test")

	suite.NoError(err, "Shouldn't error")
	suite.Equal(t, result, "should be pushing value returned from repo")
}

func (suite *PostServiceTestSuite) TestFindAll() {
	ts := []*post.Post{
		{
			ID:      "test1",
			Creator: "Joel",
		},
		{
			ID:      "test2",
			Creator: "Other",
		},
	}
	suite.postRepo.EXPECT().FindAll().Return(ts, nil)

	result, err := suite.underTest.FindAllPosts()

	suite.NoError(err, "Shouldn't error")
	suite.Len(result, 2, "Should get two results")
}
