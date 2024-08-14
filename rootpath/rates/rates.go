package rates

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"zibi94/currency-exchange-office/utils/ratesapi"
)

type Rate struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

type Response []Rate

func (d *deps) Handler(ctx *gin.Context) {
	currenciesQuery, ok := ctx.GetQuery("currencies")
	if !ok || currenciesQuery == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Deduplicate currencies.
	currencies := slices.Compact(strings.Split(currenciesQuery, ","))
	if len(currencies) <= 1 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	rates, err := d.ratesAPI.GetRates(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resp := Response{}
	for i := 0; i < len(currencies); i++ {
		for j := i + 1; j < len(currencies); j++ {
			iRate, err := rates.Get(currencies[i])
			if err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}

			jRate, err := rates.Get(currencies[j])
			if err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}

			resp = append(resp,
				Rate{
					From: currencies[i],
					To:   currencies[j],
					Rate: jRate / iRate,
				},
				Rate{
					From: currencies[j],
					To:   currencies[i],
					Rate: iRate / jRate,
				},
			)
		}
	}

	ctx.JSON(http.StatusOK, resp)
}

type deps struct {
	ratesAPI ratesapi.Interface
}

func New(ratesAPI ratesapi.Interface) *deps {
	return &deps{
		ratesAPI: ratesAPI,
	}
}
