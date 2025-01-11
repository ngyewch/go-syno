package main

import (
	"github.com/ngyewch/go-syno/api"
	"github.com/ngyewch/go-syno/api/filestation"
	"github.com/urfave/cli/v2"
)

func doList(cCtx *cli.Context) error {
	return withClient(cCtx, func(c *api.Client) error {
		fileStationApi := filestation.New(c)

		listResponse, err := fileStationApi.List(filestation.ListRequest{
			FolderPath: cCtx.Args().First(),
			Additional: []string{"size", "time"},
		})
		if err != nil {
			return err
		}

		err = dump(listResponse)
		if err != nil {
			return err
		}

		return nil
	})
}
