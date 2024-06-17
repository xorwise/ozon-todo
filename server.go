package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/xorwise/ozon-todo/graph"
	"github.com/xorwise/ozon-todo/graph/middlewares"

	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	err = initDB(db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rootHandler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.NewResolver(db)))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middlewares.AuthMiddleware(rootHandler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initDB(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS posts")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS comments")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT NOT NULL)")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE posts (id INTEGER PRIMARY KEY, title TEXT NOT NULL, content TEXT NOT NULL, comments_allowed INTEGER NOT NULL, author_id INTEGER NOT NULL, FOREIGN KEY(author_id) REFERENCES users(id))")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE comments (id INTEGER PRIMARY KEY, content TEXT NOT NULL, post_id INTEGER NOT NULL, parent_id INTEGER, author_id INTEGER NOT NULL, created_at INTEGER NOT NULL, FOREIGN KEY(post_id) REFERENCES posts(id), FOREIGN KEY(parent_id) REFERENCES comments(id), FOREIGN KEY(author_id) REFERENCES users(id))")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (id, username) VALUES (1, 'xorwise')")
	err = db.Ping()
	return err
}
