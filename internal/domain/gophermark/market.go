package gophermark

type Market struct {
	Salt string
}

func NewMarket(salt string) *Market {
	return &Market{Salt: salt}
}
