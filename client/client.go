package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	gatewayURL string
	httpClient *http.Client
}

func NewClient(gatewayURL string) *Client {
	return &Client{
		gatewayURL: gatewayURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) Invoke(ctx context.Context, fnName string, in any, out any) error {
	var body io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("marshal payload, err: %w", err)
		}
		body = bytes.NewReader(data)
	}

	url := fmt.Sprintf("%s/run/%s", c.gatewayURL, fnName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("create request, err: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("invoke function %s, err: %w", fnName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("function %s returned error %d: %s", fnName, resp.StatusCode, b)
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decode response, err: %w", err)
		}
	}

	return nil
}
