package rdb

import (
	"context"
	"encoding/json"
	"fmt"
	"ggkit_learn_service/internals/app/db"
	"ggkit_learn_service/internals/app/models"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type SubjectCache struct {
	storage *db.SubjectStorage
	redis   *redis.Client
}

func NewSubjectCache(storage *db.SubjectStorage, rdb *redis.Client) *SubjectCache {
	cache := new(SubjectCache)
	cache.storage = storage
	cache.redis = rdb
	return cache
}

func (db *SubjectCache) GetAllSubjects() ([]models.Subject, error) {
	var result []models.Subject
	cachedSubjects, err := db.redis.Get(context.Background(), "subjects").Result()
	switch {
	case err == redis.Nil:
		storageSubjects, err := db.storage.GetAllSubjects()
		if err != nil {
			logrus.Errorf("ошибка при получении предметов из БД:\n%v", err)
			return result, err
		}
		jsonData, err := json.Marshal(storageSubjects)
		if err != nil {
			logrus.Errorf("ошибка при парсинге в json:\n%v", err)
			return result, err
		}
		err = db.redis.Set(context.Background(), "subjects", jsonData, 1*time.Minute).Err()
		if err != nil {
			logrus.Errorf("ошибка кэширования: %v", err)
			return result, err
		}
		// logrus.Infof("Данные из БД: %v\n", storageSubjects)
		return storageSubjects, nil
	case err != nil:
		logrus.Errorf("ошибка получения данных из кэша redis:\n%v",err)
		return result,err
	default:
		// logrus.Infof("полученные данные из кэша:\n%v",cachedSubjects) 
		err := json.Unmarshal([]byte(cachedSubjects),&result)
		if err != nil {
			logrus.Errorf("не удалось спарсить в json:\n%v",err)
			return result,err
		}
		return result, nil
	}
}


func (db *SubjectCache) UploadStorage(title, description,image_url string) (int, error) {
	return db.storage.UploadStorage(title, description,image_url)
}


func (db *SubjectCache) GetSubjectById(num int) ([]models.Subject, error) {
	var findSubject []models.Subject
	cachedSubject, err := db.redis.Get(context.Background(),fmt.Sprintf("subject_id_%d",num)).Result()
	switch{
	case err == redis.Nil:
		storageSubject, err := db.storage.GetSubjectById(num)
		if err != nil {
			logrus.Errorf("ошибка при получении предметов из БД:\n%v", err)
			return findSubject, err
		}
		jsonData, err  := json.Marshal(storageSubject)
		if err !=  nil {
			logrus.Errorf("ошибка при парсинге в json:\n%v", err)
			return findSubject, err
		}
		err = db.redis.Set(context.Background(),fmt.Sprintf("subject_id_%d",num),jsonData,1*time.Minute).Err()
		if err != nil {
			logrus.Errorf("ошибка кэширования: %v", err)
			return findSubject, err
		}
		return storageSubject, nil
	case err != nil:
		logrus.Errorf("ошибка получения данных из кэша redis:\n%v",err)
		return findSubject,err
	default:
		err := json.Unmarshal([]byte(cachedSubject),&findSubject)
		if err != nil {
			logrus.Errorf("не удалось спарсить в json:\n%v",err)
			return findSubject,err
		}
		return findSubject, err
	}
}