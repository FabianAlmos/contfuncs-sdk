package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/FabianAlmos/contfuncs-sdk/contfuncs"
	"io"
	"net/http"
	"os"
)

type Client struct {
	gatewayURL string
	httpClient *http.Client
}

func NewClient() *Client {
	gatewayURL := os.Getenv("GATEWAY_URL")
	if gatewayURL == "" {
		gatewayURL = "http://127.0.0.1:8080"
	}
	return &Client{
		gatewayURL: gatewayURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) Invoke(ctx context.Context, input contfuncs.InvokeInput, opts ...contfuncs.InvokeOption) (*contfuncs.InvokeOutput, error) {
	var body io.Reader
	if input.Payload != nil {
		data, err := json.Marshal(input.Payload)
		if err != nil {
			return nil, fmt.Errorf("marshal payload, err: %w", err)
		}
		body = bytes.NewReader(data)
	}

	invokeOptions := contfuncs.NewInvokeOptions()
	for _, opt := range opts {
		opt(invokeOptions)
	}

	url := fmt.Sprintf("%s/run/%s", c.gatewayURL, input.FunctionName)
	if invokeOptions.Version != "" {
		url = fmt.Sprintf("%s?version=%s", url, invokeOptions.Version)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request, err: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("invoke function %s, err: %w", input.FunctionName, err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body, err: %w", err)
	}

	out := contfuncs.InvokeOutput{
		StatusCode: resp.StatusCode,
	}

	if resp.Header.Get("X-Function-Error") == "true" {
		out.Err = fmt.Errorf("%s", string(b))
		return &out, nil
	}

	out.Data = b
	return &out, nil
}
