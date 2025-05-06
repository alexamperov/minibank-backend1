package deal

import "time"

type PayInput struct {
	Sum    float32 `json:"sum"`
	DealID int
}

func (i *PayInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["sum"] = i.Sum
	m["created_at"] = time.Now()
	m["deal_id"] = i.DealID
	return m
}
