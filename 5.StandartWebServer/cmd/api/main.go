package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	api "github.com/k-airos/5.StandartWebServer/internal/app"
	"log"
	"os"
)

var (
	format     string
	configPath string
)

func init() {
	flag.StringVar(&format, "format", ".toml", "format of config file")
	//Скажем, что наше приложение на этапе запуска будет получать путь до конфиг файла из внешнего мира
	flag.StringVar(&configPath, "path", "configs/api.toml", "path to config file in .toml format")
}

func main() {
	// В этот момент происходит инициализация переменной configPath
	flag.Parse()
	log.Println("It works")
	//server instance initialization
	config := api.NewConfig()
	// Теперь надо попробовать из .toml/.env, т.к там может быть новая информация
	if format == ".toml" {
		_, err := toml.DecodeFile(configPath, config) //Десериализируете содержимое .toml файла
		if err != nil {
			log.Println("Can not find config file. Using default values", err)
		}
	}
	if format == ".env" {
		err := godotenv.Load(configPath)
		if err != nil {
			log.Println("Can not find config file. Using default values", err)
		}
		config.BindAddr = os.Getenv("bind_addr")
		config.LoggerLevel = os.Getenv("logger_level")
	}
	server := api.New(config)

	//api server start
	//if err := server.Start(); err != nil {
	//	log.Fatal(err)
	//}
	//или такав
	log.Fatal(server.Start())
}
