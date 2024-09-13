package model

type Auth struct {
	UserAccountId string `json:"userAccountId"`
	RoleId        string `json:"roleId"`
	ModuleId      string `json:"moduleId"`
	BehaviorId    string `json:"behaviorId"`
	//permissions array de numeros
	Permissions string `json:"permissions"`
}
