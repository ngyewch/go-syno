package main

import (
	"github.com/ngyewch/go-syno/api"
	"github.com/ngyewch/go-syno/api/filestation"
	"github.com/urfave/cli/v2"
)

func doListShare(cCtx *cli.Context) error {
	return withClient(cCtx, func(c *api.Client) error {
		fileStationApi := filestation.New(c)

		listShareResponse, err := fileStationApi.ListShare(filestation.ListShareRequest{})
		if err != nil {
			return err
		}

		err = dump(listShareResponse)
		if err != nil {
			return err
		}

		return nil
	})
}
