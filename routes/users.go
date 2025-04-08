package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"redcellpartners.com/users-posts-api/middleware"
	"redcellpartners.com/users-posts-api/model"
	"redcellpartners.com/users-posts-api/store"
)

type UsersResource struct {
	userStore store.UserStore
	logger    *zap.Logger
}

func NewUsersResource(userStore store.UserStore, logger *zap.Logger) *UsersResource {
	return &UsersResource{
		userStore: userStore,
		logger:    logger,
	}
}

func (resource *UsersResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", resource.ListUsers)
	r.Post("/", resource.CreateUser)

	r.Route("/{id}", func(r chi.Router) {
		userExistMiddleware := middleware.NewUserExitsMiddleware(resource.userStore, resource.logger.Named("user_middleware"))

		r.Use(userExistMiddleware.UserExists)
		r.Get("/", resource.GetUser)
		r.Put("/", resource.UpdateUser)
		r.Delete("/", resource.DeleteUser)
	})

	return r
}

func (resource *UsersResource) ListUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	users, err := resource.userStore.ListUsers()
	if err != nil {
		resource.logger.Error("unable to list users", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to list users at this time"))
		return
	}

	responseBytes, err := json.Marshal(users)
	if err != nil {
		resource.logger.Error("unable to marshal users", zap.Error(err))
		w.Write([]byte("unable to list users at this time"))
		return
	}

	w.Write(responseBytes)
}

func (resource *UsersResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		resource.logger.Error("unable to read request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var user *model.User

	if err = json.Unmarshal(body, &user); err != nil {
		resource.logger.Error("unable to unmarshal body into user", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	created, err := resource.userStore.CreateUser(user)
	if err != nil {
		resource.logger.Error("unable to create user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(created)
	if err != nil {
		resource.logger.Error("unable to marshal created user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseBody)
}

func (resource *UsersResource) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad user id sent in request"))
		return
	}

	user, err := resource.userStore.GetUser(userIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to get user at this time"))
		return
	}

	responseBody, err := json.Marshal(user)
	if err != nil {
		resource.logger.Error("unable to marshal user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

func (resource *UsersResource) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad user id sent in request"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		resource.logger.Error("unable to read request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var user *model.User

	if err = json.Unmarshal(body, &user); err != nil {
		resource.logger.Error("unable to unmarshal body into user", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.ID = userIDInt

	updatedUser, err := resource.userStore.UpdateUser(user)
	if err != nil {
		resource.logger.Error("unable to update user", zap.Int("user_id", user.ID), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to get updated user at this time"))
		return
	}

	responseBody, err := json.Marshal(updatedUser)
	if err != nil {
		resource.logger.Error("unable to marshal updated user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

func (resource *UsersResource) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad user id sent in request"))
		return
	}

	err = resource.userStore.DeleteUser(userIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to delete user at this time"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
