package db

import (
	"context"
	"errors"
	"fmt"
	"ggkit_learn_service/internals/app/models"
	"ggkit_learn_service/internals/utils"

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
	query := "SELECT id, title, image, description, iscertificated FROM subjects;"
	err := pgxscan.Select(context.Background(), db.databasePool, &result, query)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (db *SubjectStorage) UploadStorage(title, description, image_url, is_certificated string) (int, error) {
	query := "INSERT INTO subjects (title,description,image,iscertificated) VALUES ($1,$2,$3,$4) RETURNING id;"
	var id int
	err := db.databasePool.QueryRow(context.Background(), query, title, description, image_url,is_certificated).Scan(&id)
	if err != nil {
		return id, err
	}
	query = "UPDATE subjects SET image=$1  WHERE id=$2;"
	_, err = db.databasePool.Exec(context.Background(), query, fmt.Sprintf("%s%d.webp", image_url, id), id)
	if err != nil {
		return id, err
	}
	return id, err
}
func (db *SubjectStorage) UpdateSubject(subject_id,title, description,image_url, is_certificated string) (int, error) {
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
	if is_certificated != "" {
		argumentString += fmt.Sprintf(", iscertificated='%s'", is_certificated)
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
	query := "SELECT id, title, image, description,iscertificated FROM subjects WHERE id = $1;"

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
	_, err := db.databasePool.Exec(context.Background(), query, id)
	return err
}

func (db *SubjectStorage) Certificate(subjectId, userId string) error {
	allComletedThemes, err := utils.GetAllCompletedThemes(db.databasePool, userId, subjectId)
	if err != nil {
		return err
	}
	allThemes, err := utils.GetAllThemes(db.databasePool, subjectId)
	if err != nil {
		return err
	}

	if len(allComletedThemes) != len(allThemes) {
		return errors.New("все темы не пройденый")
	} else {
		fmt.Println(allComletedThemes)
		fmt.Println(allThemes)
		allComletedTests, err := utils.GetAllCompletedTest(db.databasePool, userId, subjectId)
		if err != nil {
			return err
		}
		allTestIds, err := utils.GetAllTestBySubject(db.databasePool, subjectId)
		if err != nil {
			return err
		}
		fmt.Println(allComletedTests)
		if len(allComletedTests) == len(allTestIds) {
			for i := 0; i < len(allComletedTests); i++ {
				points := float64(allComletedTests[i].Points)
				questions := float64(allComletedTests[i].QuestionCount)
				if (points/questions)*100 >= 90 {
					fmt.Printf("Тест выполнен (id теста: %d)\n", allComletedTests[i].TestId)

				} else {
					fmt.Printf("Тест не выполнен (id теста: %d)\n", allComletedTests[i].TestId)
					return errors.New("тесты выполнены неверно")
				}
			}
			return nil
		} else {
			return errors.New("не все тесты пройдены")
		}
	}

}
