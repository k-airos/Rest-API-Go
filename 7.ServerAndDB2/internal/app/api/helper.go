package api

import (
	_ "github.com/gorilla/mux"
	"github.com/k-airos/7.ServerAndDB2/internal/app/middleware"
	storage "github.com/k-airos/7.ServerAndDB2/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	prefix string = "/api/v1"
)

//Пытаемся отконфигурировать наш API инстанис (а конкретнее - поле logger)
func (a *API) configureLoggerField() error {
	logLevel, err := logrus.ParseLevel(a.config.LoggerLevel)
	if err != nil {
		return err
	}
	a.logger.SetLevel(logLevel)
	return nil
}

//Пытаемся отконфигурировать маршрутизатор (поле router API)
func (a *API) configureRouterField() {
	a.router.HandleFunc(prefix+"/articles", a.GetAllArticles).Methods("GET")
	// Было до JWT
	//a.router.HandleFunc(prefix+"/articles/{id}", a.GetArticleByID).Methods("GET")
	//Теперь требует наличия JWT
	a.router.Handle(prefix+"/articles/{id}", middleware.JWTMiddleware.Handler(
		http.HandlerFunc(a.GetArticleByID),
	)).Methods("GET")
	//
	a.router.HandleFunc(prefix+"/articles/{id}", a.DeleteArticleByID).Methods("DELETE")
	a.router.HandleFunc(prefix+"/articles", a.PostArticle).Methods("POST")
	a.router.HandleFunc(prefix+"/user/register", a.PostUserRegister).Methods("POST")
	//new pair for auth
	a.router.HandleFunc(prefix+"/user/auth", a.PostToAuth).Methods("POST")
}

//Пытаемся отконфигурировать наше храниелище (storage API)
func (a *API) configureStorageField() error {
	storageObj := storage.New(a.config.Storage)
	//Пытаемся установить соединение, если невозможно - возвращаем ошибку
	if err := storageObj.Open(); err != nil {
		return err
	}
	a.storage = storageObj
	return nil
}
