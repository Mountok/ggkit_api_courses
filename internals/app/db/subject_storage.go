package db

import (
	"context"
	"fmt"
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
func (db *SubjectStorage) UpdateSubject(subject_id, title, description, image_url string) (int, error) {
	var updatedQuestionId int
	var argumentString string
	if title != "" {
		argumentString += fmt.Sprintf(", title='%s'", title)
	}
	if description != "" {
		argumentString += fmt.Sprintf(", description='%s'", description)
	}
	if image_url != "" {
		argumentString += fmt.Sprintf(", image='%s'", image_url)
	}
	query := fmt.Sprintf("UPDATE subjects SET id=$1 %s WHERE id=$2 RETURNING id;", argumentString)

	row := db.databasePool.QueryRow(context.Background(), query, subject_id, subject_id)
	err := row.Scan(&updatedQuestionId)
	if err != nil {
		return 0, err
	}

	fmt.Println(query)

	return updatedQuestionId, nil
}

func (db *SubjectStorage) GetSubjectById(id int) ([]models.Subject, error) {
	query := "SELECT id, title, image, description FROM subjects WHERE id = $1;"

	var result []models.Subject
	err := pgxscan.Select(context.Background(), db.databasePool, &result, query, id)
	if err != nil {
		log.Errorln(err)
		return result, err
	}
	return result, nil

}


func (db *SubjectStorage) DeleteSubject(id string) error {
	query := "DELETE FROM subjects WHERE id=$1"
	_, err := db.databasePool.Exec(context.Background(),query,id)
	return err
}