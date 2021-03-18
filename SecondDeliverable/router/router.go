package router

import (
	"main/controller"

	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/token", controller.GetToken).Methods("GET")
	r.HandleFunc("/item/{id}", controller.GetItem).Methods("GET")
	r.HandleFunc("/allitems", controller.GetItems).Methods("GET")
	r.HandleFunc("/itemsAPI", controller.GetItemsAPI).Methods("GET")
	r.HandleFunc("/", controller.Index)

	return r

}
