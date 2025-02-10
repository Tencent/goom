package test

import "net/http"

type Client struct {
	c Doer
}

//go:noinline
func NewClient(options ...ClientOptionFunc) (*Client, error) {
	return &Client{}, nil
}

type ClientOptionFunc func(*Client) error

func SetHttpClient() ClientOptionFunc {
	httpClient := &http.Client{}
	return func(c *Client) error {
		if httpClient != nil {
			c.c = httpClient
		} else {
			c.c = http.DefaultClient
		}
		return nil
	}
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}
