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
	"e-commerce-api/internal/models"
	"fmt"
	"sync"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetByID(id int) (*models.Product, error)
	Create(p *models.Product) error
	Update(id int, p *models.Product) error
	Delete(id int) error
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
	return nil, fmt.Errorf("product not found")
}

func (r *InMemoryArrayRepository) Create(product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	product.ID = r.nextID
	r.products = append(r.products, *product)
	r.nextID++
	return nil

}

func (r *InMemoryArrayRepository) Update(id int, product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, items := range r.products {
		if id == items.ID {
			product.ID = id
			r.products[i] = *product
			return nil
		}
	}
	return fmt.Errorf("no item matched")
}

func (r *InMemoryArrayRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, items := range r.products {
		if items.ID == id {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("no item matched")

}
