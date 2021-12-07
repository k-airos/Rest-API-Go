package api

import (
	"github.com/gorilla/mux"
	"github.com/k-airos/7.ServerAndDB2/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

// API Base API server instance description
type API struct {
	//UNEXPORTED FIELD!
	config *Config
	logger *logrus.Logger
	router *mux.Router
	//Добавление поля для работы с хранилищем
	storage *storage.Storage
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
func (api *API) Start() error {
	//Trying to configure logger
	if err := api.configureLoggerField(); err != nil {
		return err
	}
	//Подтверждение того, что логгер сконфигурирован
	api.logger.Info("Starting api server at port:", api.config.BindAddr)

	//Конфигурируем маршрутизатор
	api.configureRouterField()
	//конфигурируем хранилище
	if err := api.configureStorageField(); err != nil {
		return err
	}
	//На этапе валидного завершения стартуем http - сервер

	return http.ListenAndServe(api.config.BindAddr, api.router)
}
