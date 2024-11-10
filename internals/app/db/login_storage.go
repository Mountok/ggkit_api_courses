package db

import (
	"context"
	"errors"
	"ggkit_learn_service/internals/app/models"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type LoginStorage struct {
	databasePool *pgxpool.Pool
}

func NewLoginStorage(databasePool *pgxpool.Pool) *LoginStorage {
	storage := new(LoginStorage)
	storage.databasePool = databasePool
	return storage
}

func (db *LoginStorage) CreateUser(user models.UserCreate) (string, error) {
	_, err := db.GetUserByEmail(user.Email)
	if err == nil {
		log.Errorf("Такой пользователь уже есть: \n %v", err)
		return "", errors.New("пользователь с таким email существует")
	}
	newUUID, err := uuid.NewV4()
	if err != nil {
		log.Errorf("Ошибка генерации uuid: \n %v", err)
		return "", err
	}
	query := "INSERT INTO users (id, email, password, create_date, role) values ($1,$2,$3,$4,'user');"
	_, err = db.databasePool.Exec(context.Background(), query, newUUID.String(), user.Email, user.Password, user.CreateDate)
	if err != nil {
		log.Errorf("Ошибка при sql запросе: \n %v", err)
		return "", err
	}
	log.Infof("Пользователь создан")
	user_id, err := db.CreateProfileForUser(user)
	if err != nil {
		log.Errorf("Ошибка при sql запросе: \n %v", err)
		return "", err
	}
	return user_id, nil
}

func (db *LoginStorage) GetUserByEmail(email string) (res []models.User, err error) {
	query := "SELECT id, email, password, role, create_date FROM users WHERE email = $1"
	err = pgxscan.Select(context.Background(), db.databasePool, &res, query, email)
	if err != nil {
		log.Errorf("Ошибка при sql запросе: \n %v", err)
		return res, err
	}
	if len(res) == 0 {
		log.Errorf("Ошибка при sql запросе: \n %v", err)
		return res, errors.New("пользователь не найдет")
	}
	log.Infof("Пользователь найден")
	return res, nil
}

func (db *LoginStorage) CreateProfileForUser(user models.UserCreate) (string, error) {
	query := "insert into profiles (user_id,description,phone,full_name, image) values ($1,$2,$3,$4,$5);"
	log.Infof("Получение пользователя по почте")
	currentUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		return "", err
	}
	log.Infof("Создание профиля для пользователя")
	_, err = db.databasePool.Exec(context.Background(), query, currentUser[0].Id, "нажмите что бы изменить описание", "-", user.Username, "admin.png")
	if err != nil {
		return "", err
	}
	log.Infof("Профиль создан")

	query = "INSERT INTO last_subjects(user_id,subjects_id) VALUES($1,$2)"
	_, err = db.databasePool.Exec(context.Background(),query,currentUser[0].Id,1)
	if err != nil{
		return currentUser[0].Id, err
	}



	return currentUser[0].Id, nil
}

// фунуция для получения профиля по id пользователя
func (db *LoginStorage) GetProfileById(user_id string) (profile []models.Profile, err error) {
	query := "SELECT id,user_id,description,phone,full_name,image,score FROM profiles WHERE user_id = $1"
	err = pgxscan.Select(context.Background(), db.databasePool, &profile, query, user_id)
	if err != nil {
		return []models.Profile{}, err
	}
	return profile, err

}

func (db *LoginStorage) Validate(uuid string) (validateUser []models.ValidateUser, err error) {
	query := "SELECT id, role FROM users WHERE id = $1"
	err = pgxscan.Select(context.Background(), db.databasePool, &validateUser, query, uuid)
	if err != nil {
		log.Errorf("Ошибка sql pf")
		return validateUser, err
	}
	if len(validateUser) == 0 {
		return validateUser, errors.New("storage: не валидный id")
	}
	return validateUser, nil
}

// id serial primary key,
// 	user_id integer not null,
// 	description varchar(125) not null,
// 	phone varchar(100) not null,
// 	full_name varchar(125) not null,
// 	image text not null,
// 	score integer not null default 0,
