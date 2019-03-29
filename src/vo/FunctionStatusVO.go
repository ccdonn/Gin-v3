package vo

// FunctionStatusVO : type define
type FunctionStatusVO struct {
	Settlement         bool `json:"settlement"`
	Maintenance        bool `json:"maintenance"`
	AlipayMaintenance  bool `json:"alipayMaintenance"`
	BankMaintenance    bool `json:"bankMaintenance"`
	CreateRegisterGate bool `json:"createRegisterGate"`
}
