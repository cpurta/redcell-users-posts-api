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

type UserExistsMiddleware struct {
	userStore store.UserStore
	logger    *zap.Logger
}

func NewUserExitsMiddleware(userStore store.UserStore, logger *zap.Logger) *UserExistsMiddleware {
	return &UserExistsMiddleware{
		userStore: userStore,
		logger:    logger,
	}
}

func (middleware *UserExistsMiddleware) UserExists(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			middleware.logger.Error("unable to convert userid provied", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid user_id provided"))
			return
		}

		_, err = middleware.userStore.GetUser(userIDInt)
		if err != nil && err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("user with user_id: %d does not exist", userIDInt)))
			return
		} else if err != nil && err != sql.ErrNoRows {
			middleware.logger.Error("unable to get user for existance check", zap.Int("user_id", userIDInt), zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
