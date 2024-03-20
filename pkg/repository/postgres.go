package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable         = "users"
	todoListsTable     = "todo_lists"
	usersListsTable    = "users_lists"
	todoItemsTable     = "todo_items"
	listsItemsTable    = "lists_items"
	rolesTable         = "roles"
	usersRolesTable    = "users_roles"
	statusesTable      = "statuses"
	usersStatusesTable = "users_statuses"
	commentsTable      = "comments"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBname   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	//подставили значения из конфига для подключения к бд
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBname, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
