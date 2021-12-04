package handlers

import (
	"4.SemiTrash_API/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func GetBookById(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		log.Println("error while parsing happened:", err)
		writer.WriteHeader(400)
		msg := models.Message{Message: "do not use parameter ID as uncased to int type"}
		json.NewEncoder(writer).Encode(msg)
		return
	}

	book, ok := models.FindBookById(id)
	log.Println("Get book with id", id)
	if !ok {
		writer.WriteHeader(404)
		msg := models.Message{Message: "book with that id does not exists in database"}
		json.NewEncoder(writer).Encode(msg)
	} else {
		writer.WriteHeader(200)
		json.NewEncoder(writer).Encode(book)
	}
}

func CreateBook(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)
	log.Println("Creating new book .... ")
	var book models.Book

	err := json.NewDecoder(request.Body).Decode(&book)
	if err != nil {
		msg := models.Message{Message: "provided json file is invalid"}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	models.DB = append(models.DB, book)

	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(book)
}

func UpdateBookById(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)
	log.Println("updating book....")
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		log.Println("error while parsing happened:", err)
		writer.WriteHeader(400)
		msg := models.Message{Message: "do not use parameter ID as uncased to int type"}
		json.NewEncoder(writer).Encode(msg)
		return
	}
	oldBook, ok := models.FindBookById(id)
	var newBook models.Book
	if !ok {
		log.Println("Book not found in database with ID", id)
		writer.WriteHeader(404)
		msg := models.Message{Message: "book with that id does not exists in database"}
		json.NewEncoder(writer).Encode(msg)
		return
	}
	err = json.NewDecoder(request.Body).Decode(&newBook)
	if err != nil {
		msg := models.Message{Message: "provided json file is invalid"}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	for i := 0; i < len(models.DB); i++ {
		if oldBook.ID == models.DB[i].ID {
			models.DB[i] = newBook
			msg := models.Message{Message: "successfully updated"}
			json.NewEncoder(writer).Encode(msg)
			break
		}
	}
}

func DeleteBookById(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)
	log.Println("deleting book....")
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		log.Println("error while parsing happened:", err)
		writer.WriteHeader(400)
		msg := models.Message{Message: "do not use parameter ID as uncased to int type"}
		json.NewEncoder(writer).Encode(msg)
		return
	}
	book, ok := models.FindBookById(id)
	if !ok {
		log.Println("Book not found in database with ID", id)
		writer.WriteHeader(404)
		msg := models.Message{Message: "book with that id does not exists in database"}
		json.NewEncoder(writer).Encode(msg)
		return
	}

	for i := 0; i < len(models.DB); i++ {
		if book.ID == models.DB[i].ID {
			models.DB = append(models.DB[:i], models.DB[i+1:]...)
			break
		}
	}
	msg := models.Message{Message: "successfully deleted request item"}
	json.NewEncoder(writer).Encode(msg)
}
