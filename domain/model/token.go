package model

type AuthenticateTokenClaims struct {
	Uaid string `json:"uaid"`
}

type AuthorizeTokenClaims struct {
	UserAccountId string `json:"uaid"`
	RoleId        string `json:"rid"`
	ModuleId      string `json:"mid"`
	BehaviorId    string `json:"bid"`
}
