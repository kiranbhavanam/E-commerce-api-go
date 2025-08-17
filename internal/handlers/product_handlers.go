package handlers

import (
	"e-commerce-api/internal/errors"
	"e-commerce-api/internal/models"
	"e-commerce-api/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}
	product, err := h.service.GetProductByID(id)
	if err != nil {
	h.handleServiceError(w,err)
	return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err = h.service.CreateProduct(&product)
	if err != nil {
		h.handleServiceError(w,err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)

}

func (h *ProductHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}
	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		return
	}
	err = h.service.UpdateProduct(id, &product)
	if err != nil {
		h.handleServiceError(w,err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteProduct(id)
	if err != nil {
		h.handleServiceError(w,err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


func (h *ProductHandler) handleServiceError(w http.ResponseWriter,err error){
	switch e:=err.(type){
	case *errors.ValidationError:
		//400 bad request
		http.Error(w,e.Error(),http.StatusBadRequest)
	case *errors.NotFoundError:
		//404 not found
		http.Error(w,e.Error(),http.StatusNotFound)
	case *errors.DuplicateError:
		//409 duplicate
		http.Error(w,e.Error(),http.StatusConflict)
	}
}