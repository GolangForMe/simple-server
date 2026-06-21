package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func init() {
	db = append(db, Pizza{
		Id:       1,
		Title:    "Peperoni",
		Diameter: 20,
		Price:    300.40,
	}, Pizza{
		Id:       2,
		Title:    "Memeroni",
		Diameter: 40,
		Price:    600.20,
	}, Pizza{
		Id:       3,
		Title:    "Teteroni",
		Diameter: 30,
		Price:    400.30,
	})
}

// Basic model
type Pizza struct {
	Id       int     `json:"id"`
	Title    string  `json:"title"`
	Diameter int     `json:"diameter"`
	Price    float64 `json:"price"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

// DB
var (
	db []Pizza
)

const (
	port string = "8080"
)

// Модельный метод
func FindPizzaById(id int) (Pizza, bool) {
	for _, p := range db {
		if p.Id == id {
			return p, true
		}
	}
	return Pizza{}, false
}

func GetAllPizza(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// логирование в консоль (stdout) с использованием встроенного пакета log
	log.Println("Get info about all pizza in database")
	w.WriteHeader(200) // status code
	// NewEncoder - подготавливает w для записи
	// Encode - делает маршалинг (сериализует слайс db - в слайс байт)
	json.NewEncoder(w).Encode(db)
}

func GetPizzaById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// vars - это map, у которой могут быть поля вида { "id": "12" }
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("Try to use wrong type of id in request")
		msg := ErrorMessage{
			Message: "Wrong id",
		}
		w.WriteHeader(400) // bad request
		json.NewEncoder(w).Encode(msg)

		return
	}

	log.Printf("Get info about pizza with id %d\n", id)
	pizza, ok := FindPizzaById(id)

	// pizza не найдена
	if !ok {
		log.Printf("Pizza with id %d not found in database\n", id)
		msg := ErrorMessage{
			Message: "Pizza not found",
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(msg)

		return
	}

	// pizza найдена
	if ok {
		log.Printf("Pizza with id %d found in database\n", id)
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(pizza)

		return
	}
}

func RemovePizzaById(w http.ResponseWriter, r *http.Request) {
	//
}

func CreateNewPizza(w http.ResponseWriter, r *http.Request) {
	//
}

func main() {
	// Создаем экзмепляр маршрутизатора
	r := mux.NewRouter()
	// Матчинг пар запрос-исполнитель
	r.HandleFunc("/pizza", GetAllPizza).Methods("GET")
	r.HandleFunc("/pizza/{id}", GetPizzaById).Methods("GET") // {id} - параметрический запрос
	r.HandleFunc("/pizza", CreateNewPizza).Methods("POST")
	r.HandleFunc("/pizza/{id}", RemovePizzaById).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+port, r))
}
