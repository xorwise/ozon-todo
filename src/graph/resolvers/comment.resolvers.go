package resolvers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xorwise/ozon-todo/dataloaders"
	"github.com/xorwise/ozon-todo/graph/model"
)

func (r *commentResolver) Parent(ctx context.Context, obj *model.Comment) (*model.Comment, error) {
	if obj.ParentID == nil {
		return nil, nil
	}
	row := r.DB.QueryRowContext(ctx, "SELECT id, post_id, parent_id, author_id, content, created_at FROM comments WHERE id = $1", *obj.ParentID)

	var comment model.Comment
	var unixTime int64
	if err := row.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.AuthorID, &comment.Content, &unixTime); err != nil {
		return nil, err
	}
	comment.CreatedAt = time.Unix(unixTime, 0)
	return &comment, nil
}

func (r *commentResolver) Author(ctx context.Context, obj *model.Comment) (*model.User, error) {
	user, err := ctx.Value("userloader").(*dataloaders.UserLoader).Load(obj.AuthorID)
	return user, err
}

func (r *mutationResolver) AddComment(ctx context.Context, input model.AddCommentInput) (*model.Comment, error) {
	userID, err := strconv.Atoi(ctx.Value("userID").(string))
	if err != nil {
		return &model.Comment{}, err
	}
	if len(input.Content) > 2000 {
		return &model.Comment{}, fmt.Errorf("comment too long")
	}
	newComment := &model.Comment{
		PostID:    input.PostID,
		ParentID:  input.ParentID,
		AuthorID:  userID,
		Content:   input.Content,
		CreatedAt: time.Now(),
	}

	rows, err := r.DB.ExecContext(
		ctx,
		`INSERT INTO comments (post_id, parent_id, author_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5)`, newComment.PostID,
		&newComment.ParentID,
		userID,
		newComment.Content,
		newComment.CreatedAt.Unix(),
	)

	if err != nil {
		return &model.Comment{}, err
	}

	insertedId, err := rows.LastInsertId()
	newComment.ID = int(insertedId)
	for id, observer := range r.newCommentCh {
		ids := strings.Split(id, "-")
		if ids[0] != fmt.Sprintf("%d", newComment.PostID) {
			continue
		}
		observer <- newComment
	}
	return newComment, nil
}

func (r *mutationResolver) ToggleComments(ctx context.Context, postID string, allow bool) (*model.Post, error) {
	userID := ctx.Value("userID").(string)
	var post model.Post

	err := r.DB.QueryRowContext(ctx, `SELECT id, title, content, author_id, comments_allowed FROM posts where id = $1`, postID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.CommentsAllowed,
	)
	if err != nil {
		return &model.Post{}, err
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		return &model.Post{}, err
	}
	if post.AuthorID != id {
		return &model.Post{}, fmt.Errorf("forbidden")
	}
	post.CommentsAllowed = allow
	return &post, nil
}

func (r *queryResolver) Comments(ctx context.Context, postID string, limit *int, offset *int) ([]*model.Comment, error) {
	var comments []*model.Comment
	rows, err := r.DB.QueryContext(ctx, `SELECT id, post_id, parent_id, author_id, content, created_at FROM comments WHERE post_id = $1 LIMIT $2 OFFSET $3`, postID, limit, offset)
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

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	id := fmt.Sprintf("%s-%s", postID, ctx.Value("user"))

	ch := make(chan *model.Comment, 1)
	go func() {
		<-ctx.Done()
	}()
	r.newCommentCh[id] = ch
	return ch, nil
}
