package main

import (
	"encoding/json"
	"github.com/ngyewch/go-syno/api"
	"github.com/ngyewch/go-syno/api/auth"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
)

func doLogin(cCtx *cli.Context) error {
	baseUrl := baseUrlFlag.Get(cCtx)
	username := usernameFlag.Get(cCtx)
	password := passwordFlag.Get(cCtx)
	c, err := api.NewClient(baseUrl, &http.Client{})
	if err != nil {
		return err
	}
	authApi := auth.NewAuthApi(c)

	loginResponse, err := authApi.Login(auth.LoginRequest{
		Account: username,
		Passwd:  password,
		Session: "FileStation1",
	})
	if err != nil {
		return err
	}

	jsonEncoder := json.NewEncoder(os.Stdout)
	jsonEncoder.SetIndent("", "  ")
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(loginResponse)
	if err != nil {
		return err
	}

	return nil
}
