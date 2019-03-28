package response

type MeResponse struct {
	AgentID  int32  `json:"agentId"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}
