package main

import (
	"e-commerce-api/internal/handlers"
	"e-commerce-api/internal/repository"
	"e-commerce-api/internal/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	//Dependency creation : creating repo,service,handler instances
	repo := repository.NewInMemoryArrayRepository()
	err:=repo.LoadFromFile("products.json")
	if err!=nil{
		log.Fatal("Failed to load products:",err)
	}
	service := service.NewProductService(repo)
	handler := handlers.NewProductHandler(service)
	//HTTP routing
	router := mux.NewRouter()
	router.HandleFunc("/products", handler.GetAllHandler).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}", handler.GetByIDHandler).Methods("GET")
	router.HandleFunc("/products", handler.CreateHandler).Methods("POST")
	router.HandleFunc("/products/{id:[0-9]+}", handler.UpdateHandler).Methods("PUT")
	router.HandleFunc("/products/{id:[0-9]+}", handler.DeleteHandler).Methods("DELETE")
	//start the server.
	fmt.Println("Starting the server:")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
