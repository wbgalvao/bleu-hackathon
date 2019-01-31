package model

// Balance represents the balance of the account for a given currency
type Balance struct {
	Currency      string
	Balance       string
	Available     string
	Pending       string
	CryptoAddress string
	IsActive      string
}

// Market holds the information about a cryptocurrency market.
type Market struct {
	Currency     string
	BaseCurrency string
	PrevDay      float64
	High         float64
	Low          float64
	Last         float64
	Average      float64
	Volume       float64
	BaseVolume   float64
	TimeStamp    string
	Bid          float64
	Ask          float64
	IsActive     string
}

// Order represents an order in the cryptocurrency exchange server.
type Order struct {
	OrderID            string
	Exchange           string
	Type               string
	Quantity           string
	QuantityRemaining  string
	QuantityBaseTraded string
	Price              string
	Status             string
	Created            string
	Comments           string
}
