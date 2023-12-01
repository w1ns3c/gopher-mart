package orders

import (
	"strings"
	"time"
)

type Item struct {
	ItemID   string
	ItemCost uint64
}

type Position struct {
	Item  Item
	Count uint64
}

type Order struct {
	ID        string
	Sum       float64
	Positions []Position
	Cashback  uint64
	Date      time.Time
	Status    OrderStatus
}

type OrderStatus string

var (
	StatusNew        OrderStatus = "NEW"        // — заказ загружен в систему, но не попал в обработку;
	StatusProcessing OrderStatus = "PROCESSING" // — вознаграждение за заказ рассчитывается;
	StatusInvalid    OrderStatus = "INVALID"    // — система расчёта вознаграждений отказала в расчёте;
	StatusDone       OrderStatus = "PROCESSED"  // — данные по заказу проверены и информация о расчёте успешно получена.
)

type AccrualSystemRegistered string

var (
	REGISTERED AccrualSystemRegistered = "REGISTERED"
	INVALID    AccrualSystemRegistered = "INVALID"
	PROCESSING AccrualSystemRegistered = "PROCESSING"
	PROCESSED  AccrualSystemRegistered = "PROCESSED"
)

func ValidateStatus(status string) bool {
	status = strings.ToTitle(status)
	return status == string(StatusNew) || status == string(StatusProcessing) ||
		status == string(StatusInvalid) || status == string(StatusDone)
}
