package api

import (
	"strings"
)

type InfoApi struct {
	client *Client
}

type QueryRequest struct {
	ApiNames []string
}

type QueryResponse map[string]*APIDescription

func (c *Client) query(req QueryRequest) (*Response[QueryResponse], error) {
	paramMap := make(map[string]string)
	if len(req.ApiNames) > 0 {
		paramMap["query"] = strings.Join(req.ApiNames, ",")
	} else {
		paramMap["query"] = "all"
	}

	var res Response[QueryResponse]
	err := c.doRequest("query.cgi", "SYNO.API.Info", 1, "query", paramMap, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
