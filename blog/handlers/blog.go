package handlers

import (
	"go-blog/data"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var Posts = []data.Post{} // Temporary in-memory store

func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", Posts)
}

func ViewPostHandler(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 || id >= len(Posts) {
		c.String(http.StatusNotFound, "Post not found")
		return
	}
	c.HTML(http.StatusOK, "post.html", Posts[id])
}

func NewPostHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "new.html", nil)
}

func CreatePostHandler(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")

	post := data.Post{
		ID:      len(Posts),
		Title:   title,
		Content: content,
	}
	Posts = append(Posts, post)

	c.Redirect(http.StatusSeeOther, "/")
}
