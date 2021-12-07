package api

import (
	_ "github.com/gorilla/mux"
	storage "github.com/k-airos/6.ServerAndDB/storage"
	"github.com/sirupsen/logrus"
	"net/http"
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
	a.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello! This is my rest api"))
	})
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
