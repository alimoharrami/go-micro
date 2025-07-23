package main

import (
	"net/http"
	"strconv"
	"text/template"
)

type Post struct {
	ID      int
	Title   string
	Content string
}

var posts = []Post{}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/post", viewPostHandler)
	http.HandleFunc("/new", newPostHandler)
	http.HandleFunc("/create", createPostHandler)
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", posts)
}

func viewPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 || id >= len(posts) {
		http.NotFound(w, r)
		return
	}
	templates.ExecuteTemplate(w, "post.html", posts[id])
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "new.html", nil)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	post := Post{
		ID:      len(posts),
		Title:   title,
		Content: content,
	}
	posts = append(posts, post)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
