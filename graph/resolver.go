package graph

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func NewResolver(db *sql.DB) Config {
	c := Config{
		Resolvers: &Resolver{DB: db},
	}
	c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		ctxUserID := ctx.Value("userID")
		if ctxUserID != nil {
			return next(ctx)
		} else {
			return nil, fmt.Errorf("not authenticated")
		}
	}
	return c
}

type Resolver struct {
	DB *sql.DB
}
