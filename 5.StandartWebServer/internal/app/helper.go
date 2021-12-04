package api

import (
	_ "github.com/gorilla/mux"
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
