package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var (
	commonErrorCodes = map[int]string{
		100: "Unknown error",
		101: "No parameter of API, method or version",
		102: "The requested API does not exist",
		103: "The requested method does not exist",
		104: "The requested version does not support the functionality",
		105: "The logged in session does not have permission",
		106: "Session timeout",
		107: "Session interrupted by duplicate login",
		119: "SID not found",
	}
)

type Client struct {
	baseUrl    string
	httpClient *http.Client
	paramMap   map[string]string
	apiMap     map[string]*APIDescription
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

func (e *Error) Error() string {
	return fmt.Sprintf("Synology API error: %d", e.Code)
}

type ErrorEntry struct {
	Code int    `json:"code"`
	Path string `json:"path,omitempty"`
}

type APIDescription struct {
	Path       string `json:"path"`
	MinVersion int    `json:"minVersion"`
	MaxVersion int    `json:"maxVersion"`
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
		paramMap:   make(map[string]string),
		apiMap:     make(map[string]*APIDescription),
	}, nil
}

func (c *Client) SetParam(key string, value string) {
	c.paramMap[key] = value
}

func (c *Client) Request(api string, version int, method string, paramMap map[string]string, response any) error {
	apiDescription, ok := c.apiMap[api]
	if !ok {
		queryResponse, err := c.query(QueryRequest{
			ApiNames: []string{api},
		})
		if err != nil {
			return err
		}
		apiDescription = (*queryResponse.Data)[api]
		c.apiMap[api] = apiDescription
	}
	return c.doRequest(apiDescription.Path, api, version, method, paramMap, response)
}

func (c *Client) RawRequest(api string, version int, method string, paramMap map[string]string) (io.ReadCloser, error) {
	apiDescription, ok := c.apiMap[api]
	if !ok {
		queryResponse, err := c.query(QueryRequest{
			ApiNames: []string{api},
		})
		if err != nil {
			return nil, err
		}
		apiDescription = (*queryResponse.Data)[api]
		c.apiMap[api] = apiDescription
	}
	return c.doRawRequest(apiDescription.Path, api, version, method, paramMap)
}

func (c *Client) doRawRequest(apiPath string, api string, version int, method string, paramMap map[string]string) (io.ReadCloser, error) {
	baseUrl, err := url.Parse(c.baseUrl)
	if err != nil {
		return nil, err
	}

	requestUrl := baseUrl.ResolveReference(&url.URL{Path: fmt.Sprintf("/webapi/%s", apiPath)})

	q := requestUrl.Query()
	q.Set("api", api)
	q.Set("version", strconv.Itoa(version))
	q.Set("method", method)
	for k, v := range c.paramMap {
		q.Set(k, v)
	}
	for k, v := range paramMap {
		q.Set(k, v)
	}
	requestUrl.RawQuery = q.Encode()

	httpResponse, err := c.httpClient.Get(requestUrl.String())
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode != 200 {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(httpResponse.Body)
		return nil, fmt.Errorf("HTTP status code %d", httpResponse.StatusCode)
	}

	return httpResponse.Body, nil
}

func (c *Client) doRequest(apiPath string, api string, version int, method string, paramMap map[string]string, response any) error {
	r, err := c.doRawRequest(apiPath, api, version, method, paramMap)
	if err != nil {
		return err
	}
	defer func(r io.ReadCloser) {
		_ = r.Close()
	}(r)

	var r0 io.Reader = r

	const trace = false
	if trace {
		httpResponseBytes, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		fmt.Println(string(httpResponseBytes))

		r0 = bytes.NewReader(httpResponseBytes)
	}

	jsonDecoder := json.NewDecoder(r0)
	err = jsonDecoder.Decode(response)
	if err != nil {
		return err
	}

	return nil
}
