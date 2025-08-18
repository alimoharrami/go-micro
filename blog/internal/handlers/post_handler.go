package handlers

import (
	"go-blog/external/protos/userpb"
	"go-blog/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	service    *service.PostService
	userClient userpb.UserServiceClient
}

func NewPostController(service *service.PostService, userClient userpb.UserServiceClient) *PostController {
	return &PostController{service, userClient}
}

func (uc *PostController) GetPostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	post, err := uc.service.GetByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	log.Printf("here")
	resp, err := uc.userClient.GetUser(c, &userpb.GetUserRequest{Id: c.Param("id")})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("here1")

	result := map[string]interface{}{
		"post": post,
		"user": resp,
	}

	c.JSON(http.StatusOK, result)
}

func (uc *PostController) CreatePost(c *gin.Context) {
	var input service.CreatePostInput // Use input struct instead of models.User directly
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call the service layer
	post, err := uc.service.Create(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with created user
	c.JSON(http.StatusCreated, post)
}

func (uc *PostController) ListPosts(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	result, err := uc.service.List(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}

func (uc *PostController) DeletePost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := uc.service.Delete(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, "Deleted Successfully!")
}

func (uc *PostController) UpdatePost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input service.UpdatePostInput // Use input struct instead of models.User directly
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	post, err := uc.service.Update(c, uint(id), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error updating post"})
		return
	}
	c.JSON(http.StatusOK, post)
}

//todo get posts paginate
