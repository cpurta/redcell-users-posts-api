package start

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"redcellpartners.com/users-posts-api/routes"
	"redcellpartners.com/users-posts-api/store/postgres"
)

const DEFAULT_TIMEOUT = time.Second * 60

type StartRunner struct {
	ListenAddr string

	PostgresHost     string
	PostgresPort     int
	PostgresUsername string
	PostgresPassword string
	PostgresDatabase string
	PostgresSSLMode  string

	LoggingProduction bool
	LoggingLevel      string

	logger *zap.Logger
}

func (runner *StartRunner) Run(cliContext *cli.Context) error {
	var (
		err error
	)

	loggerConfig := zap.NewDevelopmentConfig()

	if runner.LoggingProduction {
		loggerConfig = zap.NewProductionConfig()
	}

	if err = loggerConfig.Level.UnmarshalText([]byte(runner.LoggingLevel)); err != nil {
		log.Fatalf("unable to unmarshal zap logging level: %s", err.Error())
	}

	runner.logger, err = loggerConfig.Build()
	if err != nil {
		log.Fatalf("unable to build zap logger: %s", err.Error())
	}

	defer func() {
		if err := runner.logger.Sync(); err != nil {
			log.Fatalf("error syncing logger: %s", err.Error())
		}
	}()

	// connect to postgres
	dbConnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		runner.PostgresHost,
		runner.PostgresPort,
		runner.PostgresUsername,
		runner.PostgresPassword,
		runner.PostgresDatabase,
		runner.PostgresSSLMode,
	)

	runner.logger.Debug("checking db connection string", zap.String("db_connection_str", dbConnString))

	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("unable to connect to postgres database: %s", err.Error())
	}

	userStore, err := postgres.NewPostgresUserClient(db, runner.logger.Named("user_postgres_client"))
	if err != nil {
		log.Fatalf("unable to create new postgres user client: %s", err.Error())
	}

	postsStore, err := postgres.NewPostgresPostClient(db, runner.logger.Named("post_postgres_client"))
	if err != nil {
		log.Fatalf("unable to create new postgres post client: %s", err.Error())
	}

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(DEFAULT_TIMEOUT))

	usersResource := routes.NewUsersResource(userStore, runner.logger.Named("users_resource"))

	postsResource := routes.NewPostsResource(postsStore, runner.logger.Named("posts_resource"))

	router.Mount("/users", usersResource.Routes())
	router.Mount("/posts", postsResource.Routes())

	runner.logger.Info("starting users-posts-api REST API server")

	if err = http.ListenAndServe(runner.ListenAddr, router); err != nil {
		runner.logger.Error("error listening and serving user posts router", zap.Error(err))
	}

	return nil
}
