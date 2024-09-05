package model

type Module struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Route       string `json:"route"`
	Icon        string `json:"icon"`
}

type Modules struct {
	Modules []Module `json:"modules"`
}
