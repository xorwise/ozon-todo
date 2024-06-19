package middlewares

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/xorwise/ozon-todo/dataloaders"
	"github.com/xorwise/ozon-todo/graph/model"
)

func DataloaderMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userConfig := dataloaders.UserLoaderConfig{
			Wait:     100 * time.Millisecond,
			MaxBatch: 100,
			Fetch: func(keys []int) ([]*model.User, []error) {
				var sqlQuery string
				var stringKeys string
				for i := 0; i < len(keys); i++ {
					if i < len(keys)-1 {
						stringKeys += strconv.Itoa(keys[i]) + ","
					} else {
						stringKeys += strconv.Itoa(keys[i])
					}
				}
				sqlQuery = "SELECT * FROM users WHERE id IN (" + stringKeys + ")"
				rows, err := db.Query(sqlQuery)
				if err != nil {
					return nil, []error{err}
				}
				var users []*model.User
				for rows.Next() {
					var user model.User
					if err := rows.Scan(&user.ID, &user.Username); err != nil {
						return nil, []error{err}
					}
					users = append(users, &user)
				}
				return users, nil
			},
		}
		userloader := dataloaders.NewUserLoader(userConfig)

		ctx := context.WithValue(r.Context(), "userloader", userloader)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
