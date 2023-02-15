package model

type Card struct {
	Id    int    `json:"cardid"`
	Title string `json:"title"`
	Desc  string `json:"description"`
}
