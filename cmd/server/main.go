package main

import (
	"encoding/json"
	"fmt"
	products "e-commerce-api/internal/models"
	"net/http"
)

func main(){
	http.HandleFunc("/",homeHandler)
	http.HandleFunc("/about",aboutHandler)
	http.HandleFunc("/product",productHandler)
	fmt.Println("listening on 8081")
	http.ListenAndServe(":8081",nil)
}
func homeHandler(w http.ResponseWriter,r *http.Request){
	// return "hello world"
	w.Write([]byte("hello world"))
}
func aboutHandler(w http.ResponseWriter,r *http.Request){
	product :=products.Products{
		Name:"football",
		Price: 25.2,
	}
	data,_:=json.Marshal(product)
	w.Write(data)
}
func productHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case "GET":
		handleProductGet(w,r)
	case "POST":
		handleProductPost(w,r)
	default:
		http.Error(w,"Method not supported",http.StatusMethodNotAllowed)
	}
}

var cart []products.Products
func handleProductGet(w http.ResponseWriter,r *http.Request){
	// data:=products.Products{
	// 	Name:"bat",Price: 400,
	// }
	// data1:="kiran"
	// marshalledData,_:=json.Marshal(data)
	// w.Write(marshalledData)
	w.Header().Set("Content-type","application/json")
	if len(cart)==0{
		w.Write([]byte("no data at the moment"))
	}else{

		json.NewEncoder(w).Encode(cart)
	}
}
func handleProductPost(w http.ResponseWriter,r *http.Request){

	var item products.Products
	err:=json.NewDecoder(r.Body).Decode(&item)
	if err!=nil{
		fmt.Print("error parsing the body")
	}
	cart = append(cart, item)
	// json.NewEncoder(w).Encode(cart)
	w.Write([]byte("data inserted successfully"))
}