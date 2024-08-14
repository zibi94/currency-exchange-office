package rates

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"zibi94/currency-exchange-office/utils/ratesapi"
)

func TestHandler(t *testing.T) {
	d := deps{
		ratesAPI: &ratesapi.MockServiceClient{
			GetRatesOut: func() (ratesapi.RatesLookup, error) {
				return ratesapi.RatesLookup{
					"ETB": 105.125,
					"EUR": 0.915435,
					"GBP": 0.785423,
					"PLN": 3.958682,
					"USD": 1,
				}, nil
			},
		},
	}

	testCases := []struct {
		name, query string
		expected    []Rate
	}{
		{
			name: "three currencies", query: "USD,GBP,EUR",
			expected: []Rate{
				{From: "USD", To: "GBP", Rate: 0.785423},
				{From: "GBP", To: "USD", Rate: 1.2731992824248846},
				{From: "USD", To: "EUR", Rate: 0.915435},
				{From: "EUR", To: "USD", Rate: 1.09237684816508},
				{From: "GBP", To: "EUR", Rate: 1.1655311851066241},
				{From: "EUR", To: "GBP", Rate: 0.8579779012163616},
			},
		},
		{
			name: "two currencies", query: "USD,PLN",
			expected: []Rate{
				{From: "USD", To: "PLN", Rate: 3.958682},
				{From: "PLN", To: "USD", Rate: 0.25260932805413516},
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].query, func(t *testing.T) {
			required := require.New(t)

			respRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(respRecorder)
			ctx.Request = &http.Request{
				URL: func() *url.URL {
					u, err := url.Parse("https://example.org?currencies=" + testCases[i].query)
					required.NoError(err, "parse url")
					return u
				}(),
			}

			d.Handler(ctx)

			required.Equal(http.StatusOK, respRecorder.Code, "status code equal")

			var rates []Rate
			err := json.NewDecoder(respRecorder.Body).Decode(&rates)
			required.NoError(err, "decode body error unexpected")

			required.Equal(testCases[i].expected, rates, "rates equal")
		})
	}
}

func TestFailureHandler(t *testing.T) {
	d := deps{
		ratesAPI: &ratesapi.MockServiceClient{
			GetRatesOut: func() (ratesapi.RatesLookup, error) {
				return ratesapi.RatesLookup{
					"PLN": 3.958682,
					"USD": 1,
				}, nil
			},
		},
	}

	testCases := []struct {
		name, query  string
		expectedCode int
	}{
		{

			name: "two currencies", query: "USD,PLN",
			expectedCode: http.StatusOK,
		},
		{
			name: "one currency", query: "USD",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "one duplicated currency", query: "USD,USD",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "no currenc", query: "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for i := range testCases {
		t.Run(testCases[i].query, func(t *testing.T) {
			required := require.New(t)

			respRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(respRecorder)
			ctx.Request = &http.Request{
				URL: func() *url.URL {
					u, err := url.Parse("https://example.org?currencies=" + testCases[i].query)
					required.NoError(err, "parse url")
					return u
				}(),
			}

			d.Handler(ctx)

			required.Equal(testCases[i].expectedCode, respRecorder.Code, "status code equal")
		})
	}
}
