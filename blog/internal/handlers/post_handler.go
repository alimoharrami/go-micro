package handlers

import (
	"go-blog/internal/domain"
	"go-blog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	service *service.PostService
}

func NewPostController(service *service.PostService) *PostController {
	return &PostController{service}
}

func (uc *PostController) GetPostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	post, err := uc.service.GetByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, post)
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
	var posts []domain.Post
	posts, err := uc.service.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, posts)
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
