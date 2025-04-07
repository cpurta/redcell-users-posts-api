package start

import "github.com/urfave/cli"

func StartCommand() cli.Command {
	runner := &StartRunner{}

	return cli.Command{
		Name:        "start",
		Description: "starts the users and posts REST api",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "listen-addr",
				EnvVar:      "LISTEN_ADDR",
				Usage:       "address that the server will listen on",
				Value:       ":8080",
				Destination: &runner.ListenAddr,
			},
			cli.BoolFlag{
				Name:        "logging-production",
				EnvVar:      "LOGGING_PRODUCTION",
				Usage:       "enable logging for a system in production",
				Destination: &runner.LoggingProduction,
			},
			cli.StringFlag{
				Name:        "loggging-level",
				EnvVar:      "LOGGING_LEVEL",
				Usage:       "sets the loggging level of all logged messages",
				Destination: &runner.LoggingLevel,
			},
			cli.StringFlag{
				Name:        "postgres-username",
				EnvVar:      "POSTGRES_CONN_USERNAME",
				Usage:       "username for the postgres connection",
				Destination: &runner.PostgresUsername,
			},
			cli.StringFlag{
				Name:        "postgres-password",
				EnvVar:      "POSTGRES_CONN_PASSWORD",
				Usage:       "password for the postgres connection",
				Destination: &runner.PostgresPassword,
			},
			cli.StringFlag{
				Name:        "postgres-database",
				EnvVar:      "POSTGRES_CONN_DATABASE",
				Usage:       "database for the postgres connection",
				Destination: &runner.PostgresDatabase,
			},
			cli.StringFlag{
				Name:        "postgres-conn-host",
				EnvVar:      "POSTGRES_CONN_HOST",
				Usage:       "hostname of the postgres database",
				Destination: &runner.PostgresHost,
			},
			cli.IntFlag{
				Name:        "postgres-conn-port",
				EnvVar:      "POSTGRES_CONN_PORT",
				Usage:       "port of the postgres database",
				Destination: &runner.PostgresPort,
				Value:       5432,
			},
			cli.StringFlag{
				Name:        "postgres-conn-ssl-mode",
				EnvVar:      "POSTGRES_CONN_SSL_MODE",
				Usage:       "ssl mode of the postgres database",
				Destination: &runner.PostgresSSLMode,
				Value:       "disable",
			},
		},
		Action: runner.Run,
	}
}
