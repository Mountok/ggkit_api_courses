package db

import (
	"context"
	"ggkit_learn_service/internals/app/models"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type ThemesStorage struct {
	databasePool *pgxpool.Pool
}

func NewThemesStorage(databasePool *pgxpool.Pool) *ThemesStorage {
	storage := new(ThemesStorage)
	storage.databasePool = databasePool
	return storage
}

func (db *ThemesStorage) CreateTheme(title, description, subject_id string) (int, error) {
	query := "INSERT INTO themes (title,description,subject_id) VALUES ($1,$2,$3) RETURNING id;"
	var id int
	err := db.databasePool.QueryRow(context.Background(), query, title, description, subject_id).Scan(&id)
	log.Errorf("Ошибка при sql запросе: \n %v", err)
	return id, err
}

func (db *ThemesStorage) GetThemesBySubjectId(id int) (result []models.Theme) {
	query := "SELECT id, title, description, subject_id FROM themes WHERE subject_id = $1"
	err := pgxscan.Select(context.Background(), db.databasePool, &result, query, id)
	if err != nil {
		log.Errorf("Ошибка при sql запросе: \n %v", err)
		log.Fatalln(err)
	}
	return result
}

func (db *ThemesStorage) GetAllCompleted(user_id, subject_id string) ([]int, error) {
	var result []int
	query := "SELECT dl.theme_id FROM done_lessons dl WHERE user_id = $1 AND theme_id IN (SELECT id FROM themes WHERE subject_id = $2)"
	err := pgxscan.Select(context.Background(), db.databasePool, &result, query, user_id, subject_id)
	if err != nil {
		log.Errorf("Ошибка при sql запросе: \n %v", err)
		return nil, err
	}
	log.Print(result)
	return result, nil
}

// SELECT dl.theme_id FROM done_lessons dl WHERE user_id = 'b43a1721-2bc3-4421-8e70-b7bd932ad802' AND theme_id IN (SELECT id FROM themes WHERE subject_id = 2);
// SELECT dl.theme_id FROM done_lessons dl JOIN themes t ON t.id = dl.theme_id AND t.subject_id = 2 AND dl.user_id = 'b43a1721-2bc3-4421-8e70-b7bd932ad802';