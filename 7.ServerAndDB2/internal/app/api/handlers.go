package api

import (
	"encoding/json"
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/k-airos/7.ServerAndDB2/internal/app/middleware"
	"github.com/k-airos/7.ServerAndDB2/internal/app/models"
	"net/http"
	"strconv"
	"time"
)

// Message Вспомогательная структура для формирования сообщений
type Message struct {
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
		msg := Message{
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
		msg := Message{
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
		msg := Message{
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
		msg := Message{
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
		msg := Message{
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
		msg := Message{
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
		msg := Message{
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
		msg := Message{
			StatusCode: 501,
			Message:    "We have some troubles to accessing database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(202)
	msg := Message{
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
		msg := Message{
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
		msg := Message{
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
		msg := Message{
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
		msg := Message{
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
		msg := Message{
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
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles to accessing database. Try later",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	msg := Message{
		StatusCode: 201,
		Message:    fmt.Sprintf("User {login: %s} successfully registred!", userAdded.Login),
		IsError:    false,
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(msg)
}

func (a *API) PostToAuth(writer http.ResponseWriter, r *http.Request) {
	initHeaders(writer)
	a.logger.Info("Post to Auth POST /api/v1/user/auth")
	var userFromJSON models.User
	//Обрабатывается случай, когда json - вовсе не json или в нем какие-то проблемы
	if err := json.NewDecoder(r.Body).Decode(&userFromJSON); err != nil {
		if err != nil {
			a.logger.Info("Invalid json recieved from client", err)
			msg := Message{
				StatusCode: 400,
				Message:    "Unapropriate provided json",
				IsError:    true,
			}
			writer.WriteHeader(400)
			json.NewEncoder(writer).Encode(msg)
			return
		}
	}
	// Необходимо попытаться обнаружить пользователя с таким логин в бд
	userInDB, ok, err := a.storage.User().FindByLogin(userFromJSON.Login)
	// проблема доступа к БД
	if err != nil {
		a.logger.Info("Can not make user search in database")
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles while accessing databasse",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	// Если такого пользователя нет
	if !ok {
		a.logger.Info("Can not find user with that ID in database")
		msg := Message{
			StatusCode: 404,
			Message:    "user with this id doesn't exist in database.",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//проверим совпадает пароль юзера с пришедшим
	if userInDB.Password != userFromJSON.Password {
		a.logger.Info("Invalid credentials to auth")
		msg := Message{
			StatusCode: 404,
			Message:    "Your password is invalid",
			IsError:    true,
		}
		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	//Теперь выбиваем токен как знак успешной аутентификации
	token := jwt.New(jwt.SigningMethodHS256)             // Тот же метод подписания токена, что и в JWTMiddleware.go
	claims := token.Claims.(jwt.MapClaims)               //дополнительные действия для шифрования
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() //Время жизни токена
	claims["name"] = userInDB.Login
	claims["admin"] = true
	tokenString, err := token.SignedString(middleware.SecretKey)
	if err != nil {
		a.logger.Info("Can not claim jwt-token")
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles. Try again",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//В случае, если токен успешно выбит - отдаем его клинету
	msg := Message{
		StatusCode: 201,
		Message:    tokenString,
		IsError:    false,
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(msg)
}
