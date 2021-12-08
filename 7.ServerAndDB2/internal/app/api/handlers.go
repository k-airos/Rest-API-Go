package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/k-airos/7.ServerAndDB2/internal/app/models"
	"net/http"
	"strconv"
)

// Messagge Вспомогательная структура для формирования сообщений
type Messagge struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

//Full API Handler initialisation file
func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-type", "application/json")
}

//Returns all current articles from DB
func (api *API) GetAllArticles(writer http.ResponseWriter, r *http.Request) {
	//Initialise headers
	initHeaders(writer)
	//logging the start moment of processing request
	api.logger.Info("Get all articles GET /api/v1/articles")

	articles, err := api.storage.Article().SelectAll()
	if err != nil {
		//Что делаем, если была ошибка на этапе подключения?
		api.logger.Info("Error while Articles.SelectAll: ", err)
		msg := Messagge{
			StatusCode: 501,
			Message:    "We have some troubles to accessing database. Try again later",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(articles)
}

func (api *API) GetArticleByID(writer http.ResponseWriter, r *http.Request) {
	initHeaders(writer)
	api.logger.Info("Get Article by ID /api/v1/articles{id}")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		api.logger.Info("Troubles while parsing {id} param", err)
		msg := Messagge{
			StatusCode: 400,
			Message:    "Unapropriate id value. Don't use ID as uncasting to int value",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	article, ok, err := api.storage.Article().FindArticleById(id)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (articles) with id, err", err)
		msg := Messagge{
			StatusCode: 500,
			Message:    "We have some troubles to accessing database. Try later",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if !ok {
		api.logger.Info("Can not find article with that ID in database")
		msg := Messagge{
			StatusCode: 404,
			Message:    "Article with this id does't exist in database.",
			IsError:    true,
		}
		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
	}
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(article)
}
func (api *API) DeleteArticleByID(writer http.ResponseWriter, r *http.Request) {
	initHeaders(writer)
	api.logger.Info("Delete Article by ID DELETE /api/v1/articles/{id}")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		api.logger.Info("Troubles while parsing {id} param", err)
		msg := Messagge{
			StatusCode: 400,
			Message:    "Unapropriate id value. Don't use ID as uncasting to int value",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	_, ok, err := api.storage.Article().FindArticleById(id)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (articles) with id, err", err)
		msg := Messagge{
			StatusCode: 500,
			Message:    "We have some troubles to accessing database. Try later",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if !ok {
		api.logger.Info("Can not find article with that ID in database")
		msg := Messagge{
			StatusCode: 404,
			Message:    "Article with this id doesn't exist in database.",
			IsError:    true,
		}
		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	_, err = api.storage.Article().DeleteByID(id)
	if err != nil {
		api.logger.Info("Troubles while deleting database table (articles) with id. err:", err)
		msg := Messagge{
			StatusCode: 501,
			Message:    "We have some troubles to accessing database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(202)
	msg := Messagge{
		StatusCode: 202,
		Message:    fmt.Sprintf("Acticle with ID %d successfully deleted.", id),
		IsError:    false,
	}
	json.NewEncoder(writer).Encode(msg)
}
func (api *API) PostArticle(writer http.ResponseWriter, r *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post Article POST /api/v1/articles")
	var article models.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		api.logger.Info("Invalid json received from client")
		msg := Messagge{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	a, err := api.storage.Article().Create(&article)
	if err != nil {
		api.logger.Info("Troubles while creating new article: ", err)
		msg := Messagge{
			StatusCode: 501,
			Message:    "We have some troubles to accessing database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(a)
}
func (api *API) PostUserRegister(writer http.ResponseWriter, r *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post User Register POST /api/v1/user/register")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		api.logger.Info("Invalid json received from client")
		msg := Messagge{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//Пытаемся найти пользователя с таким логином в DB
	_, ok, err := api.storage.User().FindByLogin(user.Login)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (users) with id, err", err)
		msg := Messagge{
			StatusCode: 500,
			Message:    "We have some troubles to accessing database. Try later",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if ok {
		api.logger.Info("User with this ID already exists")
		msg := Messagge{
			StatusCode: 400,
			Message:    "User with this ID already exists",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//Теперь пытаемся добавить в бд
	userAdded, err := api.storage.User().Create(&user)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (users) with id, err", err)
		msg := Messagge{
			StatusCode: 500,
			Message:    "We have some troubles to accessing database. Try later",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	msg := Messagge{
		StatusCode: 201,
		Message:    fmt.Sprintf("User {login: %s} successfully registred!", userAdded.Login),
		IsError:    false,
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(msg)
}
