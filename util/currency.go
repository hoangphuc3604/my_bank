package util

const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
	JPY = "YEN"
	AUD = "AUD"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case USD, EUR, GBP, JPY, AUD:
		return true
	}
	return false
}