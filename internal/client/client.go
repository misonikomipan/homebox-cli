package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/misonikomipan/homebox-cli/internal/config"
)

type Client struct {
	base       string
	token      string
	httpClient *http.Client
}

func New(authenticated bool) (*Client, error) {
	c := &Client{
		base:       config.GetEndpoint() + "/api",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	if authenticated {
		token := config.GetToken()
		if token == "" {
			return nil, fmt.Errorf("not authenticated — run 'hb login' first")
		}
		c.token = token
	}
	return c, nil
}

func (c *Client) do(method, path string, query url.Values, body any) ([]byte, int, error) {
	u := c.base + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, u, bodyReader)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return data, resp.StatusCode, nil
}

func (c *Client) Get(path string, query url.Values) ([]byte, error) {
	return c.request("GET", path, query, nil)
}

func (c *Client) Post(path string, body any) ([]byte, error) {
	return c.request("POST", path, nil, body)
}

func (c *Client) Put(path string, body any) ([]byte, error) {
	return c.request("PUT", path, nil, body)
}

func (c *Client) Delete(path string) ([]byte, error) {
	return c.request("DELETE", path, nil, nil)
}

func (c *Client) request(method, path string, query url.Values, body any) ([]byte, error) {
	data, status, err := c.do(method, path, query, body)
	if err != nil {
		return nil, err
	}
	if status == 204 {
		return []byte("{}"), nil
	}
	if status >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", status, string(data))
	}
	return data, nil
}

// PrintJSON pretty-prints JSON bytes to stdout.
func PrintJSON(data []byte) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Println(string(data))
		return
	}
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}

// Die prints an error and exits.
func Die(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
