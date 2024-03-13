package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type StatusPostgres struct {
	db *sqlx.DB
}

func NewStatusPostgres(db *sqlx.DB) *StatusPostgres {
	return &StatusPostgres{db: db}
}

func (r *StatusPostgres) Create(description string) (int, error) {
	var id int
	insertQuery := fmt.Sprintf("INSERT INTO %s (description) VALUES ($1) RETURNING id", statusesTable)
	err := r.db.QueryRow(insertQuery, description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *StatusPostgres) GetAll() ([]Status, error) {
	var statuses []Status
	query := fmt.Sprintf("SELECT id, description FROM %s", statusesTable)
	err := r.db.Select(&statuses, query)
	return statuses, err
}

func (r *StatusPostgres) Delete(statusId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Удаление связей статуса
	deleteUserStatusQuery := fmt.Sprintf("DELETE FROM %s WHERE status_id=$1", usersStatusesTable)
	_, err = tx.Exec(deleteUserStatusQuery, statusId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление самого статуса
	deleteStatusQuery := fmt.Sprintf("DELETE FROM %s WHERE id=$1", statusesTable)
	_, err = tx.Exec(deleteStatusQuery, statusId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *StatusPostgres) GetUsersStatuses() ([][]string, error) {
	var usersStatuses [][]string
	query := fmt.Sprintf(`
		SELECT u.name, s.description 
		FROM %s u 
		LEFT JOIN %s us ON u.id = us.user_id 
		LEFT JOIN %s s ON us.status_id = s.id
	`, usersTable, usersStatusesTable, statusesTable)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name, description string
		if err := rows.Scan(&name, &description); err != nil {
			return nil, err
		}
		usersStatuses = append(usersStatuses, []string{name, description})
	}
	return usersStatuses, rows.Err()
}

func (r *StatusPostgres) SetStatus(userId int, statusId int) error {
	// Проверка, есть ли уже статус для пользователя
	checkStatusExistsQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id=$1", usersStatusesTable)
	var count int
	err := r.db.QueryRow(checkStatusExistsQuery, userId).Scan(&count)
	if err != nil {
		return err
	}

	// Если статус уже установлен, обновляем его
	if count > 0 {
		updateStatusQuery := fmt.Sprintf("UPDATE %s SET status_id=$1 WHERE user_id=$2", usersStatusesTable)
		_, err = r.db.Exec(updateStatusQuery, statusId, userId)
		return err
	}

	// Если статус еще не установлен, устанавливаем его
	insertStatusQuery := fmt.Sprintf("INSERT INTO %s (user_id, status_id) VALUES ($1, $2)", usersStatusesTable)
	_, err = r.db.Exec(insertStatusQuery, userId, statusId)
	return err
}

func (r *StatusPostgres) DropStatus(userId int) error {
	dropStatusQuery := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", usersStatusesTable)
	_, err := r.db.Exec(dropStatusQuery, userId)
	return err
}
