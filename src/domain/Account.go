package domain

type Account struct {
	AgentID            int32  `json:"agentId"`
	Nickname           string `json:"nickname"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	Phone              string `json:"phone"`
	SettlementPassword string `json:"settlementPassword"`
}
