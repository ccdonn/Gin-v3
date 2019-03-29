package request

type TutorialRequest struct {
	PageSize int16  `form:"pageSize,default=10"`
	PageNum  int16  `form:"pageNum,default=1"`
	Query    string `form:"q"`
	Order    string `form:"o,default=createTime"`
	Desc     string `form:"d,default=desc"`
}
