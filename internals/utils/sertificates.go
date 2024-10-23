package utils

import (
	"context"
	"ggkit_learn_service/internals/app/models"
	log	"github.com/sirupsen/logrus"

	"github.com/georgysavva/scany/pgxscan"

)


func GetAllCompletedThemes(db pgxscan.Querier, userId, subjectId string) ([]int, error) {
	var result []int
	query := "SELECT dl.theme_id FROM done_lessons dl WHERE user_id = $1 AND theme_id IN (SELECT id FROM themes WHERE subject_id = $2)"
	err := pgxscan.Select(context.Background(), db, &result, query, userId, subjectId)
	if err != nil {
		log.Errorf("Ошибка при sql запросе: \n %v", err)
		return []int{}, err
	}
	return result, nil
}

func GetAllThemes(db pgxscan.Querier, subjectId string) ([]int, error) {
	var result []int
	query := "SELECT id  FROM themes WHERE subject_id = $1"
	err := pgxscan.Select(context.Background(), db, &result, query, subjectId)
	if err != nil {
		log.Errorf("Ошибка при sql запросе: \n %v", err)
		log.Fatalln(err)
		return []int{}, err
	}
	return result, nil
}

func GetAllCompletedTest(db pgxscan.Querier, userId, subjectId string) ([]models.CompletedTestCheck, error) {
	var completedTestId []models.CompletedTestCheck
	log.Printf("Получение выполненых тестов для userId: %s по subjectID: %d", userId, subjectId)
	query := `
	SELECT 
		dt.id, 
		dt.test_id, 
		st.subject_id, 
		dt.user_id, 
		dt.points,
		COUNT(tq.id) AS question_count
	FROM 
		done_test dt 
	JOIN 
		subject_test st ON st.id = dt.test_id  
	LEFT JOIN 
		test_questions tq ON tq.test_id = dt.test_id
	WHERE 
		dt.user_id = $1
		AND st.subject_id = $2
	GROUP BY 
		dt.id, dt.test_id, st.subject_id, dt.user_id, dt.points
	ORDER BY 
		dt.test_id ASC;
	`
	err := pgxscan.Select(context.Background(), db, &completedTestId, query, userId, subjectId)
	if err != nil {
		return completedTestId, err
	}
	return completedTestId, nil
}

func GetAllTestBySubject(db pgxscan.Querier, subjectId string) ([]int, error) {
	var subjectTests []int
	query := "SELECT id  FROM subject_test WHERE subject_id=$1;"
	log.Infof("Запрос на получение тестов по курсу с id %s", subjectId)
	err := pgxscan.Select(context.Background(), db, &subjectTests, query, subjectId)
	if err != nil {
		log.Errorf("Запрос на получение тестов по курсу с id %s не удался:\n%v", subjectId, err)
		return subjectTests, err
	}
	log.Infof("Запрос на получение тестов по курсу с id %s успешен", subjectId)
	return subjectTests, nil
}
