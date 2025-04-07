package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"redcellpartners.com/users-posts-api/store"
)

type PostExistsMiddleware struct {
	postStore store.PostStore
	logger    *zap.Logger
}

func NewPostExistsMiddleware(postStore store.PostStore, logger *zap.Logger) *PostExistsMiddleware {
	return &PostExistsMiddleware{
		postStore: postStore,
		logger:    logger,
	}
}

func (middleware *PostExistsMiddleware) PostExists(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "id")

		postIDInt, err := strconv.Atoi(postID)
		if err != nil {
			middleware.logger.Error("unable to convert post_id provied", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid post_id provided"))
			return
		}

		_, err = middleware.postStore.GetPost(postIDInt)
		if err != nil && err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("post with post_id: %d does not exist", postIDInt)))
			return
		} else if err != nil && err != sql.ErrNoRows {
			middleware.logger.Error("unable to get post for existance check", zap.Int("post_id", postIDInt), zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
