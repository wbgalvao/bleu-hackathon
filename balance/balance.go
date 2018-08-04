package balance

// Balance represents the balance of the account for a given currency
type Balance struct {
	Currency      string
	Balance       float32
	Available     float32
	Pending       float32
	CryptoAddress string
	IsActive      bool
}
