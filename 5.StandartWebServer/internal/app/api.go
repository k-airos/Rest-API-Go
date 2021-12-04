package api

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

// API Base API server instance description
type API struct {
	//UNEXPORTED FIELD!
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

//API constructor: build base API instance
func New(config *Config) *API {
	return &API{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

//Start http server/configure loggers, router, database connection and etc...
func (a *API) Start() error {
	//Trying to configure logger
	if err := a.configureLoggerField(); err != nil {
		return err
	}
	//Подтверждение того, что логгер сконфигурирован
	a.logger.Info("Starting api server at port:", a.config.BindAddr)

	//Конфигурируем маршрутизатор
	a.configureRouterField()
	//На этапе валидного завершения стартуем http - сервер

	return http.ListenAndServe(a.config.BindAddr, a.router)
}
