package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryPosts(t *testing.T) {
	// Prepare the GraphQL query
	query := `
        query {
            posts(limit: 10, offset: 0) {
                id,
                title,
                content,
                author {
                    id,
                    username
                }
            }
        }`

	// Create a request to your GraphQL server
	reqBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Assert the response body
	var actual map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	// Assert the response status code

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestQueryComments(t *testing.T) {
	// Prepare the GraphQL query
	query := `
		query {
			comments(postId: 1, limit: 10, offset: 0) {
				id,
				content,
				author {
					id,
					username
				},
				parent {
					id,
					content
				}
			}
		}`

	// Create a request to your GraphQL server
	reqBody, err := json.Marshal(map[string]string{"query": query})
	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var actual map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestQueryPost(t *testing.T) {
	// Prepare the GraphQL query
	query := `
		query {
			post(id: 1) {
				id,
				title,
				content,
				author {
					id,
					username
				},
				comments {
					id,
					content,
					author {
						id,
						username
					},
					parent {
						id,
						content
					}
				}
			}
		}`

	// Create a request to your GraphQL server
	reqBody, err := json.Marshal(map[string]string{"query": query})
	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var actual map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
