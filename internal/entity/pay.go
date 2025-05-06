package entity

import "time"

type EPay struct {
	ID          int       `json:"id"`
	Status      string    `json:"status"`
	Amount      float64   `json:"amount"`
	Method      string    `json:"method"`
	DealID      int       `json:"deal_id"`
	PaymentDate time.Time `json:"payment_date"`
}
