package info

import (
	"github.com/ngyewch/go-syno/api"
	"strings"
)

type InfoApi struct {
	client *api.Client
}

type QueryRequest struct {
	ApiNames []string
}

type QueryResponse map[string]APIDescription

type APIDescription struct {
	Path       string `json:"path"`
	MinVersion int    `json:"minVersion"`
	MaxVersion int    `json:"maxVersion"`
}

func NewInfoApi(client *api.Client) *InfoApi {
	return &InfoApi{
		client: client,
	}
}

func (c *InfoApi) Query(req QueryRequest) (*api.Response[QueryResponse], error) {
	paramMap := make(map[string]string)
	if len(req.ApiNames) > 0 {
		paramMap["query"] = strings.Join(req.ApiNames, ",")
	} else {
		paramMap["query"] = "all"
	}

	var res api.Response[QueryResponse]
	err := c.client.Request("query.cgi", "SYNO.API.Info", "1", "query", paramMap, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
