package model

type Student struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	University string `json:"university"`
	Major      string `json:"major"`
}
