package handler

import (
	"fmt"
	"net/http"
)

type Order struct{}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create order")
}
func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("list orders")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete order")
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get by id")
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update by id")
}
