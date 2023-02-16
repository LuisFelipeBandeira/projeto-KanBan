package model

type Card struct {
	Id        int    `json:"cardid"`
	Desc      string `json:"description"`
	DateLimit string `json:"date"`
	HourLimit string `json:"hour"`
}
