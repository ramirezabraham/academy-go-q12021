package usecase

import (
	"math"
	"net/url"
	"sync"

	"main/model"

	csvService "main/service/csv"
	httpService "main/service/http"
)

type ItemsUseCase struct {
	csvservice  csvService.NewCsvService
	httpservice httpService.NewHttpService
}

type NewItemsUseCase interface {
	GetItems() ([]model.Item, *model.Error)
	GetItem(itemID string) (model.Item, *model.Error)
	GetToken() (model.Token, *model.Error)
	GetItemsAPI(token string, query url.Values) ([]model.ApiItem, *model.Error)
	GetConcurrentlyItems(typeNumber string, items int, itemsPerWorker int) ([]model.Item, *model.Error)
}

func New(s csvService.NewCsvService, h httpService.NewHttpService) *ItemsUseCase {
	return &ItemsUseCase{s, h}
}

// ItemsPoolSize - Get Items PoolSize
func ItemsPoolSize(items int, itemsPerWorker int, totalItems int) int {
	var poolSize int
	if items%itemsPerWorker != 0 {
		poolSize = int(math.Ceil(float64(items) / float64(itemsPerWorker)))
	} else {
		poolSize = int(items / itemsPerWorker)
	}

	// If we go over the number of workers more than the half of number
	// of items it's gonna get into an infinit looop
	if poolSize > (totalItems / 2) {
		poolSize = totalItems / 2
	}
	return poolSize
}

// MaxItems - Get Max Items
func MaxItems(totalItems int) int {
	var maxItems int

	if totalItems%2 == 0 {
		maxItems = totalItems / 2
	} else {
		maxItems = totalItems/2 + 1
	}
	return maxItems
}

// GetItems - Get Items from Service
func (us *ItemsUseCase) GetItems() ([]model.Item, *model.Error) {
	return us.csvservice.GetItemsFromCSV()
}

// GetItem - Get Item from Service
func (us *ItemsUseCase) GetItem(itemID string) (model.Item, *model.Error) {
	return us.csvservice.GetItemFromCSV(itemID)
}

// GetToken - Get Token from Service
func (us *ItemsUseCase) GetToken() (model.Token, *model.Error) {
	return us.httpservice.CreateUrlToken()
}

// GetItemsAPI - Get Items from Service
func (us *ItemsUseCase) GetItemsAPI(token string, query url.Values) ([]model.ApiItem, *model.Error) {
	newItems, err := us.httpservice.GetItemAPI(token, query)

	if err != nil {
		return nil, err
	}

	errorSaveItem := us.csvservice.SaveItems(&newItems)

	if errorSaveItem != nil {
		return nil, errorSaveItem
	}

	return newItems, nil
}

// GetConcurrentlyItems - Get items Concurrently
func (us *ItemsUseCase) GetConcurrentlyItems(typeNumber string, items int, itemsPerWorker int) ([]model.Item, *model.Error) {
	csvItems, err := us.csvservice.GetItemsFromCSV()

	if err != nil {
		return nil, err
	}

	totalItems := len(csvItems)
	poolSize := ItemsPoolSize(items, itemsPerWorker, totalItems)
	maxItems := MaxItems(totalItems)

	//Creating Channels
	values := make(chan int)
	jobs := make(chan int, poolSize)
	shutdown := make(chan struct{})

	startIndex := 0
	var limit = int(math.Ceil(float64(totalItems) / float64(poolSize)))
	lastLimit := (totalItems % limit)

	var wg sync.WaitGroup
	wg.Add(poolSize)

	for i := 0; i < poolSize; i++ {
		go func(jobs <-chan int) {
			for {
				var id int
				var limitRecalculated int
				start := <-jobs

				// Recalculated the limit based on the number of jobs left.
				// lastLimit sometimes can be 0.
				if limit+start >= totalItems && lastLimit != 0 {
					limitRecalculated = start + lastLimit
				} else {
					limitRecalculated = start + limit
				}

				for j := start; j < limitRecalculated; j++ {
					id = csvItems[j].ItemNumber

					select {
					case values <- id:
					case <-shutdown:
						wg.Done()
						return
					}
				}
			}
		}(jobs)
	}

	for i := 0; i < poolSize; i++ {
		jobs <- startIndex
		startIndex += limit
	}
	close(jobs)

	var itemsFiltered []model.Item
	bucket := make(map[int]int, totalItems+1)
	for elem := range values {
		if typeNumber == "odd" {
			if elem%2 != 0 && bucket[elem] == 0 {
				itemRec := model.Item{ItemNumber: csvItems[elem-1].ItemNumber, ItemID: csvItems[elem-1].ItemID, ItemName: csvItems[elem-1].ItemName, ItemType: csvItems[elem-1].ItemType}
				itemsFiltered = append(itemsFiltered, itemRec)
				bucket[elem] = elem // Marking the selected items that has been added to the collection
			}
		} else if typeNumber == "even" {
			if elem%2 == 0 && bucket[elem] == 0 {
				itemRec := model.Item{ItemNumber: csvItems[elem-1].ItemNumber, ItemID: csvItems[elem-1].ItemID, ItemName: csvItems[elem-1].ItemName, ItemType: csvItems[elem-1].ItemType}
				itemsFiltered = append(itemsFiltered, itemRec)
				bucket[elem] = elem // Marking the selected items that has been added to the collection
			}
		}
		if len(itemsFiltered) >= items || len(itemsFiltered) >= maxItems {
			break // Finally when we reach the number of items or the possibly half that we can take, we break the loop
		}
	}

	// Send the signal to all goroutines to be finished
	close(shutdown)

	wg.Wait()

	return itemsFiltered, nil
}
