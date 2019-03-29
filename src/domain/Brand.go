package domain

type Brand struct {
	ID     int32  `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status int8   `json:"status"`
	IsDel  int8
}
