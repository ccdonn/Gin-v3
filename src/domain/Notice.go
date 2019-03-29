package domain

import "time"

// Notice : type define
type Notice struct {
	ID              int32  `json:"id"`
	Title           string `json:"title"`
	Type            int8   `json:"type"`
	CreateTimeValue int64  `json:"createTime"`
	FromTimeValue   int64  `json:"fromTime"`
	EndTimeValue    int64  `json:"toTime"`
	AgentID         int32  `json:"agentId"`
	ReadState       int8   `json:"readState"`

	createTime *time.Time
	FromTime   *time.Time
	EndTime    *time.Time
}
