package db

import (
	"context"
	"fmt"
	"ggkit_learn_service/internals/app/models"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type SubjectTestStorage struct {
	databasePool *pgxpool.Pool
}

func NewSubjectTestStorage(databasePool *pgxpool.Pool) *SubjectTestStorage {
	storage := new(SubjectTestStorage)
	storage.databasePool = databasePool
	return storage
}

// ----------------------- ТЕСТЫ

func (db *SubjectTestStorage) CreateTestForSubject(testTitle, subjectId string) (int, error) {
	var testId int
	query := "INSERT INTO subject_test(title,subject_id) VALUES($1,$2) RETURNING id;"
	logrus.Infof("Запрос на создание теста по курсу с id %s", subjectId)

	row := db.databasePool.QueryRow(context.Background(), query, testTitle, subjectId)
	err := row.Scan(&testId)
	if err != nil {
		logrus.Errorf("Ошбика при создании теста по id %s: %v", subjectId, err)
		return 0, err
	}
	logrus.Infof("Создан тест по курсу с id %s", subjectId)

	return testId, nil
}

func (db *SubjectTestStorage) ReadTestsForSubject(subjectId string) ([]models.SubjectTest, error) {
	var subjectTests []models.SubjectTest
	query := "SELECT id, title, subject_id FROM subject_test WHERE subject_id=$1;"
	logrus.Infof("Запрос на получение тестов по курсу с id %s", subjectId)
	err := pgxscan.Select(context.Background(), db.databasePool, &subjectTests, query, subjectId)
	if err != nil {
		logrus.Errorf("Запрос на получение тестов по курсу с id %s не удался:\n%v", subjectId, err)
		return subjectTests, err
	}
	logrus.Infof("Запрос на получение тестов по курсу с id %s успешен", subjectId)
	return subjectTests, nil
}

func (db *SubjectTestStorage) UpdateTestForSubject(newTitle, subjectId, lastTitle string) (int, error) {
	var testId int
	query := "UPDATE subject_test SET title=$1 WHERE subject_id=$2 and title=$3 RETURNING id;"
	logrus.Infof("Запрос на обновление названия теста по курсу id %s", subjectId)

	row := db.databasePool.QueryRow(context.Background(), query, newTitle, subjectId, lastTitle)
	err := row.Scan(&testId)
	if err != nil {
		logrus.Errorf("Ошибка при обновлени названия теста для курса с id %s: %v", subjectId, err)
		return 0, err
	}
	logrus.Infof("Обновлено название теста с %s на %s по курсу с id %s", lastTitle, newTitle, subjectId)

	return testId, nil
}

func (db *SubjectTestStorage) DeleteTestForSubject(testTitle, subjectId string) error {
	query := "DELETE FROM subject_test WHERE title=$1 and subject_id=$2"
	logrus.Infof("Запрос на удаление теста по курсу с id %s", subjectId)
	_, err := db.databasePool.Exec(context.Background(), query, testTitle, subjectId)
	if err != nil {
		logrus.Errorf("Ошбика при удалении теста по id %s: %v", subjectId, err)
		return err
	}
	logrus.Infof("Удален тест по курсу с id %s", subjectId)

	return nil
}

// ----------------------- ВОПРОСЫ В ТЕСТАХ

func (db *SubjectTestStorage) CreateQuestionForTest(testId, question, options, answer string) (int, error) {
	var questionId int
	query := "INSERT INTO test_questions(question,options,answer,test_id) VALUES($1,$2,$3,$4) RETURNING id;"
	logrus.Infof("Запрос на создание вопроса для теста с id %s", testId)
	row := db.databasePool.QueryRow(context.Background(), query, question, options, answer, testId)
	err := row.Scan(&questionId)
	if err != nil {
		logrus.Errorf("Запрос на создание вопроса для теста с id %s - не удался: \n %v", testId, err)
		return questionId, err
	}
	logrus.Infof("Создан вопрос для теста с id %s", testId)

	return questionId, nil
}

func (db *SubjectTestStorage) ReadQuestionForTest(testId string) ([]models.TestQuestion, error) {
	var questionsList []models.TestQuestion
	query := "SELECT tq.id, tq.question, tq.options, tq.test_id, st.subject_id FROM test_questions tq JOIN subject_test st ON tq.test_id = st.id WHERE tq.test_id=$1;"
	logrus.Infof("Запрос на получение вопросов для теста с id %s", testId)
	err := pgxscan.Select(context.Background(), db.databasePool, &questionsList, query, testId)
	if err != nil {
		logrus.Errorf("Запрос на получения вопросов для теста с id %s - не удался: \n %v", testId, err)
		return questionsList, err
	}
	logrus.Infof("Получение вопрос для теста с id %s - прошел успешно", testId)
	return questionsList, nil
}

func (db *SubjectTestStorage) UpdateQuestionForTest(testId, question, options, answer string) (int, error) {
	var updatedQuestionId int
	var argumentString string
	if question != "" {
		argumentString += fmt.Sprintf(", question='%s'", question)
	}
	if options != "" {
		argumentString += fmt.Sprintf(", options='%s'", options)
	}
	if answer != "" {
		argumentString += fmt.Sprintf(", answer='%s'", answer)
	}
	query := fmt.Sprintf("UPDATE test_questions SET test_id=$1 %s WHERE test_id=$2 RETURNING id;", argumentString)

	row := db.databasePool.QueryRow(context.Background(), query, testId, testId)
	err := row.Scan(&updatedQuestionId)
	if err != nil {
		return 0, err
	}

	fmt.Println(query)

	return updatedQuestionId, nil
}

func (db *SubjectTestStorage) DeleteQuestionForTest(testId, questionId string) error {
	query := "DELETE FROM test_questions WHERE test_id=$1 and id=$2;"
	logrus.Infof("Запрос на удаление вопроса с id=%s для теста с id=%s", questionId, testId)
	_, err := db.databasePool.Exec(context.Background(), query, testId, questionId)
	if err != nil {
		logrus.Errorf("Запрос на удаление вопроса с id=%s для теста с id=%s - не удался: \n %v", questionId, testId, err)
		return err
	}
	logrus.Infof("Запрос на удаление вопроса с id=%s для теста с id=%s - успешно прошол", questionId, testId)
	return nil
}

// ----------------------- ВЫПОЛНЕНЫЕ ТЕСТЫ

func (db *SubjectTestStorage) CreateCompletedTest(userId, testId string, points int) (int, error) {
	var completedTestId int
	tx, err := db.databasePool.Begin(context.Background())
	if err != nil {
		return 0, err
	}

	defer func ()  {
		if err != nil {
			tx.Rollback(context.Background())
			log.Printf("Transaction rolback with error: %s \n", err)
		} else {
			err := tx.Commit(context.Background())
			if err != nil {
				log.Fatalf("Unable to commit transaction: %v", err)
			}
		}
	}()

	row := tx.QueryRow(context.Background(), "INSERT INTO done_test(test_id,user_id,points) VALUES ($1,$2,$3) RETURNING id;", testId, userId, points)
	err = row.Scan(&completedTestId)

	_, err = tx.Exec(context.Background(),"UPDATE profiles SET score=score+$1 WHERE user_id=$2",points,userId)
	if err != nil {
		return 0, err
	}

	return completedTestId, nil
}

func (db *SubjectTestStorage) ReadCompletedTest(subject_id int, userId string) ([]models.CompletedTestCheck, error) {
	var completedTestId []models.CompletedTestCheck
	logrus.Printf("Получение выполненых тестов для userId: %s по subjectID: %d", userId, subject_id)
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
	err := pgxscan.Select(context.Background(), db.databasePool, &completedTestId, query, userId, subject_id)
	if err != nil {
		return completedTestId, err
	}
	return completedTestId, nil
}

func (db *SubjectTestStorage) ReadCompletedTestById(subject_id int, userId, test_id string) ([]models.CompletedTestCheck, error) {
	var completedTestId []models.CompletedTestCheck
	logrus.Printf("Получение выполненых тестов для userId: %s по subjectID: %d", userId, subject_id)
	query := "SELECT dt.id, dt.test_id, st.subject_id, dt.user_id, dt.points FROM done_test dt JOIN subject_test st ON st.id = dt.test_id  WHERE dt.user_id = $1 AND st.subject_id = $2 AND st.id = $3 ORDER BY dt.test_id ASC;"
	err := pgxscan.Select(context.Background(), db.databasePool, &completedTestId, query, userId, subject_id, test_id)
	if err != nil {
		return completedTestId, err
	}
	return completedTestId, nil
}

func (db *SubjectTestStorage) UpdateCompletedTest(testId, userId string, points int) (int, error) {
	var updatedTestId int

	tx, err := db.databasePool.Begin(context.Background())

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

	var lastPoints int = 0
	row := tx.QueryRow(context.Background(), "SELECT points FROM done_test WHERE user_id=$1 and test_id=$2;", userId, testId)
	err = row.Scan(&lastPoints)
	if err != nil {
		return lastPoints, err
	}

	_, err = tx.Exec(context.Background(), "UPDATE profiles SET score=score-$1 WHERE user_id=$2;", lastPoints, userId)
	if err != nil {
		return 0, err
	}

	query := "UPDATE done_test SET points=$1 WHERE user_id=$2 and test_id=$3 RETURNING id;"
	row = tx.QueryRow(context.Background(), query, points, userId, testId)
	err = row.Scan(&updatedTestId)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(context.Background(),"UPDATE profiles SET score=score+$1 WHERE user_id=$2;",points,userId)
	if err != nil {
		return 0, err
	}

	return updatedTestId, nil
}

func (db *SubjectTestStorage) DeleteCompletedTest(testId, userId, completedId string) error {
	query := "DELETE FROM done_test WHERE test_id=$1 and user_id=$2 and id=$3;"
	_, err := db.databasePool.Exec(context.Background(), query, testId, userId, completedId)
	if err != nil {
		return err
	}
	return nil
}

// --------------------- ПРОВЕРКА ВОПРОСОВ

func (db *SubjectTestStorage) CheckQuestion(answers []models.QuestionCheckReq, test_id, user_id string, subject_id int) (int, error) {
	var resultPoints int
	var checkNum []int
	for i := 0; i < len(answers); i++ {
		var query string = fmt.Sprint("SELECT COUNT(*) FROM test_questions WHERE id=$1 and answer=$2")
		err := pgxscan.Select(context.Background(), db.databasePool, &checkNum, query,answers[i].QuestionId, answers[i].Answer)
		if err != nil {
			return 0, err
		}
		if checkNum[0] != 0 {
			resultPoints++
		}
		logrus.Println(query)
		logrus.Println(resultPoints)
	}
	isCompleted, err := db.ReadCompletedTestById(subject_id, user_id, test_id)
	if err != nil {
		return 0, err
	}
	if len(isCompleted) == 1 {
		_, err = db.UpdateCompletedTest(test_id, user_id, resultPoints)
		if err != nil {
			return 0, err
		}
	} else {
		_, err = db.CreateCompletedTest(user_id, test_id, resultPoints)
		if err != nil {
			return 0, err
		}

	}
	return resultPoints, err
}

// ---------------------------  ВЫПОЛНЕНЫЕ ТЕСТЫ

func (db *SubjectTestStorage) GetAllCompleted(user_id string) (int, error) {
	var (
		query string = "SELECT count(id) from done_test WHERE user_id=$1;"
		num   int    =  0
	)
	row := db.databasePool.QueryRow(context.Background(), query, user_id)
	err := row.Scan(&num)
	return num, err
}
