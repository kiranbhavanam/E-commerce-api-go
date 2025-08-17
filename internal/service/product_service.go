package service

import (
	"e-commerce-api/internal/errors"
	"e-commerce-api/internal/models"
	"e-commerce-api/internal/repository"
	"strings"
	"unicode/utf8"
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
	if err:=s.validateID(id);err!=nil{
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	//business logic for product validation
	if err:=s.validateProduct(product);err!=nil{
		return err
	}
	if s.repo.ExistsByName(product.Name){
		return errors.NewDuplicateError("product",product.Name)
	}
	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(id int, product *models.Product) error {
	if err:=s.validateID(id);err!=nil{
		return err
	}
	return s.repo.Update(id, product)
}

func (s *ProductService) DeleteProduct(id int) error {
	if err:=s.validateID(id);err!=nil{
		return  err
	}
	return s.repo.Delete(id)
}

func (s *ProductService) validateID(id int) error{
	if id <= 0 {
		return  errors.NewValidationError("id","must be a positive number")
	}
	return nil
}
func (s *ProductService) validateProduct(product *models.Product) error {
    if product == nil {
        return errors.NewValidationError("product", "cannot be nil")
    }
    
    // Validate name
    if err := s.validateProductName(product.Name); err != nil {
        return err
    }
    
    // Validate price
    if err := s.validateProductPrice(product.Price); err != nil {
        return err
    }
    
    return nil
}

func (s *ProductService) validateProductName(name string) error {
    // Trim whitespace for validation
    trimmedName := strings.TrimSpace(name)
    
    // Check if empty
    if trimmedName == "" {
        return errors.NewValidationError("name", "is required and cannot be empty")
    }
    
    // Check minimum length
    if utf8.RuneCountInString(trimmedName) < 2 {
        return errors.NewValidationError("name", "must be at least 2 characters long")
    }
    
    // Check maximum length
    if utf8.RuneCountInString(trimmedName) > 100 {
        return errors.NewValidationError("name", "cannot exceed 100 characters")
    }
    
    // Check for invalid characters (basic example)
    if strings.Contains(trimmedName, "<") || strings.Contains(trimmedName, ">") {
        return errors.NewValidationError("name", "cannot contain HTML tags")
    }
    
    return nil
}

func (s *ProductService) validateProductPrice(price float32) error {
    // Check if price is positive
    if price <= 0 {
        return errors.NewValidationError("price", "must be greater than 0")
    }
    
    // Check reasonable upper limit (example: $1 million)
    if price > 1000000 {
        return errors.NewValidationError("price", "cannot exceed $1,000,000")
    }
    
    // Check for reasonable precision (max 2 decimal places)
    if price*100 != float32(int(price*100)) {
        return errors.NewValidationError("price", "cannot have more than 2 decimal places")
    }
    
    return nil
}