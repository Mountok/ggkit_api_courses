package cfg

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Cfg struct {
	Port   string
	DBname string
	DBuser string
	DBpass string
	DBhost string
	DBport string
}

func LoadConfig() Cfg {
	v := viper.New()
	v.SetEnvPrefix("GGKIT_SERV")
	v.Set("PORT", os.Getenv("PORT"))
	// v.Set("DBNAME", "postgres")
	// v.Set("DBUSER", "postgres")
	// v.Set("DBPASS", "admin")
	// v.Set("DBHOST", "127.0.0.1")
	// v.Set("DBPORT", "5436")
	v.Set("DBNAME", "railway")
	v.Set("DBUSER", "postgres")
	v.Set("DBPASS", "JMBBpmeyasyiQWhdpLxjESwTwsocyehv")
	v.Set("DBHOST", "junction.proxy.rlwy.net")
	v.Set("DBPORT", "38705")
	v.AutomaticEnv()

	var cfg Cfg

	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Panic(err)
	}
	return cfg

}

func (cfg *Cfg) GetDBConnetcUrl() string { //маленький метод для сборки строки соединения с БД
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBuser,
		cfg.DBpass,
		cfg.DBhost,
		cfg.DBport,
		cfg.DBname,
	)
}

func (cfg *Cfg) NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "autorack.proxy.rlwy.net:52228",    // адрес вашего Redis сервера
		Password: "ndauicaHkkzYaCjcFlZZcaMzeHZetGWk", // пароль, если установлен
		DB:       0,                                  // номер базы данных
	})

	// Проверка подключения
	log.Println("Connecting to Redis")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	log.Printf("Connect to Redis %s\n", rdb.Options().Addr)

	return rdb
}

// redis://default:ndauicaHkkzYaCjcFlZZcaMzeHZetGWk@autorack.proxy.rlwy.net:52228