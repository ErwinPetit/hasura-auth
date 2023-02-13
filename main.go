package main

import (
	"log"
	"os"

	"github.com/nhost/hasura-auth/cmd"
	"github.com/nhost/hasura-auth/controller"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{ //nolint: exhaustruct
		Name:                 "hasura-auth",
		EnableBashCompletion: true,
		Version:              controller.Version(),
		Description:          "Manages and operate tenants' infrastructure",
		Commands: []*cli.Command{
			cmd.Command(),
		},
		Flags: []cli.Flag{},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
