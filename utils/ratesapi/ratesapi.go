package ratesapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Interface interface {
	GetRates(ctx context.Context) (RatesLookup, error)
}

var _ Interface = (*serviceClient)(nil)

type serviceClient struct {
	appID  string
	client interface {
		Do(req *http.Request) (*http.Response, error)
	}
}

func NewClient(appID string) *serviceClient {
	return &serviceClient{
		appID:  appID,
		client: http.DefaultClient,
	}
}

var ErrUnexpectedStatusCode = errors.New("unexpected status code")

func (svc *serviceClient) GetRates(ctx context.Context) (RatesLookup, error) {
	queryVal := make(url.Values)
	queryVal.Add("app_id", svc.appID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://openexchangerates.org/api/latest.json?"+queryVal.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request with context: %w", err)
	}

	req.Header.Add("Accept", "application/json")

	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d: %w", resp.StatusCode, ErrUnexpectedStatusCode)
	}

	var out struct {
		Rates RatesLookup `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode resp body: %w", err)
	}

	return out.Rates, nil
}
