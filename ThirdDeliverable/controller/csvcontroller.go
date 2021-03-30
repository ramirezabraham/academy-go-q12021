package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"main/model"
	"main/usecase"

	"github.com/gorilla/mux"
)

type ItemController struct {
	useCase usecase.NewItemsUseCase
}

type NewItemController interface {
	Index(w http.ResponseWriter, r *http.Request)
	GetItems(w http.ResponseWriter, r *http.Request)
	GetItem(w http.ResponseWriter, r *http.Request)
	GetToken(w http.ResponseWriter, r *http.Request)
	GetItemsAPI(w http.ResponseWriter, r *http.Request)
	GetItemsConcurrently(w http.ResponseWriter, r *http.Request)
}

func New(ic usecase.NewItemsUseCase) *ItemController {
	return &ItemController{ic}
}

// Index ..
func (ic *ItemController) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my Diablo 3 Items-API")
}

// GetItems - Get items from CSV
func (ic *ItemController) GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := ic.useCase.GetItems()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(err.Code)
		fmt.Fprint(w, `{"Message": "Error - %v "}`, err.Message)
	}
	json.NewEncoder(w).Encode(items)
}

// GetItem - Get an item from CSV
func (ic *ItemController) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID := vars["id"]
	item, err := ic.useCase.GetItem(itemID)

	if err != nil {
		w.WriteHeader(err.Code)
		fmt.Fprint(w, `{"Message": "Error - %v "}`, err.Message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// GetToken - Get a Token to consume API
func (ic *ItemController) GetToken(w http.ResponseWriter, r *http.Request) {
	token, err := ic.useCase.GetToken()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(err.Code)
		fmt.Fprint(w, `{"Message": "Error - %v "}`, err.Message)
	}
	json.NewEncoder(w).Encode(token)
}

// GetItemsAPI - Get items from external API
func (ic *ItemController) GetItemsAPI(w http.ResponseWriter, r *http.Request) {
	var token string
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		token = tokens[0]
		token = strings.TrimPrefix(token, "Bearer ")
	}

	// If the token is empty...
	if token == "" {
		// If we get here, the required token is missing
		http.Error(w, http.StatusText(http.StatusUnauthorized)+" - Missing token", http.StatusUnauthorized)
		return
	}

	//Check if region param is added to the url.
	regionParam, ok := r.URL.Query()["region"]

	if !ok || len(regionParam[0]) < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest)+" - Missing param 'region'", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := ic.useCase.GetItemsAPI(token, r.URL.Query())
	if err != nil {
		w.WriteHeader(err.Code)
		fmt.Fprintf(w, "There was some errors, please try again.")
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ic *ItemController) GetItemsConcurrently(w http.ResponseWriter, r *http.Request) {
	typeNumber := r.FormValue("type")
	w.Header().Set("Content-Type", "application/json")

	if typeNumber == "even" || typeNumber == "odd" {
		itemsFormValue := r.FormValue("items")
		itemsPerWorkerFormValue := r.FormValue("items_per_worker")

		items, _ := strconv.Atoi(itemsFormValue)
		itemsPerWorker, _ := strconv.Atoi(itemsPerWorkerFormValue)

		csvItems, _ := ic.useCase.GetConcurrentlyItems(typeNumber, items, itemsPerWorker)
		w.WriteHeader(http.StatusOK)
		concurrencyResponse := model.ConcurrencyResponse{TypeNumber: typeNumber, ItemNumber: items, ItemPerWorker: itemsPerWorker, Items: csvItems}
		json.NewEncoder(w).Encode(concurrencyResponse)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{ "message": "Only support "odd" or "even"" }`)
	}
}
