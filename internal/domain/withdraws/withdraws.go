package withdraws

import "time"

type Withdraw struct {
	OrderID string
	Sum     float64
	Date    time.Time
}
