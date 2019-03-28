package domain

import "time"

type Suggestion struct {
	ID              int32  `json:"id"`
	AgentID         int32  `json:"agentId"`
	Nickname        string `json:"nickname"`
	Username        string `json:"username"`
	Type            int32  `json:"type"`
	Content         string `json:"content"`
	createTime      *time.Time
	CreateTimeValue int64  `json:"createTime"`
	DeviceInfo      string `json:"deviceInfo"`
	ReplyContent    string `json:"replyContent"`
	replyTime       *time.Time
	ReplyTimeValue  int64 `json:"replyTime"`
	Status          int32 `json:"status"`
	// Image           string `json:"image"`
}
