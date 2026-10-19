package catalog

import (
	"fmt"
	"math"
	"strings"
)

type Currency string

const (
	RON Currency = "RON"
	EUR Currency = "EUR"
	HUF Currency = "HUF"
)

// fxRates are indicative rates, in units of the currency for 1 RON.
// They are fixed for the PoC.
var fxRates = map[Currency]float64{
	RON: 1,
	EUR: 0.2012,
	HUF: 79.45,
}

// fractionDigits is the display precision of each currency.
var fractionDigits = map[Currency]int{
	RON: 2,
	EUR: 2,
	HUF: 0,
}

// ParseCurrency accepts a currency code in any letter case.
func ParseCurrency(s string) (Currency, error) {
	c := Currency(strings.ToUpper(strings.TrimSpace(s)))
	if _, ok := fxRates[c]; !ok {
		return "", fmt.Errorf("unsupported currency %q", s)
	}
	return c, nil
}

// Convert turns an amount in RON bani into the target currency, rounded half
// away from zero to the currency's display precision.
func Convert(bani int64, to Currency) float64 {
	scale := math.Pow10(fractionDigits[to])
	return math.Round(float64(bani)/100*fxRates[to]*scale) / scale
}
