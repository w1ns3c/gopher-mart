package orders

import (
	"time"
)

type Order struct {
	ID       string
	Sum      float64
	Cashback float64
	Date     time.Time
	Status   OrderStatus
}

type OrderStatus string

var (
	StatusNew        OrderStatus = "NEW"        // — заказ загружен в систему, но не попал в обработку;
	StatusProcessing OrderStatus = "PROCESSING" // — вознаграждение за заказ рассчитывается;
	StatusInvalid    OrderStatus = "INVALID"    // — система расчёта вознаграждений отказала в расчёте;
	StatusDone       OrderStatus = "PROCESSED"  // — данные по заказу проверены и информация о расчёте успешно получена.
)
