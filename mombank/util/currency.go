package util

const (
	USD = "USD"
	KWR = "KWR"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, KWR:
		return true
	}
	return false
}
