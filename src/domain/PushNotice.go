package domain

import "time"

// PushNotice : type define
type PushNotice struct {
	ID            int32  `json:"id"`
	Title         string `json:"title"`
	PushTimeValue int64  `json:"pushTime"`
	Content       string `json:"content"`
	CreateDate    int64  `json:"createDate"`
	NoticeID      int32  `json:"noticeId"`
	Extra         string `json:"extra"`

	pushTime   *time.Time
	createDate *time.Time
}
