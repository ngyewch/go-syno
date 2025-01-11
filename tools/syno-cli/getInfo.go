package main

import (
	"github.com/ngyewch/go-syno/api"
	"github.com/ngyewch/go-syno/api/filestation"
	"github.com/urfave/cli/v2"
)

func doGetInfo(cCtx *cli.Context) error {
	return withClient(cCtx, func(c *api.Client) error {
		fileStationApi := filestation.New(c)

		getInfoResponse, err := fileStationApi.GetInfo(filestation.GetInfoRequest{
			Path:       cCtx.Args().Slice(),
			Additional: []string{"size", "time"},
		})
		if err != nil {
			return err
		}

		err = dump(getInfoResponse)
		if err != nil {
			return err
		}

		return nil
	})
}
