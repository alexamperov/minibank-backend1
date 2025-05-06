package entity

import "time"

type EApply struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	EmployeeID int       `json:"employee_id"`
	Sum        float32   `json:"sum"`
	Percent    int       `json:"percent"`
	ReturnAt   time.Time `json:"return_at"`
}
