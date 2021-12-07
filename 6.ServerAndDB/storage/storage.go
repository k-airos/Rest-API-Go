package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

//Instance of storage
type Storage struct {
	config *Config //Отвечает за информацию "как" подключиться к базе данных
	db     *sql.DB //Database file descriptor
}

//Storage Constructor
func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

//Open connection method
func (storage *Storage) Open() error {
	db, err := sql.Open("postgres", storage.config.DatabaseURI)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	storage.db = db
	log.Println("Database connection created successfully")
	return nil
}

//Close connection
func (storage *Storage) Close() {
	if err := storage.db.Close(); err != nil {
		log.Println("Connection already close")
	}
}
