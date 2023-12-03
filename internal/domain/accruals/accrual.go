package accruals

import "strings"

type Accrual struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

type AccrualSystemRegistered string

var (
	REGISTERED AccrualSystemRegistered = "REGISTERED"
	INVALID    AccrualSystemRegistered = "INVALID"
	PROCESSING AccrualSystemRegistered = "PROCESSING"
	PROCESSED  AccrualSystemRegistered = "PROCESSED"
)

func ValidateStatus(status string) bool {
	status = strings.ToTitle(status)
	return status == string(REGISTERED) || status == string(INVALID) ||
		status == string(PROCESSING) || status == string(PROCESSED)
}
