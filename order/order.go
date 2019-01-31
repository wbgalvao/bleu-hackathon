package order

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
