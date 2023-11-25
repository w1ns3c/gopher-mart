package withdraws

import "time"

type Withdraw struct {
	OrderID string
	Sum     uint64
	Date    time.Time
}
