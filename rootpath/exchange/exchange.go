package exchange

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"zibi94/currency-exchange-office/utils/num"
)

const (
	from   = "from"
	to     = "to"
	amount = "amount"
)

type Response struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func Handler(ctx *gin.Context) {
	var params = map[string]struct{ value string }{
		from:   {},
		to:     {},
		amount: {},
	}

	// Check if each required query parameter exists.
	// and assign value of param.
	for key := range params {
		query, ok := ctx.GetQuery(key)
		if !ok || query == "" {
			ctx.AbortWithStatus(400)
			return
		}
		params[key] = struct{ value string }{value: query}
	}

	f, ok := base[params[from].value]
	if !ok {
		ctx.AbortWithStatus(400)
		return
	}

	t, ok := base[params[to].value]
	if !ok {
		ctx.AbortWithStatus(400)
		return
	}

	a, err := strconv.ParseFloat(params[amount].value, 64)
	if err != nil {
		ctx.AbortWithStatus(400)
		return
	}

	ctx.JSON(200, &Response{
		From:   f.CryptoCurrency,
		To:     t.CryptoCurrency,
		Amount: num.Round(a*f.Rate/t.Rate, t.DecimalPlaces),
	})
}

var base = map[string]struct {
	CryptoCurrency string
	DecimalPlaces  int
	Rate           float64
}{
	"BEER":  {CryptoCurrency: "BEER", DecimalPlaces: 18, Rate: 0.00002461},
	"FLOKI": {CryptoCurrency: "FLOKI", DecimalPlaces: 18, Rate: 0.0001428},
	"GATE":  {CryptoCurrency: "GATE", DecimalPlaces: 18, Rate: 6.87},
	"USDT":  {CryptoCurrency: "USDT", DecimalPlaces: 6, Rate: 0.999},
	"WBTC":  {CryptoCurrency: "WBTC", DecimalPlaces: 8, Rate: 57037.22},
}
