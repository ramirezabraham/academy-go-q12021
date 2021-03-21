package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"main/usecase"

	"github.com/gorilla/mux"
)

// Index ..
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my Diablo 3 Items-API")
}

// GetItems - Get items from CSV
func GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := usecase.GetItems()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(err.Code)
		fmt.Fprint(w, "Something happened..%v", err.Message)
	}
	json.NewEncoder(w).Encode(items)
}

// GetItem - Get an item from CSV
func GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID := vars["id"]
	item, err := usecase.GetItem(itemID)

	if err != nil {
		w.WriteHeader(err.Code)
		fmt.Fprint(w, "Something Happened.. %v", err.Message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// GetToken - Get a Token to consume API
func GetToken(w http.ResponseWriter, r *http.Request) {
	token, err := usecase.GetToken()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(err.Code)
		fmt.Fprint(w, "Something happened..%v", err.Message)
	}
	json.NewEncoder(w).Encode(token)
}

// GetItemsAPI - Get items from external API
func GetItemsAPI(w http.ResponseWriter, r *http.Request) {
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
	response, err := usecase.GetItemsAPI(token, r.URL.Query())
	if err != nil {
		w.WriteHeader(err.Code)
		fmt.Fprintf(w, "There was some errors, please try again.")
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
