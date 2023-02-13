package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/nhost/hasura-auth/controller"
	"github.com/urfave/cli/v2"
)

const (
	flagHasuraJSURL   = "hasura-js-url"
	flagLogFormatJSON = "log-format-json"
	flagDebug         = "debug"
	flagBind          = "bind"
	flagJWTSecret     = "jwt-secret"
)

func Command() *cli.Command { //nolint:funlen
	return &cli.Command{ //nolint: exhaustruct
		Name:  "serve",
		Usage: "serve the application",
		Flags: []cli.Flag{
			&cli.StringFlag{ //nolint: exhaustruct
				Name:     flagHasuraJSURL,
				Usage:    "Hasura JS URL",
				Value:    "http://localhost:4001",
				Category: "Javascript Service",
			},
			&cli.StringFlag{ //nolint: exhaustruct
				Name:     flagBind,
				Usage:    "bind address",
				Value:    ":4000",
				Category: "server",
			},
			&cli.BoolFlag{ //nolint: exhaustruct
				Name:     flagDebug,
				Usage:    "enable debug logging",
				Category: "general",
			},
			&cli.StringFlag{ //nolint: exhaustruct
				Name:     flagJWTSecret,
				Usage:    "JWT secret to use",
				Category: "JWT",
				Required: true,
				EnvVars:  []string{"JWT_SECRET", "HASURA_GRAPHQL_JWT_SECRET"},
			},
		},
		Action: serve,
	}
}

func serve(cCtx *cli.Context) error {
	rand.Seed(time.Now().UnixNano())

	logger := getLogger(cCtx.Bool(flagDebug), cCtx.Bool(flagLogFormatJSON))
	logger.Info(cCtx.App.Name + " v" + cCtx.App.Version)
	logFlags(logger, cCtx)

	authJS, err := url.Parse(cCtx.String(flagHasuraJSURL))
	if err != nil {
		return fmt.Errorf("invalid hasura-js-url: %w", err)
	}

	jwtSecret := controller.JWTSecret{}
	if err := json.Unmarshal([]byte(cCtx.String(flagJWTSecret)), &jwtSecret); err != nil {
		return fmt.Errorf("problem unmarshalling jwt secret: %w", err)
	}
	ctrl := controller.New(authJS.Scheme, authJS.Host, jwtSecret, logger)

	router, err := ctrl.SetupRouter([]string{}, "/", ginLogger(logger))
	if err != nil {
		return fmt.Errorf("failed to setup router: %w", err)
	}

	if err := router.Run(cCtx.String(flagBind)); err != nil {
		return fmt.Errorf("failed to run router: %w", err)
	}
	return nil
}
