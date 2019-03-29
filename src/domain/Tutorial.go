package domain

import "time"

type Tutorial struct {
	ID                  int32  `json:"id"`
	Title               string `json:"title"`
	TitleImg            string `json:"titleImg"`
	Content             string `json:"content"`
	createTime          *time.Time
	CreateTimeValue     int64 `json:"createTime"`
	Del                 int8  `json:"del"`
	LastUpdateUser      int32 `json:"lastUpdateUser"`
	lastUpdateTime      *time.Time
	LastUpdateTimeValue int64 `json:"lastUpdateTime"`
	// TopAndOrder    int
}
