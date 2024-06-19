package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPost(t *testing.T) {
	// Prepare the GraphQL query
	query := `
		mutation {
			addPost(input: {
				title: "Test Post"
				content: "This is a test post"
			}) {
				id
				title
				content
				author {
					id
					username
				},
			}
		}`

	// Create a request to your GraphQL server
	reqBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/query", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var actual map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(actual)
	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAddCommentWithoutParent(t *testing.T) {
	// Prepare the GraphQL query
	query := `
		mutation {
			addComment(input: {
				postId: 1
				content: "Test Comment"
			}) {
				id
				content
				author {
					id
					username
				},
				parent {
					id
					content
				}
				createdAt
			}
		}`
	reqBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/query", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var actual map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(actual)
	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAddCommentWithParent(t *testing.T) {
	// Prepare the GraphQL query
	query := `
		mutation {
			addComment(input: {
				postId: 1
				content: "Test Comment"
				parentId: 1
			}) {
				id
				content
				author {
					id
					username
				},
				parent {
					id
					content
				}
				createdAt
			}
		}`
	reqBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/query", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var actual map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(actual)
	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestToggleComments(t *testing.T) {
	// Prepare the GraphQL query
	query := `
		mutation {
	toggleComments(postId: 1, allow: true) {
				id
				title
				commentsAllowed
			}
		}`
	reqBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/query", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var actual map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(actual)
	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
