package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"ggkit_learn_service/internals/app/models"
	"ggkit_learn_service/internals/app/processor"
	"net/http"
	"strings"
)

type LoginHandler struct {
	process *processor.LoginProcessor
}

func NewLoginhandler(processor *processor.LoginProcessor) *LoginHandler {
	handler := new(LoginHandler)
	handler.process = processor
	return handler
}

func (handler *LoginHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	var input models.UserCreate
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		WrapError(w, err)
		return
	}

	userId, err := handler.process.CreateUser(input)
	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"id": userId.Id,
	}
	WrapOK(w, m)

}
func (handler *LoginHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input models.UserSignIn
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		WrapError(w, err)
		return
	}
	fmt.Println(input)
	token, err := handler.process.GenerateToken(input.Email, input.Password)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}

	var m = map[string]interface{}{
		"token": token,
	}
	WrapOK(w, m)
}

func (handler *LoginHandler) Authorization(w http.ResponseWriter, r *http.Request) {
	w, r, err := UserIdentify(w,r)
	if err != nil {
		WrapErrorWithStatus(w,err,http.StatusUnauthorized)
		return
	}
	var m = map[string]interface{}{
		"user_id": w.Header().Get(UserCtx),
		"user_role": w.Header().Get(UserRole),
	}
	WrapOK(w,m)
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
		"result":    "ok",
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
		"result":    "OK",
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
		WrapError(w, err)
		return
	}

	validateUser, err := handler.process.Validate(userUUID.Id)
	if err != nil {
		WrapError(w, err)
		return
	}
	if len(validateUser) == 1 {
		var m = map[string]interface{}{
			"result":   "OK",
			"is_valid": true,
			"user":     validateUser,
		}
		WrapOK(w, m)
		return
	} else {
		WrapErrorWithStatus(w, errors.New("не авторизированный пользователь"), http.StatusUnauthorized)
	}

}

func (handler *LoginHandler) Profile(w http.ResponseWriter, r *http.Request) {
	w,r,err := UserIdentify(w,r)
	if err != nil {
		WrapErrorWithStatus(w,err,http.StatusUnauthorized)
	}
	userID := w.Header().Get(UserCtx)

	data, err := handler.process.GetProfileByUserId(userID)
	if err != nil {
		WrapError(w, err)
		return
	}

	m := map[string]interface{}{
		"result": "ok",
		"data":   data,
	}

	WrapOK(w, m)
}

const (
	authorizationHeader = "Authorization"
	UserCtx             = "userId"
	UserRole 			= "userRole"
)

func UserIdentify(w http.ResponseWriter, r *http.Request) (rw http.ResponseWriter, rs *http.Request, err error) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff") //даем понять что ответ приходит в формате json
			w.WriteHeader(http.StatusUnauthorized)
			return w,r,errors.New("empty header")
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff") //даем понять что ответ приходит в формате json
			w.WriteHeader(http.StatusUnauthorized)
			return w,r,errors.New("invalid header")
		}

		userId,userRole, err := processor.ParseToken(string(headerParts[1]))
		if err != nil {

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff") //даем понять что ответ приходит в формате json
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err)
			return w,r,errors.New("invalid token")
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set(UserCtx, userId)
		w.Header().Set(UserRole, userRole)
		return w,r,nil
}
