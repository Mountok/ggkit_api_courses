package db

import (
	"context"
	"ggkit_learn_service/internals/app/models"
	"github.com/georgysavva/scany/pgxscan"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type UserStorage struct {
	databasePool *pgxpool.Pool
}

func NewUserStorage(databasePool *pgxpool.Pool) *UserStorage {
	storage := new(UserStorage)
	storage.databasePool = databasePool
	return storage
}

func (db *UserStorage) LastSubject(userId string) ([]models.LastSubject, error) {
	var lastCourses []models.LastSubject
	query := "select id, user_id, subjects_id from last_subjects where user_id = $1"
	logrus.Infof("Заппрос в бд для получения последнего курса для пользователя id=%s", userId)
	err := pgxscan.Select(context.Background(), db.databasePool, &lastCourses, query, userId)
	if err != nil {
		logrus.Errorf("Ошибка при запросе последнего курса для пользователя по id: %s \n содержание ошибки: %v", userId, err)
		return nil, err
	}
	logrus.Infof("Последний курс для пользователя id=%s получен.", userId)
	logrus.Print(lastCourses)
	return lastCourses, nil
}

func (db *UserStorage) SetLastSubject(userId, subjectId string) error {
	query := "update last_subjects set subjects_id=$1 where user_id=$2"
	_, err := db.databasePool.Exec(context.Background(), query, subjectId, userId)
	return err
}

func (db *UserStorage) UploadAvatar(userId, image_url string) error {
	query := "update profiles set image = $1 where user_id = $2;"
	_, err := db.databasePool.Exec(context.Background(), query, image_url, userId)
	logrus.Infof("обновление данных в бд: в процессе")
	if err != nil {
		return err
	}
	logrus.Infof("обновление данных в бд: успешно")

	return nil
}

func (db *UserStorage) ChangeUserName(userId, newName string) (userName string, err error) {
	query := "update profiles set full_name = $1 where user_id = $2"
	_, err = db.databasePool.Exec(context.Background(), query, newName, userId)
	logrus.Infof("обновление данных в бд")
	if err != nil {
		logrus.Errorf("ошибка при обновлении: \n %v", err)
		return "", err
	}
	logrus.Infof("все прошло успешло")
	return newName, nil
}

func (db *UserStorage) ChangeDescription(userId, newDescription string) error {
	var query string = "update profiles set description = $1 where user_id = $2"
	_, err := db.databasePool.Exec(context.Background(), query, newDescription, userId)
	if err != nil {
		return err
	}
	return nil
}

func (db *UserStorage) GetPoint(user_id, theme_id string) error {
	// Открытие транзакции
	tx, err := db.databasePool.Begin(context.Background())
	if err != nil {
		log.Fatalf("Unable to begin transaction: %v", err)
		return err // Возвращаем ошибку
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
			log.Printf("Transaction rolled back: %v", err)
		} else {
			err = tx.Commit(context.Background())
			if err != nil {
				log.Fatalf("Unable to commit transaction: %v", err)
			}
		}
	}()

	// Выполнение первого запроса
	_, err = tx.Exec(context.Background(), "INSERT INTO done_lessons (theme_id, user_id) VALUES ($1, $2)", theme_id, user_id)
	if err != nil {
		return err // Ошибка будет обработана в deferred функции
	}

	// Выполнение второго запроса
	_, err = tx.Exec(context.Background(), "UPDATE profiles SET score = score + 1 WHERE user_id = $1", user_id)
	if err != nil {
		return err
	}

	return nil
}

func (db *UserStorage) Rating() ([]models.RatingUser, error) {
	var lists []models.RatingUser
	var query string = "SELECT p.user_id, p.full_name, p.image, p.score FROM users u JOIN profiles p ON u.id = p.user_id ORDER BY p.score DESC;"
	err := pgxscan.Select(context.Background(), db.databasePool, &lists, query)
	return lists, err
}

func (db *UserStorage) CheckDoneLessons(user_id, theme_id string) (int, error) {
	var record []int
	var query string = "SELECT COUNT(*) AS record_exists FROM done_lessons WHERE user_id = $1 AND theme_id = $2"
	err := pgxscan.Select(context.Background(), db.databasePool, &record, query, user_id, theme_id)
	if err != nil {
		return record[0], err
	}
	return record[0], nil
}
