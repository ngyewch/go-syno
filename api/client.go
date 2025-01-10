package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	baseUrl    string
	httpClient *http.Client
}

type Response[T any] struct {
	Success bool   `json:"success"`
	Error   *Error `json:"error,omitempty"`
	Data    *T     `json:"data,omitempty"`
}

type Error struct {
	Code   int          `json:"code"`
	Errors []ErrorEntry `json:"errors,omitempty"`
}

type ErrorEntry struct {
	Code int    `json:"code"`
	Path string `json:"path,omitempty"`
}

func NewClient(baseUrl string, httpClient *http.Client) (*Client, error) {
	_, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		baseUrl:    baseUrl,
		httpClient: httpClient,
	}, nil
}

func (c *Client) Request(relativePath string, api string, version string, method string, paramMap map[string]string, response any) error {
	baseUrl, err := url.Parse(c.baseUrl)
	if err != nil {
		return err
	}
	requestUrl := baseUrl.ResolveReference(&url.URL{Path: fmt.Sprintf("/webapi/%s", relativePath)})

	q := requestUrl.Query()
	q.Set("api", api)
	q.Set("version", version)
	q.Set("method", method)
	for k, v := range paramMap {
		q.Add(k, v)
	}
	requestUrl.RawQuery = q.Encode()

	fmt.Println(requestUrl.String())
	httpResponse, err := c.httpClient.Get(requestUrl.String())
	if err != nil {
		return err
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(httpResponse.Body)

	jsonDecoder := json.NewDecoder(httpResponse.Body)
	err = jsonDecoder.Decode(response)
	if err != nil {
		return err
	}
	return nil
}
