package shipping

import (
	"strings"
)

// RateQuote is a carrier/service option (live APIs can replace stub math when keys are configured).
type RateQuote struct {
	Carrier   string  `json:"carrier"`
	Service   string  `json:"service"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	DaysMin   int     `json:"daysMin"`
	DaysMax   int     `json:"daysMax"`
	Code      string  `json:"code"` // machine id: fedex_ground, ups_ground, dhl_express
}

// QuoteCarriers returns deterministic stub quotes for FedEx / UPS / DHL based on weight and destination country.
// Replace with EasyPost / Shippo / carrier REST SDKs in production for live rating.
func QuoteCarriers(destCountry string, totalWeightKg float64) []RateQuote {
	dest := strings.ToUpper(strings.TrimSpace(destCountry))
	if dest == "" {
		dest = "US"
	}
	w := totalWeightKg
	if w <= 0 {
		w = 0.5
	}
	base := 6.0 + w*4.2
	if dest != "US" {
		base += 18 + w*9
	}
	return []RateQuote{
		{Carrier: "FedEx", Service: "Ground", Amount: round2(base * 1.05), Currency: "USD", DaysMin: 3, DaysMax: 6, Code: "fedex_ground"},
		{Carrier: "UPS", Service: "Ground", Amount: round2(base * 1.02), Currency: "USD", DaysMin: 3, DaysMax: 7, Code: "ups_ground"},
		{Carrier: "DHL", Service: "Express Worldwide", Amount: round2(base * 1.35 + 6), Currency: "USD", DaysMin: 2, DaysMax: 5, Code: "dhl_express"},
	}
}

func round2(x float64) float64 {
	return float64(int64(x*100+0.5)) / 100
}

func FindQuote(code string, destCountry string, weight float64) (RateQuote, bool) {
	code = strings.ToLower(strings.TrimSpace(code))
	for _, q := range QuoteCarriers(destCountry, weight) {
		if strings.EqualFold(q.Code, code) {
			return q, true
		}
	}
	return RateQuote{}, false
}
