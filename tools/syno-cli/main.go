package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	baseUrlFlag = &cli.StringFlag{
		Name:    "base-url",
		Usage:   "base URL",
		EnvVars: []string{"SYNOLOGY_BASE_URL"},
	}
	usernameFlag = &cli.StringFlag{
		Name:    "username",
		Usage:   "username",
		EnvVars: []string{"SYNOLOGY_USERNAME"},
	}
	passwordFlag = &cli.StringFlag{
		Name:    "password",
		Usage:   "password",
		EnvVars: []string{"SYNOLOGY_PASSWORD"},
	}

	app = &cli.App{
		Name:  "syno-cli",
		Usage: "Synology CLI",
		Commands: []*cli.Command{
			{
				Name:   "login",
				Usage:  "login",
				Action: doLogin,
				Flags: []cli.Flag{
					baseUrlFlag,
					usernameFlag,
					passwordFlag,
				},
			},
			{
				Name:   "info",
				Usage:  "info",
				Action: doInfo,
				Flags: []cli.Flag{
					baseUrlFlag,
				},
			},
		},
	}
)

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
