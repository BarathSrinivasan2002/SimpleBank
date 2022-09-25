package util

//const for all supported currency

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

//IsSupportedCurrency returns true if the input currencyis supported

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
