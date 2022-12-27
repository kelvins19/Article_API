package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelvins19/Article_API/article"
	_ "github.com/lib/pq"
)

func main() {

	r := mux.NewRouter()

	article := &article.ArticleAPI{}

	// Create an article
	r.HandleFunc("/articles", article.CreateArticleHandler).Methods("POST")

	// List all articles
	r.HandleFunc("/articles", article.ListArticlesHandler).Methods("GET")
	r.HandleFunc("/articles/{id}", article.GetArticleHandler).Methods("GET")

	http.ListenAndServe(":8000", r)
}
