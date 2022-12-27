package article

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kelvins19/Article_API/constants"
)

type Article struct {
	ID     int64  `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type ApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ArticleAPI struct{}

func (articleApi *ArticleAPI) CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	var article Article
	var apiResponse ApiResponse

	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		constants.DB_HOST, constants.DB_PORT, constants.DB_USER, constants.DB_PASSWORD, constants.DB_NAME)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	_, err = db.Exec("INSERT INTO articles (author, title, body, created_at) VALUES ($1, $2, $3, now())",
		article.Author, article.Title, article.Body)
	if err != nil {
		apiResponse.Status = http.StatusInternalServerError
		apiResponse.Message = "There is an error when creating a new article"
		http.Error(w, err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiResponse)
		return
	}

	defer db.Close()

	apiResponse.Status = http.StatusCreated
	apiResponse.Message = "Your article has been created"

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiResponse)
}

func (articleApi *ArticleAPI) ListArticlesHandler(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	// Retrieve articles from database
	articles := []Article{}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		constants.DB_HOST, constants.DB_PORT, constants.DB_USER, constants.DB_PASSWORD, constants.DB_NAME)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	script := "SELECT id, author, title, body FROM articles"
	if params["author"] != nil && params["query"] != nil {
		script += " WHERE LOWER(author) LIKE '%" + params["author"][0] + "%' AND (LOWER(title) LIKE '%" + params["query"][0] + "%' OR LOWER(body) LIKE '%" + params["query"][0] + "%')"
		fmt.Print(script)
	} else {
		if params["author"] != nil {
			script += " WHERE LOWER(author) '%" + params["author"][0] + "%'"
			fmt.Print(script)
		}
		if params["query"] != nil {
			script += " WHERE LOWER(title) LIKE '%" + params["query"][0] + "%' OR LOWER(body) LIKE '%" + params["query"][0] + "%'"
			fmt.Print(script)
		}
	}

	err = db.Select(&articles, script)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer db.Close()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

var cache = &sync.Map{}

func getCachedArticle(id int) (*Article, bool) {
	if value, ok := cache.Load(id); ok {
		return value.(*Article), true
	}
	return nil, false
}

func setCachedArticle(article *Article) {
	cache.Store(article.ID, article)
}

func (articleApi *ArticleAPI) GetArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID of the article from the URL parameters
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Check the cache for the requested article
	articleCache, ok := getCachedArticle(id)
	if ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articleCache)
		return
	}

	// Retrieve the article from the database
	var article Article

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		constants.DB_HOST, constants.DB_PORT, constants.DB_USER, constants.DB_PASSWORD, constants.DB_NAME)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = db.Get(&article, "SELECT id, author, title, body FROM articles WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}

	defer db.Close()

	// Cache the article and return it
	setCachedArticle(&article)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}
