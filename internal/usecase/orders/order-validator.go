package orders

import (
	"context"
)

type OrderValidator interface {
	ValidateOrderFormat(ctx context.Context, orderNumber string) bool
	GetMaxRequestsPerMinute() uint64
	SetMaxRequestsPerMinute(max uint64)
}
