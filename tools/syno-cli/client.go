package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ngyewch/go-syno/api"
	"github.com/ngyewch/go-syno/api/auth"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
)

func withClient(cCtx *cli.Context, handler func(c *api.Client) error) error {
	baseUrl := baseUrlFlag.Get(cCtx)
	username := usernameFlag.Get(cCtx)
	password := passwordFlag.Get(cCtx)

	c, err := api.NewClient(baseUrl, &http.Client{})
	if err != nil {
		return err
	}
	authApi, err := auth.New(c)
	if err != nil {
		return err
	}

	sessionId := uuid.New().String()

	loginResponse, err := authApi.Login(auth.LoginRequest{
		Account: username,
		Passwd:  password,
		Session: sessionId,
	})
	if err != nil {
		return err
	}

	c.SetParam("_sid", loginResponse.Data.Sid)

	defer func() {
		_ = func() error {
			_, err := authApi.Logout(auth.LogoutRequest{
				Session: sessionId,
			})
			return err
		}
	}()

	return handler(c)
}

func dump(o any) error {
	jsonEncoder := json.NewEncoder(os.Stdout)
	jsonEncoder.SetIndent("", "  ")
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(o)
	if err != nil {
		return err
	}
	return nil
}
