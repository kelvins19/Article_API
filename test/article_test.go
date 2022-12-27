package article

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

type Article struct {
	ID     int64  `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func TestCreateArticle(t *testing.T) {

	// Create a new article
	article := Article{Author: "Spongebob", Title: "Test Title", Body: "Test Content"}
	body, err := json.Marshal(article)
	if err != nil {
		t.Fatalf("Error marshaling JSON: %v", err)
	}

	// Send a request to the server to create the article
	res, err := http.Post("http://127.0.0.1:8000/articles", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer res.Body.Close()

	// Check that the response has the correct status code
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, res.StatusCode)
	}
}

func TestListArticles(t *testing.T) {

	// Send a request to the server to list the articles
	res, err := http.Get("http://127.0.0.1:8000/articles")
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer res.Body.Close()

	// Check that the response has the correct status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	// Parse the response body and check that it contains the test article
	var articles []Article
	err = json.NewDecoder(res.Body).Decode(&articles)
	if err != nil {
		t.Fatalf("Error parsing response body: %v", err)
	}
	if articles[0].Author != "author" || articles[0].Title != "Hello World" || articles[0].Body != "Hello World" {
		t.Errorf("Unexpected response body: %v", articles)
	}
}
