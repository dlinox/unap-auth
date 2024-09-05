package model

type Role struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Roles struct {
	Roles []Role `json:"roles"`
}
