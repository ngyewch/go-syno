package main

import (
	"encoding/json"
	"github.com/ngyewch/go-syno/api"
	"github.com/ngyewch/go-syno/api/info"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
)

func doInfo(cCtx *cli.Context) error {
	baseUrl := baseUrlFlag.Get(cCtx)
	//username := usernameFlag.Get(cCtx)
	//password := passwordFlag.Get(cCtx)
	c, err := api.NewClient(baseUrl, &http.Client{})
	if err != nil {
		return err
	}
	infoApi := info.NewInfoApi(c)

	queryResponse, err := infoApi.Query(info.QueryRequest{})
	if err != nil {
		return err
	}

	jsonEncoder := json.NewEncoder(os.Stdout)
	jsonEncoder.SetIndent("", "  ")
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(queryResponse)
	if err != nil {
		return err
	}

	return nil
}
