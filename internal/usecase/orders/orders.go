package orders

type ID string
type Item struct {
	ItemID   string
	ItemCost uint64
}

type Position struct {
	Item  Item
	Count uint64
}

type Order struct {
	ID        ID
	Sum       float64
	Positions []Position
}
