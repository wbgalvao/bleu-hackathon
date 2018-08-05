package market

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
