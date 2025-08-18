package repository

import (
	"e-commerce-api/internal/errors"
	"e-commerce-api/internal/models"
	"sync"
)

type InMemoryMapRepository struct{
	Products map[int]*models.Product
	nextID int
	mu sync.RWMutex
}
func NewInMemoryMapRepository() ProductRepository{
	return &InMemoryMapRepository{
		Products: make(map[int]*models.Product),
		nextID: 100,
		mu:sync.RWMutex{},
	}
}

func (r *InMemoryMapRepository) GetAll()([]models.Product,error){
	r.mu.RLock()
	defer r.mu.RUnlock()
	productArray:=make([]models.Product,0,len(r.Products))
	for _,p:=range r.Products{
		productArray=append(productArray, *p)
	}
	return productArray,nil
}
func (r *InMemoryMapRepository) GetByID(id int)(*models.Product,error){
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p,ok:=r.Products[id];ok{
		return p,nil
	}
	return nil,errors.NewNotFoundError(id,"product")
}
func (r *InMemoryMapRepository) Create(product *models.Product) error{
	r.mu.Lock()
	defer r.mu.Unlock()
	product.ID=r.nextID
	r.Products[r.nextID]=product
	r.nextID++
	return nil
}
func (r *InMemoryMapRepository) Update(id int,product *models.Product)error{
	r.mu.Lock()
	defer r.mu.Unlock()
	if _,ok:=r.Products[id];ok{
		product.ID=id
		r.Products[id]=product
		return nil
	}
	return errors.NewNotFoundError(id,"product")
}
func (r *InMemoryMapRepository) Delete(id int) error{
	r.mu.Lock()
	defer r.mu.Unlock()
	if _,ok:=r.Products[id];ok{
	delete(r.Products,id)
	return nil
	}
	return errors.NewNotFoundError(id,"product")
}

func (r *InMemoryMapRepository)ExistsByName(name string)bool{
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _,p:= range r.Products{
		if p.Name==name{
			return true
		}
	}
	return false
}