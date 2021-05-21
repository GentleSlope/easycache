package auth_service

import "easycache/models"

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() (bool, models.Auth, error) {
	return models.CheckAuth(a.Username, a.Password)
}
