package auth_service

import (
	"easycache/models"
	"easycache/pkg/logger"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() (bool, models.Auth, error) {
	return models.CheckAuth(a.Username, a.Password)
}

func Login(username string, password string) (flag bool, err error) {
	pwd := []byte(password)
	pwdDatabase, err := models.GetPassword(username)
	if err != nil || pwdDatabase == "" {
		logger.Error("GetPassword err", err)
		return flag, errors.New("用户不存在")
	}
	err = bcrypt.CompareHashAndPassword([]byte(pwdDatabase), pwd)
	if err != nil {
		logger.Error("Check err", err)
		return flag, errors.New("密码错误")
	}
	flag = true
	return flag, err
}

func GetUser(username string) (user models.Auth, err error) {
	user, err = models.GetUser(username)
	if err != nil {
		logger.Error("Check err", err)
		return user, errors.New("密码错误")
	}
	return user, err
}

func Register(username string, password string) (flag bool, err error) {
	pwd := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		logger.Error("bcrypt err", err)
		return flag, err
	}
	pwdDatabase, err := models.GetPassword(username)
	if err != nil {
		logger.Error("GetPassword err", err)
		return flag, err
	}
	if pwdDatabase != "" {
		return flag, errors.New("该用户已注册")
	}
	err = models.Register(username, string(hashedPassword))
	if err != nil {
		logger.Error("Register err", err)
		return flag, err
	}
	flag = true
	return flag, err
}
