package service

import (
	"e-commerce-api/internal/models"
	"e-commerce-api/internal/repository"
	"errors"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product ID")
	}
	return s.repo.GetByID(id)
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	//business logic for product validation
	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(id int, product *models.Product) error {
	if id <= 0 {
		return errors.New("invalid product ID")
	}
	return s.repo.Update(id, product)
}

func (s *ProductService) DeleteProduct(id int) error {
	if id <= 0 {
		return errors.New("invalid product ID")
	}
	return s.repo.Delete(id)
}
