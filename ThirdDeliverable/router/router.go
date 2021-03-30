package router

import (
	"main/controller"

	"github.com/gorilla/mux"
)

type Router struct {
	itemController controller.NewItemController
}

type IRouter interface {
	InitRouter() *mux.Router
}

func New(ctrl controller.NewItemController) *Router {
	return &Router{ctrl}
}

func (router *Router) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", router.itemController.Index)
	r.HandleFunc("/allitems", router.itemController.GetItems).Methods("GET")
	r.HandleFunc("/item/{id}", router.itemController.GetItem).Methods("GET")
	r.HandleFunc("/token", router.itemController.GetToken).Methods("GET")
	r.HandleFunc("/itemsAPI", router.itemController.GetItemsAPI).Methods("GET")
	r.HandleFunc("/concurrency", router.itemController.GetItemsConcurrently).Methods("GET")
	return r

}
