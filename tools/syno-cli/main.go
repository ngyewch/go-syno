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
		Flags: []cli.Flag{
			baseUrlFlag,
			usernameFlag,
			passwordFlag,
		},
		Commands: []*cli.Command{
			{
				Name:   "list-share",
				Usage:  "list share",
				Action: doListShare,
			},
			{
				Name:   "list",
				Usage:  "list",
				Action: doList,
			},
			{
				Name:   "get-info",
				Usage:  "get info",
				Action: doGetInfo,
			},
			{
				Name:   "download",
				Usage:  "download",
				Action: doDownload,
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
