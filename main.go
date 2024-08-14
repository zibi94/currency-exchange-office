package main

import (
	"os"
	"zibi94/currency-exchange-office/rootpath/exchange"
	"zibi94/currency-exchange-office/rootpath/rates"
	"zibi94/currency-exchange-office/utils/ratesapi"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/rates", rates.New(ratesapi.NewClient(os.Getenv("APP_ID"))).Handler)
	r.GET("/exchange", exchange.Handler)

	r.Run()
}
