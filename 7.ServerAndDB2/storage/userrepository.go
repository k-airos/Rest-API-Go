package storage

import (
	"fmt"
	"github.com/k-airos/7.ServerAndDB2/internal/app/models"
	"log"
)

// UserRepository User repository instance (model interface)
type UserRepository struct {
	storage *Storage
}

var (
	tableUser string = "users"
)

//Create user in db
func (ur *UserRepository) Create(u *models.User) (*models.User, error) {
	query := fmt.Sprintf("INSERT INTO %s (login, password) VALUES ($1, $2) RETURNING id", tableUser)
	if err := ur.storage.db.QueryRow(query, u.Login, u.Password).Scan(&u.ID); err != nil {
		return nil, err
	}
	return u, nil
}

//Find user by login
func (ur *UserRepository) FindByLogin(login string) (*models.User, bool, error) {
	users, err := ur.SelectAll()
	var founded bool
	if err != nil {
		return nil, founded, err
	}
	var userFinded *models.User
	for _, u := range users {
		if u.Login == login {
			userFinded = u
			founded = true
			break
		}
	}

	return userFinded, founded, nil
}

// Select all users in db
func (ur *UserRepository) SelectAll() ([]*models.User, error) {
	query := fmt.Sprintf("SELECT * from %s", tableUser)
	rows, err := ur.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//Подготовим, куда будем считать
	users := make([]*models.User, 0)
	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.ID, &u.Login, &u.Password)
		if err != nil {
			log.Println(err) // нам незачем фаталиться здесь
			continue
		}
		users = append(users, &u)
	}
	return users, nil
}
