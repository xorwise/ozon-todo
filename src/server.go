package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/xorwise/ozon-todo/graph"
	"github.com/xorwise/ozon-todo/graph/resolvers"
	"github.com/xorwise/ozon-todo/middlewares"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	mode := os.Getenv("MODE")
	var db *sql.DB
	if mode == "postgres" {
		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbname := os.Getenv("POSTGRES_DB")
		fmt.Println(host, port, user, password, dbname)
		d, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
		if err != nil {
			log.Fatal(err)
		}
		err = initPostgresDB(d)
		if err != nil {
			log.Fatal(err)
		}
		db = d
	} else if mode == "memory" {
		d, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			log.Fatal(err)
		}
		err = initMemoryDB(d)
		if err != nil {
			log.Fatal(err)
		}
		db = d
	} else {
		log.Fatal("MODE must be 'postgres' or 'memory'")
	}

	defer db.Close()

	rootHandler := handler.NewDefaultServer(graph.NewExecutableSchema(resolvers.NewResolver(db)))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middlewares.DataloaderMiddleware(db, middlewares.AuthMiddleware(rootHandler)))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initMemoryDB(db *sql.DB) error {
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
	_, err = db.Exec("INSERT INTO users (id, username) VALUES (1, 'root')")
	err = db.Ping()
	return err
}

func initPostgresDB(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users CASCADE")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS posts CASCADE")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS comments CASCADE")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE users (id SERIAL PRIMARY KEY, username TEXT NOT NULL)")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE posts (id SERIAL PRIMARY KEY, title TEXT NOT NULL, content TEXT NOT NULL, comments_allowed BOOLEAN NOT NULL, author_id INTEGER NOT NULL, FOREIGN KEY(author_id) REFERENCES users(id))")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE comments (id SERIAL PRIMARY KEY, content VARCHAR(2000) NOT NULL, post_id INTEGER NOT NULL, parent_id INTEGER, author_id INTEGER NOT NULL, created_at INTEGER NOT NULL, FOREIGN KEY(post_id) REFERENCES posts(id), FOREIGN KEY(parent_id) REFERENCES comments(id), FOREIGN KEY(author_id) REFERENCES users(id))")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (id, username) VALUES (1, 'root')")
	err = db.Ping()
	return err
}
