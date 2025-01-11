package auth

import (
	"github.com/ngyewch/go-syno/api"
)

type Api struct {
	client *api.Client
}

type LoginRequest struct {
	Account string
	Passwd  string
	Session string
}

type LoginResponse struct {
	Sid string `json:"sid"`
}

type LogoutRequest struct {
	Session string
}

type LogoutResponse struct{}

func New(client *api.Client) (*Api, error) {
	return &Api{
		client: client,
	}, nil
}

func (a *Api) Login(req LoginRequest) (*api.Response[LoginResponse], error) {
	paramMap := make(map[string]string)
	paramMap["account"] = req.Account
	paramMap["passwd"] = req.Passwd
	paramMap["session"] = req.Session
	paramMap["format"] = "sid"

	var res api.Response[LoginResponse]
	err := a.client.Request("SYNO.API.Auth", 3, "login", paramMap, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (a *Api) Logout(req LogoutRequest) (*api.Response[LogoutResponse], error) {
	paramMap := make(map[string]string)
	paramMap["session"] = req.Session

	var res api.Response[LogoutResponse]
	err := a.client.Request("SYNO.API.Auth", 1, "logout", paramMap, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
