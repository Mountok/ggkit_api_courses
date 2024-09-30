package processor

import (
	"errors"
	"fmt"
	"ggkit_learn_service/internals/app/db"
	"ggkit_learn_service/internals/app/models"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type LoginProcessor struct {
	storage *db.LoginStorage
}

func NewLoginProcessor(storage *db.LoginStorage) *LoginProcessor {
	process := new(LoginProcessor)
	process.storage = storage
	return process
}

func (processor *LoginProcessor) CreateUser(user models.User) (models.User, error) {
	currentTime := time.Now()

	user.CreateDate = fmt.Sprintf("%.2d.%.2d.%d-%.2d:%.2d", currentTime.Day(), currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute())
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashPass)
	_, err = processor.storage.CreateUser(user)
	if err != nil {
		return models.User{}, err
	}

	findUser, err := processor.storage.GetUserByEmail(user)
	if err != nil {
		return models.User{}, err
	}
	return findUser[0], nil
}

func (processor *LoginProcessor) Auth(user models.User) (models.User, error) {
	// vadition
	findUser, err := processor.storage.GetUserByEmail(user)
	if err != nil {
		return models.User{}, errors.New("такого пользователя не сущетвует")
	}
	err = bcrypt.CompareHashAndPassword([]byte(findUser[0].Password), []byte(user.Password))
	if err != nil {
		log.Println("Пароль неверный")
		return models.User{}, errors.New("пароль не верный")
	}
	log.Println("Пароль верный")
	log.Println(findUser)
	return findUser[0], nil
}

func (processor *LoginProcessor) Validate(uuid string) ([]models.ValidateUser, error) {
	if uuid == "" || uuid == " " {
		return []models.ValidateUser{}, errors.New("process: не валидный id")
	}
	return processor.storage.Validate(uuid)
}

func (processor *LoginProcessor) GetProfileByUserId(userid string) ([]models.Profile, error) {
	return processor.storage.GetProfileById(userid)
}
