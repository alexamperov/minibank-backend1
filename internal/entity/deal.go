package entity

import "time"

type EDeal struct {
	ID       int       `json:"ID"`
	Status   string    `json:"status"`
	Sum      float32   `json:"sum"`
	Currency string    `json:"currency"`
	Percent  int       `json:"percent"`
	IssuedAt time.Time `json:"issued_at"`
	ReturnAt time.Time `json:"return_at"`

	UserID     int `json:"user_id"`     // TODO in response for employee and admin
	EmployeeID int `json:"employee_id"` // TODO in response for user and admin
}
