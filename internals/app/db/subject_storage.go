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

func (db *SubjectStorage) GetAllSubjects(userId string) ([]models.SubjectResponse, error) {
	var result []models.SubjectResponse
	query := `
			SELECT 
				subjects.id AS id,
				subjects.title AS title,
				subjects.image  as image,
				subjects.description as description,
				subjects.isCertificated AS isCertificated,
				COUNT(themes.id) AS total_themes,
				COUNT(done_lessons.theme_id) AS completed_themes
			FROM 
				subjects
			LEFT JOIN 
				deleted_subjects ON subjects.id = deleted_subjects.subject_id
			LEFT JOIN 
				themes ON subjects.id = themes.subject_id
			LEFT JOIN 
				done_lessons ON themes.id = done_lessons.theme_id AND done_lessons.user_id = $1 -- замените 'USER_ID_HERE' на реальный ID пользователя
			WHERE 
				deleted_subjects.subject_id IS NULL
			GROUP BY 
				subjects.id, subjects.title, subjects.image, subjects.description, subjects.isCertificated;`
	err := pgxscan.Select(context.Background(), db.databasePool, &result, query, userId)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (db *SubjectStorage) UploadStorage(title, description, image_url, is_certificated string) (int, error) {
	query := "INSERT INTO subjects (title,description,image,iscertificated) VALUES ($1,$2,$3,$4) RETURNING id;"
	var id int
	err := db.databasePool.QueryRow(context.Background(), query, title, description, image_url, is_certificated).Scan(&id)
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
func (db *SubjectStorage) UpdateSubject(subject_id, title, description, image_url, is_certificated string) (int, error) {
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
	query := "INSERT INTO deleted_subjects (subject_id) VALUES ($1);"
	_, err := db.databasePool.Exec(context.Background(), query, id)
	return err
}
func (db *SubjectStorage) GetDeletedSubject() ([]int, error) {
	var ids []int
	query := "SELECT subject_id FROM seleted_subjects;"
	err := pgxscan.Select(context.Background(), db.databasePool, &ids, query)
	return ids, err
}

func (db *SubjectStorage) Certificate(subjectId, userId string) (interface{}, error) {

	allComletedThemes, err := utils.GetAllCompletedThemes(db.databasePool, userId, subjectId)
	if err != nil {
		return "", err
	}
	allThemes, err := utils.GetAllThemes(db.databasePool, subjectId)
	if err != nil {
		return "", err
	}

	if len(allComletedThemes) != len(allThemes) {
		return "", errors.New("все темы не пройденый")
	} else {
		fmt.Println(allComletedThemes)
		fmt.Println(allThemes)
		allComletedTests, err := utils.GetAllCompletedTest(db.databasePool, userId, subjectId)
		if err != nil {
			return "", err
		}
		allTestIds, err := utils.GetAllTestBySubject(db.databasePool, subjectId)
		if err != nil {
			return "", err
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
					return "", errors.New("тесты выполнены неверно")
				}
			}
			var checkCertfRes []string
			checkCertfQuery := "SELECT user_id FROM certificates WHERE user_id=$1"
			err := pgxscan.Select(context.Background(), db.databasePool, &checkCertfRes, checkCertfQuery, userId)
			if err != nil {
				return "", errors.New(
					fmt.Sprintf("ошибка на стороне сервера (SELECT user_id FROM certificates WHERE user_id=$1): \n%s\n", err.Error()),
				)
			}
			if len(checkCertfRes) == 0 {
				query := "INSERT INTO certificates (user_id,subject_id) VALUES ($1,$2);"
				_, err = db.databasePool.Exec(context.Background(), query, userId, subjectId)
				if err != nil {
					return "", errors.New(
						fmt.Sprintf("ошибка на стороне сервера INSERT INTO certificates (user_id,subject_id) VALUES ($1,$2);: \n%s\n", err.Error()),
					)
				}
				var certfRes []interface{}
				certfQuery := "SELECT get_date FROM certificates WHERE user_id=$1"
				err = pgxscan.Select(context.Background(), db.databasePool, &certfRes, certfQuery, userId)
				if len(certfRes) == 0 {
					return "", err
				}
				return certfRes[0], err

			}
			var certfRes []interface{}
			certfQuery := "SELECT get_date FROM certificates WHERE user_id=$1"
			err = pgxscan.Select(context.Background(), db.databasePool, &certfRes, certfQuery, userId)
			if len(certfRes) == 0 {
				return "", err
			}
			return certfRes[0], err
		} else {
			return "", errors.New("не все тесты пройдены")
		}
	}

}
