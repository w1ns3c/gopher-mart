package orders

import "time"

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
	Cashback  uint64
	Date      time.Time
}

type OrderStatus string

var (
	StatusNew        OrderStatus = "NEW"        // — заказ загружен в систему, но не попал в обработку;
	StatusProcessing OrderStatus = "PROCESSING" // — вознаграждение за заказ рассчитывается;
	StatusInvalid    OrderStatus = "INVALID"    // — система расчёта вознаграждений отказала в расчёте;
	StatusDone       OrderStatus = "PROCESSED"  // — данные по заказу проверены и информация о расчёте успешно получена.
)
