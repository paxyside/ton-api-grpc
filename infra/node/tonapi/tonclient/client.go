package tonclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"emperror.dev/errors"
)

type Client struct {
	URL    string
	APIKey string
	client *http.Client
}

func NewClient(rpcURL, apiKey string, client *http.Client) *Client {
	return &Client{
		URL:    rpcURL,
		APIKey: apiKey,
		client: client,
	}
}

func (c *Client) GetRequest(ctx context.Context, link string) ([]byte, int, error) {
	reqURL, err := url.Parse(link)
	if err != nil {
		return nil, 0, errors.Wrap(err, "url.Parse")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "http.NewRequestWithContext")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "c.client.Do")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "io.ReadAll")
	}

	return body, resp.StatusCode, nil
}

func (c *Client) PostRequest(ctx context.Context, link string, body []byte) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, link, bytes.NewReader(body))
	if err != nil {
		return nil, 0, errors.Wrap(err, "http.NewRequestWithContext")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "client.Do")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "io.ReadAll")
	}

	return respBody, resp.StatusCode, nil
}

func (c *Client) HandleBusinessError(req []byte, code int) error {
	var errResp struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal(req, &errResp); err != nil {
		return errors.Errorf("json.Unmarshal: status: %d, raw: %s", code, string(req))
	}

	return errors.Errorf("error: %s", errResp.Error)
}

func (c *Client) HandleBusinessErrorWithMap(
	req []byte, code int, errorMap map[string]error,
) error {
	var errResp struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal(req, &errResp); err != nil {
		return errors.Errorf("json.Unmarshal: status: %d, raw: %s", code, string(req))
	}

	if mappedErr, ok := errorMap[errResp.Error]; ok {
		return mappedErr
	}

	return errors.Errorf("error: %s", errResp.Error)
}
