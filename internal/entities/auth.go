package entities

import "fmt"

var (
	ErrWrongLoginOrPassword = fmt.Errorf("wrong login or password")
)

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Token struct {
	Access string `json:"access_token"`
	RT     string `json:"refresh_token"`
}
