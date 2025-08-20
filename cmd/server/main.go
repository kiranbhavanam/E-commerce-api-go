package main

import (
	"e-commerce-api/internal/config"
	"e-commerce-api/internal/handlers"
	"e-commerce-api/internal/repository"
	"e-commerce-api/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	//Dependency creation : creating repo,service,handler instances
	repo := createRepository(os.Getenv("ENV"))
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

func createRepository(repoType string) repository.ProductRepository{
	switch repoType{
	case "array":
		repo:=repository.NewInMemoryArrayRepository()
		if fileRepo,ok:=repo.(*repository.InMemoryArrayRepository);ok{
			fileRepo.LoadFromFile("products.json")
		}
		return repo
	case "map":
		repo:=repository.NewInMemoryMapRepository()
		return repo
	case "db":
		dbconfig:=config.LoadDBConfig()
		repo,err:=repository.NewPostgresRepository(dbconfig.ConnectionString())
		if err!=nil{
			log.Fatal("error while loading postgres repo ",err)
		}
		return repo
	default:
		return repository.NewInMemoryMapRepository()
	}
}