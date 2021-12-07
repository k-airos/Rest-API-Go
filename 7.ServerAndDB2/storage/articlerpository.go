package storage

import (
	"fmt"
	"github.com/k-airos/7.ServerAndDB2/internal/app/models"
	"log"
)

// ArticleRepository Article repository instance (model interface)
type ArticleRepository struct {
	storage *Storage
}

var (
	tableArticle string = "articles"
)

//Добавить статью в бд
func (ar *ArticleRepository) Create(a *models.Article) (*models.Article, error) {
	query := fmt.Sprintf("INSERT INTO %s (title, author, content) VALUES ($1,$2,$3) RETURNING id", tableArticle)
	if err := ar.storage.db.QueryRow(query, a.Title, a.Author, a.Content).Scan(&a.ID); err != nil {
		return nil, err
	}
	return a, nil
}

//удалять по id
func (ar *ArticleRepository) DeleteByID(id int) (*models.Article, error) {
	article, ok, err := ar.FindArticleById(id)
	if err != nil {
		return nil, err
	}
	if ok {
		query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", tableArticle)
		_, err := ar.storage.db.Exec(query)
		if err != nil {
			return nil, err
		}
	}
	return article, nil
}

//Получаем статью по id
func (ar *ArticleRepository) FindArticleById(id int) (*models.Article, bool, error) {
	articles, err := ar.SelectAll()
	var founded bool
	if err != nil {
		return nil, founded, err
	}
	var articleFinded *models.Article
	for _, a := range articles {
		if a.ID == id {
			articleFinded = a
			founded = true
			break
		}
	}
	return articleFinded, founded, nil
}

//Получим все статьи в бд
func (ar *ArticleRepository) SelectAll() ([]*models.Article, error) {
	query := fmt.Sprintf("SELECT * from %s", tableArticle)
	rows, err := ar.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//Подготовим, куда будем считать
	articles := make([]*models.Article, 0)
	for rows.Next() {
		a := models.Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Author, &a.Content)
		if err != nil {
			log.Println(err) // нам незачем фаталиться здесь
			continue
		}
		articles = append(articles, &a)
	}
	return articles, nil
}
