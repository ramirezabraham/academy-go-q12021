package usecase

import (
	"main/model"
	csvService "main/service/csv"
	httpService "main/service/http"
	"net/url"
)

// GetItems - Get Items from Service
func GetItems() ([]model.Item, *model.Error) {
	return csvService.GetItemsFromCSV()
}

// GetItem - Get Item from Service
func GetItem(itemID string) (model.Item, *model.Error) {
	return csvService.GetItemFromCSV(itemID)
}

func GetToken() (model.Token, *model.Error) {
	return httpService.CreateUrlToken()
}

func GetItemsAPI(token string, query url.Values) ([]model.ApiItem, *model.Error) {
	newItems, err := httpService.GetItemAPI(token, query)

	if err != nil {
		return nil, err
	}

	errorSaveItem := csvService.SaveItems(&newItems)

	if errorSaveItem != nil {
		return nil, errorSaveItem
	}

	return newItems, nil
}
