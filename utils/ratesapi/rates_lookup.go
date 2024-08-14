package ratesapi

import (
	"errors"
	"fmt"
)

var ErrCurrencyNotFound = errors.New("currency not found")

type RatesLookup map[string]float64

func (rl RatesLookup) Get(currency string) (float64, error) {
	r, ok := rl[currency]
	if !ok {
		return 0, fmt.Errorf("%s: %w", currency, ErrCurrencyNotFound)
	}

	return r, nil
}
