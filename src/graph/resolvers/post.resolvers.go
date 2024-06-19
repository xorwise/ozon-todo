package resolvers

import (
	"context"
	"strconv"
	"time"

	"github.com/xorwise/ozon-todo/dataloaders"
	"github.com/xorwise/ozon-todo/graph/model"
)

func (r *mutationResolver) AddPost(ctx context.Context, input model.AddPostInput) (*model.Post, error) {
	userID, err := strconv.Atoi(ctx.Value("userID").(string))
	if err != nil {
		return &model.Post{}, err
	}
	newPost := &model.Post{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: userID,
	}

	rows, err := r.DB.ExecContext(
		ctx, `INSERT INTO posts (title, content, author_id, comments_allowed) 
			VALUES ($1, $2, $3, $4)`,
		newPost.Title,
		newPost.Content,
		userID,
		true,
	)

	if err != nil {
		return &model.Post{}, err
	}

	insertedId, err := rows.LastInsertId()
	newPost.ID = int(insertedId)
	return newPost, nil
}

func (r *postResolver) Author(ctx context.Context, obj *model.Post) (*model.User, error) {
	user, err := ctx.Value("userloader").(*dataloaders.UserLoader).Load(obj.AuthorID)
	return user, err
}

func (r *postResolver) Comments(ctx context.Context, obj *model.Post, limit *int, offset *int) ([]*model.Comment, error) {
	var comments []*model.Comment
	rows, err := r.DB.QueryContext(ctx, `SELECT id, post_id, parent_id, author_id, content, created_at FROM comments WHERE post_id = $1 LIMIT $2 OFFSET $3`, obj.ID, limit, offset)
	defer rows.Close()
	if err != nil {
		return comments, err
	}
	for rows.Next() {
		var unixTime int64
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.AuthorID, &comment.Content, &unixTime); err != nil {
			return comments, err
		}
		comment.CreatedAt = time.Unix(unixTime, 0)
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (r *queryResolver) Posts(ctx context.Context, limit *int, offset *int) ([]*model.Post, error) {
	var posts []*model.Post
	rows, err := r.DB.QueryContext(ctx, `SELECT id, title, content, author_id, comments_allowed FROM posts LIMIT $1 OFFSET $2`, limit, offset)
	defer rows.Close()
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CommentsAllowed); err != nil {
			return posts, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	var post model.Post
	err := r.DB.QueryRowContext(ctx, `SELECT id, title, content, author_id, comments_allowed FROM posts WHERE id = $1`, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.CommentsAllowed,
	)
	if err != nil {
		return &post, err
	}
	return &post, nil
}
