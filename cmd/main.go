package main

import (
	todo "ToDoApp"
	"ToDoApp/pkg/handler"
	"ToDoApp/pkg/repository"
	"ToDoApp/pkg/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"

	_ "github.com/lib/pq"

	"github.com/spf13/viper"
)

func main() {
	//сразу храним логи в json формате
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	//передали параметр из окружения(пароли пользователей) в конфиг
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBname:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	//конструктор
	repos := repository.NewRepository(db)    //зависимость от бд
	services := service.NewService(repos)    //зависимость сервиса от репозитория
	handlers := handler.NewHandler(services) //зависимость запросов от сервиса

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
