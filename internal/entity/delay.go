package entity

import "time"

type EDelay struct {
	ID          int       `json:"id"`
	Status      string    `json:"status"`
	Amount      float64   `json:"amount"`
	AccrualDate time.Time `json:"accrual_date"`
	DealID      int       `json:"deal_id"`
	EmployeeID  int       `json:"employee_id"` // Employee who created Delay for User
}
