package user_service

import (
	"easycache/models"
	"easycache/pkg/logger"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

func GetRedisUserById(redisUserId int64) (user models.RedisUser, err error) {
	user, err = models.GetRedisUserById(redisUserId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetRedisUsers() (users []models.RedisUser, err error) {
	users, err = models.GetRedisUsers()
	if err != nil {
		return users, err
	}
	return users, nil
}

func InsertUser(user *models.RedisUser) (err error) {
	uuidStr := uuid.Must(uuid.NewV4(), err).String()
	if err != nil {
		return err
	}
	logger.Info("uuidStr: ", uuidStr)
	logger.Info("name:", user.Name)
	//todo 查看username是否重复
	userFind, err := models.FindRedisUserByName(user.Name)
	if err != nil {
		logger.Error("service_user InsertUser FindRedisUserByName err:", err)
		return err
	}
	if userFind.ID != 0 {
		errMSg := "user existed"
		logger.Error(errMSg)
		return errors.New(errMSg)
	}
	user.Uuid = uuidStr
	err = models.InsertRedisUser(user)
	if err != nil {
		logger.Error("service_user InsertRedisUser models.InsertRedisUser err:", err)
		return err
	}
	return nil
}

func DeleteUserById(id int64) (err error) {
	err = models.DeleteUserById(id)
	if err != nil {
		return err
	}
	return nil
}
