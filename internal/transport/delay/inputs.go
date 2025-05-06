package delay

import "time"

type DelayInput struct {
	EmployeeID  int       `json:"-"` // Only for toMap method; not in request body
	Sum         float32   `json:"sum"`
	AccrualDate time.Time `json:"accrual_date"`
}

func (i *DelayInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["employee_id"] = i.EmployeeID
	m["sum"] = i.Sum
	m["accrual_date"] = i.AccrualDate
	return m
}
