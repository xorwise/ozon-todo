package model

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postId"`
	ParentID  *int      `json:"-,omitempty"`
	AuthorID  int       `json:"-"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%d", t.Unix()))
	})
}

func UnmarshalTime(v interface{}) (time.Time, error) {
	return time.Unix(v.(int64), 0), nil
}

type Post struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	AuthorID        int    `json:"-"`
	CommentsAllowed bool   `json:"commentsAllowed"`
	CommentsID      []int  `json:"-,omitempty"`
}

type AddPostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type AddCommentInput struct {
	PostID   int    `json:"postId"`
	ParentID *int   `json:"parentId,omitempty"`
	Content  string `json:"content"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
