package balance

// Balance represents the balance of the account for a given currency
type Balance struct {
	Currency      string
	Balance       string
	Available     string
	Pending       string
	CryptoAddress string
	IsActive      string
}
