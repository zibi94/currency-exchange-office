package ratesapi

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRates(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		required := require.New(t)
		svc := serviceClient{
			appID:  "my_app_id",
			client: new(mockHTTP),
		}

		rates, err := svc.GetRates(context.Background())
		required.NoError(err, "get rates error unexpected")
		required.Equal(4, len(rates))
		for _, curr := range []string{"EUR", "GBP", "PLN", "USD"} {
			_, err := rates.Get(curr)
			required.NoError(err, "%s get error unexpected")
		}
	})

	t.Run("missing app id", func(t *testing.T) {
		required := require.New(t)
		svc := serviceClient{
			appID:  "",
			client: new(mockHTTP),
		}

		_, err := svc.GetRates(context.Background())
		required.Error(err, "get rates error expected")
		required.ErrorContains(err, "403", "error contains")
	})
}

type mockHTTP struct{}

func (m *mockHTTP) Do(req *http.Request) (*http.Response, error) {
	if req.URL.Query().Get("app_id") == "" {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       io.NopCloser(strings.NewReader(`{"error":true,"status":403,"message":"missing_app_id"}`)),
		}, nil
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(
			strings.NewReader(`{"timestamp":1723327211,"base":"USD","rates":{"EUR":0.91529,"GBP":0.785484,"PLN":3.956881,"USD":1}}`),
		),
	}, nil
}
