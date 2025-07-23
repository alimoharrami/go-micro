package handlers

import (
	"go-blog/data"
	"html/template"
	"net/http"
	"strconv"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))
var Posts = []data.Post{} // Temporary in-memory store

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", Posts)
}

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 || id >= len(Posts) {
		http.NotFound(w, r)
		return
	}
	templates.ExecuteTemplate(w, "post.html", Posts[id])
}

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "new.html", nil)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	post := data.Post{
		ID:      len(Posts),
		Title:   title,
		Content: content,
	}
	Posts = append(Posts, post)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
