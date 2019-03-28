package request

type SuggestionRequest struct {
	AgentID   int32  `form:"agentId"`
	Type      int8   `form:"type"`
	Status    int8   `form:"status"`
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	PageNum   int16  `form:"pageNum,default=1"`
	PageSize  int16  `form:"pageSize,default=10"`
}
