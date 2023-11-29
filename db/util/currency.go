package util

const (
	GHS = "GHS"
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case GHS, USD, EUR, GBP:
		return true
	}
	return false
}