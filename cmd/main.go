package main

import (
	"os"

	gamereview "github.com/Neokrid/game-review"
	"github.com/Neokrid/game-review/pkg/handler"
	"github.com/Neokrid/game-review/pkg/repository"
	"github.com/Neokrid/game-review/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка инициализации конфигов: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Ошибка загрузки env переменных: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWRD"),
		DBName:   viper.GetString("db.dbname"),
		SSlMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("Ошибка при подключении к БД: %s", err.Error())
	}

	redis, err := repository.NewRedisDB(repository.ConfigRedis{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	if err != nil {
		logrus.Fatalf("Ошибка при подключении к Redis: %s", err.Error())
	}

	repos := repository.NewRepository(db, redis)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(gamereview.Server)
	if err := server.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Произошла ошибка при запуске HTTP-сервера")
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
