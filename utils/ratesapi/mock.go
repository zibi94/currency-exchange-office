package ratesapi

import "context"

type MockServiceClient struct {
	GetRatesOut func() (RatesLookup, error)
}

func (m *MockServiceClient) GetRates(ctx context.Context) (RatesLookup, error) {
	if m.GetRatesOut != nil {
		return m.GetRatesOut()
	}

	return RatesLookup{}, nil
}
