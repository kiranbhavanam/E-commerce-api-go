package repository

/*

import "ecommerce-api-go/internal/models"

type ProductRepository interface {
    GetAll() ([]models.Product, error)
    GetByID(id int) (*models.Product, error)
    Create(p *models.Product) error
    Update(id int, p *models.Product) error
    Delete(id int) error
}
*/

import (
	"e-commerce-api/internal/errors"
	"e-commerce-api/internal/models"
	"encoding/json"
	"os"
	"sync"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetByID(id int) (*models.Product, error)
	Create(p *models.Product) error
	Update(id int, p *models.Product) error
	Delete(id int) error
    ExistsByName(name string)bool
	LoadFromFile(name string)error
}

/*
Thread Safety: Structs encapsulate thread-safe state

Configuration: Dependency injection allows environment-specific setup

Testing: Mock implementations make tests fast and reliable

Maintenance: Change storage implementation without touching business logic

Performance: Efficient resource management and connection pooling

Scalability: Easy to swap implementations as requirements grow
*/

type InMemoryArrayRepository struct {
	products []models.Product
	mu       sync.RWMutex
	nextID   int
}

/*
func NewInMemoryArrayRepository() ProductRepository { //returns interface }
func NewInMemoryArrayRepositoryConcrete() *InMemoryArrayRepository {// returns concrete  }
*/
func NewInMemoryArrayRepository() ProductRepository {
	return &InMemoryArrayRepository{
		products: make([]models.Product, 0),
		mu:       sync.RWMutex{},
		nextID:   1,
	}
}

func (r *InMemoryArrayRepository) GetAll() ([]models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.products, nil
}
func (r *InMemoryArrayRepository) GetByID(id int) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i, product := range r.products {
		if product.ID == id {
			// return &product,nil dangerous:same address will be updated again and agian
			return &r.products[i], nil
		}
	}
	return nil, errors.NewNotFoundError(id,"product")
}

func (r *InMemoryArrayRepository) Create(product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	product.ID = r.nextID
	r.products = append(r.products, *product)
	r.nextID++
	return r.SaveToFile("products.json")

}

func (r *InMemoryArrayRepository) Update(id int, product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, items := range r.products {
		if id == items.ID {
			product.ID = id
			r.products[i] = *product
			return r.SaveToFile("products.json")
		}
	}
	return errors.NewNotFoundError(id,"product")
}

func (r *InMemoryArrayRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, items := range r.products {
		if items.ID == id {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return r.SaveToFile("products.json")
		}
	}
	return errors.NewNotFoundError(id,"product")

}

func (r *InMemoryArrayRepository) ExistsByName(name string)bool{
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    for _,product:=range r.products{
        if product.Name==name{
            return true
        }
    }
    return false
}
// func (r *InMemoryArrayRepository) ExistsByNameExcludingID()

//Method to save and load data to and from a file
func(r *InMemoryArrayRepository) SaveToFile(filename string)error{
	data,err:=json.MarshalIndent(r.products,"","")
	if err!=nil{
		return err
	}
	return os.WriteFile(filename,data,0644)
}

func (r *InMemoryArrayRepository) LoadFromFile(filename string)error{
	file,err:=os.Open(filename)
	if err!=nil{
		if os.IsNotExist(err){
			r.products=[]models.Product{}
			return nil
		}
		return err
	}
	defer file.Close()

	bytes,err:=os.ReadFile(filename)
	if err!=nil{
		return err
	}
	return json.Unmarshal(bytes,&r.products)
}