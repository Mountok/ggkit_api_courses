package processor

import (
	"errors"
	"fmt"
	"ggkit_learn_service/internals/app/db"
	"ggkit_learn_service/internals/app/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
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
const (
	signingString = "39jgr0gtir9mc289g"
	tokenTTL      = 1 * time.Minute
	salt          = "&5a2@4D$5~54dC*&"
)






func (processor *LoginProcessor) GenerateToken(email string, password string) (string, error) {
	user, err := processor.storage.GetUser(email,password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
		UserRole: user.Role,
	})
	return token.SignedString([]byte(signingString))

}

func ParseToken(token string) (id, role string,err error) {
	accessToken, err := jwt.ParseWithClaims(token, &models.TokenClaims{},func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid singing method")
		}
		return []byte(signingString), nil
	})
	if err != nil {
		return "","", err
	}
	claims, ok := accessToken.Claims.(*models.TokenClaims)
	if !ok {
		return "","", err
	}
	return claims.UserId, claims.UserRole, nil
}


func (processor *LoginProcessor) CreateUser(user models.UserCreate) (models.User, error) {
	currentTime := time.Now()
	newUUID, err := uuid.NewV4()
	if err != nil {
		return models.User{},err
	}
	user.Id = newUUID.String()
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

	findUser, err := processor.storage.GetUserByEmail(user.Email)
	if err != nil {
		return models.User{}, err
	}
	return findUser[0], nil
}

// func (processor *LoginProcessor) CreateUser(user models.UserCreate) (models.User, error) {
// 	currentTime := time.Now()

// 	user.CreateDate = fmt.Sprintf("%.2d.%.2d.%d-%.2d:%.2d", currentTime.Day(), currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute())
// 	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return models.User{}, err
// 	}
// 	user.Password = string(hashPass)
// 	_, err = processor.storage.CreateUser(user)
// 	if err != nil {
// 		return models.User{}, err
// 	}

// 	findUser, err := processor.storage.GetUserByEmail(user.Email)
// 	if err != nil {
// 		return models.User{}, err
// 	}
// 	return findUser[0], nil
// }

func (processor *LoginProcessor) Auth(user models.User) (models.User, error) {
	// vadition
	findUser, err := processor.storage.GetUserByEmail(user.Email)
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
