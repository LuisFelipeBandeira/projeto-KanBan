package model

import (
	"errors"
	"strings"
)

type User struct {
	Id       int    `json:"id"`
	Nome     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (usuario *User) Prepare() error {
	if errValidation := usuario.validation(); errValidation != nil {
		return errValidation
	}

	usuario.format()
	return nil
}

func (user *User) validation() error {
	if user.Nome == "" {
		return errors.New("o campo nome é obrigatório")
	}

	if user.Username == "" {
		return errors.New("o campo username é obrigatório")
	}

	if user.Password == "" {
		return errors.New("o campo password é obrigatório")
	}

	return nil
}

func (user *User) format() {
	user.Nome = strings.TrimSpace(user.Nome)
	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)
}
