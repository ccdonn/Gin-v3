package domain

import "time"

// Tip : type define
type Tip struct {
	ID             int32  `json:"id"`
	Name           string `json:"name"`
	Content        string `json:"content"`
	status         int8
	createUser     int32
	createTime     *time.Time
	lastUpdateUser int32
	lastUpdateTime *time.Time
}
