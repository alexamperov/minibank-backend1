package apply

import "time"

type ApplyInput struct {
	UserID   int       `json:"-"`
	Sum      float32   `json:"sum"`
	Percent  int       `json:"percent"`
	ReturnAt time.Time `json:"return_at"`
}

func (i *ApplyInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["user_id"] = i.UserID
	m["sum"] = i.Sum
	m["percent"] = i.Percent
	m["return_at"] = i.ReturnAt
	return m
}
