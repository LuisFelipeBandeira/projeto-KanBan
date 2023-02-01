package model

type User struct {
	Id         int    `json:"id"`
	Nome       string `json:"name"`
	Email      string `json:"email"`
	Senha      string `json:"password"`
	Permission int    `json:"permission"`
}
