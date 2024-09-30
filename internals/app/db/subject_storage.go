package db

import (
	"context"
	"ggkit_learn_service/internals/app/models"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type SubjectStorage struct {
	databasePool *pgxpool.Pool
}

func NewSubjectStorage(pool *pgxpool.Pool) *SubjectStorage {
	storage := new(SubjectStorage)
	storage.databasePool = pool
	return storage
}

func (db *SubjectStorage) GetAllSubjects() ([]models.Subject, error) {
	var result []models.Subject
	query := "SELECT id, title, image, description FROM subjects;"
	err := pgxscan.Select(context.Background(), db.databasePool, &result, query)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (db *SubjectStorage) UploadStorage(title, description, image_url string) (int, error) {
	query := "INSERT INTO subjects (title,description,image) VALUES ($1,$2,$3) RETURNING id;"
	var id int
	err := db.databasePool.QueryRow(context.Background(), query, title, description, image_url).Scan(&id)
	return id, err
}

func (db *SubjectStorage) GetSubjectById(id int) []models.Subject {
	query := "SELECT id, title, image, description FROM subjects WHERE id = $1;"

	var result []models.Subject
	err := pgxscan.Select(context.Background(), db.databasePool, &result, query, id)
	if err != nil {
		log.Errorln(err)
	}
	return result

}
