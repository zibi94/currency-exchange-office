package exchange

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	required := require.New(t)

	respRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(respRecorder)
	ctx.Request = &http.Request{
		URL: func() *url.URL {
			u, err := url.Parse("https://example.org/exchange?from=WBTC&to=USDT&amount=1.0")
			required.NoError(err, "parse url")
			return u
		}(),
	}

	Handler(ctx)

	required.Equal(http.StatusOK, respRecorder.Code, "status code equal")

	var resp Response
	err := json.NewDecoder(respRecorder.Body).Decode(&resp)
	required.NoError(err, "decode body error unexpected")

	required.Equal(57094.314314, resp.Amount, "rates equal")
}
