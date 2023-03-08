package model

type Card struct {
	Id        int    `json:"cardid"`
	Company   string `json:"company"`
	Desc      string `json:"description"`
	DateLimit string `json:"date"`
	HourLimit string `json:"hour"`
}
