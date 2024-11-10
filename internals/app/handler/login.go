package handler

import (
	"encoding/json"
	"errors"
	"ggkit_learn_service/internals/app/models"
	"ggkit_learn_service/internals/app/processor"
	"net/http"

	"github.com/gorilla/mux"
)

type LoginHandler struct {
	process *processor.LoginProcessor
}

func NewLoginhandler(processor *processor.LoginProcessor) *LoginHandler {
	handler := new(LoginHandler)
	handler.process = processor
	return handler
}

func (handler *LoginHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user models.UserCreate
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WrapError(w, err)
	}

	userInfo, err := handler.process.CreateUser(user)
	if err != nil {
		WrapError(w, err)
	}
	var m = map[string]interface{}{
		"result": "ok",
		"user_info": userInfo,

	}
	WrapOK(w, m)
}

func (handler *LoginHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WrapError(w, err)
		return
	}

	userInfo, err := handler.process.Auth(user)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"user_info": userInfo,
	}
	WrapOK(w, m)
}


type UserUUID struct {
	Id string `json:"user_uuid"`
}

func (handler *LoginHandler) Validate(w http.ResponseWriter, r *http.Request) {
	var userUUID UserUUID
	err := json.NewDecoder(r.Body).Decode(&userUUID)
	if err != nil {
		WrapError(w,err)
		return
	}

	validateUser, err := handler.process.Validate(userUUID.Id)
	if err != nil {
		WrapError(w, err)
		return
	}
	if len(validateUser) == 1 {
		var m = map[string]interface{}{
			"result": "OK",
			"is_valid": true,
			"user": validateUser,
		}
		WrapOK(w, m)
		return
	} else {
		WrapErrorWithStatus(w,errors.New("не авторизированный пользователь"),http.StatusUnauthorized)
	}
	
}



func (handler *LoginHandler) Profile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	data, err := handler.process.GetProfileByUserId(userID)
	if err != nil {
		WrapError(w,err)
		return
	}

	m := map[string]interface{}{
		"result": "ok",
		"data": data,
	}

	WrapOK(w,m)
}