package models

type Products struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Price float32 `json:"price"`
}