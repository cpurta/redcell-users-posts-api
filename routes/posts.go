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

type PostsResource struct {
	postStore store.PostStore
	logger    *zap.Logger
}

func NewPostsResource(postStore store.PostStore, logger *zap.Logger) *PostsResource {
	return &PostsResource{
		postStore: postStore,
		logger:    logger,
	}
}

func (resource *PostsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", resource.ListPosts)
	r.Post("/", resource.CreatePost)

	r.Route("/{id}", func(r chi.Router) {
		postExistsMiddleware := middleware.NewPostExistsMiddleware(resource.postStore, resource.logger.Named("post_middleware"))

		r.Use(postExistsMiddleware.PostExists)
		r.Get("/", resource.GetPost)
		r.Put("/", resource.UpdatePost)
		r.Delete("/", resource.DeletePost)
	})

	return r
}

func (resource *PostsResource) ListPosts(w http.ResponseWriter, r *http.Request) {
	users, err := resource.postStore.ListPosts()
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

func (resource *PostsResource) CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		resource.logger.Error("unable to read request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var post *model.Post

	if err = json.Unmarshal(body, &post); err != nil {
		resource.logger.Error("unable to unmarshal body into post", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	created, err := resource.postStore.CreatePost(post)
	if err != nil {
		resource.logger.Error("unable to create post", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(created)
	if err != nil {
		resource.logger.Error("unable to marshal created post", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

func (resource *PostsResource) GetPost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad user id sent in request"))
		return
	}

	post, err := resource.postStore.GetPost(postIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to get user at this time"))
		return
	}

	responseBody, err := json.Marshal(post)
	if err != nil {
		resource.logger.Error("unable to marshal post", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

func (resource *PostsResource) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	postIDInt, err := strconv.Atoi(postID)
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

	var post *model.Post

	if err = json.Unmarshal(body, &post); err != nil {
		resource.logger.Error("unable to unmarshal body into user", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	post.ID = postIDInt

	updatedUser, err := resource.postStore.UpdatePost(post)
	if err != nil {
		resource.logger.Error("unable to update post", zap.Int("post_id", post.ID), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to get updated post at this time"))
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

func (resource *PostsResource) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad user id sent in request"))
		return
	}

	err = resource.postStore.DeletePost(postIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to delete user at this time"))
		return
	}
}
