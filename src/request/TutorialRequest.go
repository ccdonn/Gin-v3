package request

type TutorialRequest struct {
	PageSize int16 `form:"pageSize,default=10"`
	PageNo   int16 `form:"pageNo,default=1"`
}
