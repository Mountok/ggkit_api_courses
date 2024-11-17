package db

import (
	"context"
	"ggkit_learn_service/internals/app/models"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type LessonsStorage struct {
	databasePool *pgxpool.Pool
}

func NewLessonsStorage(pool *pgxpool.Pool) *LessonsStorage {
	storage := new(LessonsStorage)
	storage.databasePool = pool
	return storage
}

func (db *LessonsStorage) CreateLesson(theme_id, theme_html string) error {

	tx, err := db.databasePool.Begin(context.Background())
	if err != nil {
		return err
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

	var isLessonCreate int = 0
	row := tx.QueryRow(context.Background(),"SELECT count(id) FROM lessons WHERE theme_id=$1;",theme_id);
	err = row.Scan(&isLessonCreate)
	if err != nil {
		return err
	}
	
	if isLessonCreate == 0 {
		_, err := tx.Exec(context.Background(), "INSERT INTO lessons (upkeep,theme_id) VALUES ($1,$2);", theme_html, theme_id)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.Exec(context.Background(),"UPDATE lessons SET upkeep=$1 WHERE theme_id=$2",theme_html,theme_id)
		if err != nil {
			return err
		}
	}


	
	return err

}

func (db *LessonsStorage) GetLesson(subjectId, themeId int) []models.Lesson {
	var result []models.Lesson
	query := "select * from lessons where theme_id in (select id from themes where subject_id = $1 and id = $2);"
	err := pgxscan.Select(context.Background(), db.databasePool, &result, query, subjectId, themeId)
	if err != nil {
		log.Errorf("Ошибка при sql запросе: \n %v", err)
		return nil
	}
	return result
}
