package main

import (
	"github.com/ngyewch/go-syno/api"
	"github.com/ngyewch/go-syno/api/filestation"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"path/filepath"
)

func doDownload(cCtx *cli.Context) error {
	return withClient(cCtx, func(c *api.Client) error {
		fileStationApi := filestation.New(c)

		r, err := fileStationApi.Download(filestation.DownloadRequest{
			Path: cCtx.Args().Slice(),
			Mode: "download",
		})
		if err != nil {
			return err
		}
		defer func(r io.ReadCloser) {
			_ = r.Close()
		}(r)

		var w io.WriteCloser
		if cCtx.NArg() > 1 {
			w, err = os.Create("download.zip")
			if err != nil {
				return err
			}
		} else {
			w, err = os.Create(filepath.Base(cCtx.Args().First()))
			if err != nil {
				return err
			}
		}
		defer func(w io.WriteCloser) {
			_ = w.Close()
		}(w)

		_, err = io.Copy(w, r)
		return err
	})
}
