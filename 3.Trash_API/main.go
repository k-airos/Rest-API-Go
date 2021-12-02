package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var (
	port string = "8080"
	// Наша "база данных"
	db []Pizza
)

func init() {
	pizza1 := Pizza{
		ID:       1,
		Diameter: 22,
		Price:    500.50,
		Title:    "Pepperoni",
	}

	pizza2 := Pizza{
		ID:       2,
		Diameter: 25,
		Price:    650.50,
		Title:    "BBQ",
	}
	pizza3 := Pizza{
		ID:       3,
		Diameter: 22,
		Price:    300.50,
		Title:    "Cheese",
	}
	db = append(db, pizza1, pizza2, pizza3)
}

//наша модель
type Pizza struct {
	ID       int     `json:"id"`
	Diameter int     `json:"diameter"`
	Price    float64 `json:"price"`
	Title    string  `json:"title"`
}

// Вспомогательная функция для модели (модельный метод)
func FindPizzaById(id int) (Pizza, bool) {
	var pizza Pizza
	var found bool

	for _, p := range db {
		if p.ID == id {
			pizza = p
			found = true
			break
		}
	}
	return pizza, found
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func GetAllPizzas(writer http.ResponseWriter, request *http.Request) {
	//Прописываем хежеры
	writer.Header().Set("Content-Type", "application/json")
	log.Println("Get infos about all pizzas in database")
	writer.WriteHeader(200)            //Статус код для пиццы
	json.NewEncoder(writer).Encode(db) //сериализация + запись в writer
}

func GetPizzaById(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	//считаем id из строки запроса и конвертируем его в int
	vars := mux.Vars(request) // {"id":"12"} -- все параметры запроса в мапу
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := ErrorMessage{Message: "do not use ID not supported int casting"}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
	}
	log.Println("Trying to send to client to pizza with id #:", id)

	pizza, ok := FindPizzaById(id)
	if ok {
		writer.WriteHeader(200)
		json.NewEncoder(writer).Encode(pizza)
	} else {
		msg := ErrorMessage{Message: "pizza with that id doesnt exist in database"}
		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
	}
}

func main() {
	log.Println("Trying to start REST API pizza")
	//Инициализируем маршрутизатор
	router := mux.NewRouter()
	//1, Если на вход пришел запрос /pizzas
	router.HandleFunc("/pizzas", GetAllPizzas).Methods("GET")
	//2. Если на вход приходит запрос вида /pizza/{id}
	router.HandleFunc("/pizza/{id}", GetPizzaById).Methods("GET")

	log.Println("Router configured successfully! Let's go")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
