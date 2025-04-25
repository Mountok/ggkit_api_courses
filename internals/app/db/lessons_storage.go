package db

import (
	"context"
	"fmt"
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
	row := tx.QueryRow(context.Background(), "SELECT count(id) FROM lessons WHERE theme_id=$1;", theme_id)
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
		_, err = tx.Exec(context.Background(), "UPDATE lessons SET upkeep=$1 WHERE theme_id=$2", theme_html, theme_id)
		if err != nil {
			return err
		}
	}

	return err

}

func (db *LessonsStorage) GetLessonHTML(theme_id string) (string, error) {
	var themeHtml string
	var query string = "SELECT upkeep FROM lessons where theme_id=$1;"
	row := db.databasePool.QueryRow(context.Background(), query, theme_id)
	err := row.Scan(&themeHtml)
	if err != nil {
		return "", err
	}
	return themeHtml, err
}

func (db *LessonsStorage) UpdateLesson(theme_id, theme_html string) error {

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
	log.Infof("TX - Started (%s,%s)", theme_id, theme_html)

	var isLessonCreate int = 0
	row := tx.QueryRow(context.Background(), "SELECT count(id) FROM lessons WHERE theme_id=$1;", theme_id)
	err = row.Scan(&isLessonCreate)
	if err != nil {
		return err
	}

	log.Info("TX 1 - COMPLETE")

	if isLessonCreate == 0 {
		_, err := tx.Exec(context.Background(), "INSERT INTO lessons (upkeep,theme_id) VALUES ($1,$2);", theme_html, theme_id)
		if err != nil {
			return err
		}
		log.Info("TX 2 - COMPLETE")

	} else {
		_, err = tx.Exec(context.Background(), fmt.Sprintf("UPDATE lessons SET upkeep=CONCAT(upkeep,'%s') WHERE theme_id=$1", theme_html), theme_id)
		if err != nil {
			return err
		}
		log.Info("TX 3 - COMPLETE")

	}

	return err

}

func (db *LessonsStorage) GetLesson(subjectId, themeId int) []models.LessonResponse {
	var result []models.LessonResponse
	query := "select ls.id, ls.upkeep, ls.theme_id, th.title from lessons ls LEFT JOIN themes th ON th.id = ls.theme_id WHERE theme_id = $1;"
	err := pgxscan.Select(context.Background(), db.databasePool, &result, query, themeId)
	if err != nil {
		log.Errorf("Ошибка при sql запросе: \n %v", err)
		return nil
	}
	return result
}
