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
	//Subfields for repo interfacing
	userRepository    *UserRepository
	articleRepository *ArticleRepository
}

//Storage Constructor
func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

//Open connection method
func (s *Storage) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURI)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db
	log.Println("Database connection created successfully")
	return nil
}

//Close connection
func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		log.Println("Connection already close")
	}
}

//Public repo for article
func (s *Storage) Article() *ArticleRepository {
	if s.articleRepository != nil {
		return s.articleRepository
	}
	s.articleRepository = &ArticleRepository{
		storage: s,
	}
	return s.articleRepository
}

//Public repo for User
func (s *Storage) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{
		storage: s,
	}
	return s.userRepository
}
