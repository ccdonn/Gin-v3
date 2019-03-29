package vo

// LoginLogVO : type define
type LoginLogVO struct {
	AgentID int32  `json:"agentId"`
	IP      string `json:"ip"`
	IPArea  string `json:"ipArea"`
	Date    string `json:"date"`
}
