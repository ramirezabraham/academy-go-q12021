package csvservice

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"main/model"

	"github.com/spf13/viper"
)

type CsvService struct{}

type NewCsvService interface {
	GetItemsFromCSV() ([]model.Item, *model.Error)
	GetItemFromCSV(itemID string) (model.Item, *model.Error)
	SaveItems(newItems *[]model.ApiItem) *model.Error
}

func New() *CsvService {
	return &CsvService{}
}

// ReadCSV returns records, err
func ReadCSV(path string) ([][]string, *model.Error) {
	// Open the file
	csvfile, err := os.Open(path)
	if err != nil {
		e := model.Error{
			Code:    http.StatusInternalServerError,
			Message: "Couldn't open the csv file",
		}
		return nil, &e
	}
	defer csvfile.Close()

	// Parse the file
	reader := csv.NewReader(csvfile)
	reader.LazyQuotes = false
	reader.FieldsPerRecord = -1
	reader.Comma = ','
	reader.Comment = '#'
	records, err := reader.ReadAll()
	if err != nil {
		e := model.Error{
			Code:    http.StatusInternalServerError,
			Message: "Couldn't read the csv file",
		}
		return nil, &e
	}

	return records, nil
}

// GetConfig - Function for reading config.yaml file.
func GetConfig(key string) (string, *model.Error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	config_error := viper.ReadInConfig()
	if config_error != nil {
		e := model.Error{
			Code:    http.StatusInternalServerError,
			Message: "Couldn't open the config file",
		}
		return "", &e
	}
	keyData := viper.GetString(key)
	if keyData == "" {
		e := model.Error{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("key %v not found in config file", key),
		}
		return "", &e
	}
	return keyData, nil
}

// OpenAndWrite - Function for opening and writing in the csv file
func OpenAndWrite(path string) (*os.File, *model.Error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, &model.Error{
			Code:    http.StatusInternalServerError,
			Message: "There was an error opening the file",
		}
	}
	return f, nil
}

// AddLine - Add a new line of item data in csv file
func AddLine(f *os.File, newItems *[]model.ApiItem) *model.Error {

	w := csv.NewWriter(f)
	for _, item := range *newItems {
		var row []string
		row = append(row, item.ItemID)
		row = append(row, item.ItemName)
		row = append(row, item.ItemType)
		w.Write(row)
	}
	defer w.Flush()

	return nil
}

// GetItemsFromCSV - Gets all items from CSV
func (Csv *CsvService) GetItemsFromCSV() ([]model.Item, *model.Error) {
	var items []model.Item
	//Getting Config
	pathFile, errorPath := GetConfig("api.path")
	// Checking if pathFile exists
	if errorPath != nil {
		return nil, errorPath
	}
	records, err := ReadCSV(pathFile)
	// Iterate through the records
	for _, record := range records {
		if err != nil {
			err := model.Error{
				Code:    http.StatusInternalServerError,
				Message: "Error getting the records",
			}
			return nil, &err
		}
		id, intErr := strconv.Atoi(record[0])
		if intErr != nil {
			return nil, &model.Error{
				Message: intErr.Error(),
				Code:    http.StatusInternalServerError,
			}
		}
		itemRec := model.Item{ItemNumber: id, ItemID: record[1], ItemName: record[2], ItemType: record[3]}
		items = append(items, itemRec)
	}

	if err != nil {
		err := model.Error{
			Code:    http.StatusInternalServerError,
			Message: "Couldn't convert data to json",
		}
		return nil, &err
	}
	return items, nil
}

// GetItemFromCSV - Gets one item from CSV
func (Csv *CsvService) GetItemFromCSV(itemID string) (model.Item, *model.Error) {
	var item model.Item
	//Getting Config
	pathFile, errPath := GetConfig("api.path")
	// Checking if config exists
	if errPath != nil {
		return model.Item{}, errPath
	}
	records, err := ReadCSV(pathFile)
	// Iterate through the records
	for _, record := range records {
		if err != nil {
			return model.Item{}, &model.Error{
				Code:    http.StatusInternalServerError,
				Message: "Error getting the records",
			}
		}
		id, intErr := strconv.Atoi(record[0])
		if intErr != nil {
			return model.Item{}, &model.Error{
				Message: intErr.Error(),
				Code:    http.StatusInternalServerError,
			}
		}
		if itemID != "" {
			if itemID == record[1] {
				itemRec := model.Item{ItemNumber: id, ItemID: record[1], ItemName: record[2], ItemType: record[3]}
				item = itemRec
				break
			}
		} else {
			return model.Item{}, &model.Error{
				Code:    http.StatusAccepted,
				Message: "The Item does not exists",
			}
		}
	}
	if err != nil {
		err := model.Error{
			Code:    http.StatusInternalServerError,
			Message: "Couldn't convert data to json",
		}
		return model.Item{}, &err
	}
	return item, nil
}

// SaveItems - Function for saving items in the csv file
func (Csv *CsvService) SaveItems(newItems *[]model.ApiItem) *model.Error {
	//Getting Config
	pathFile, errPath := GetConfig("api.pathSave")
	if errPath != nil {
		return errPath
	}
	fileOpenAndWrite, _ := OpenAndWrite(pathFile) // Write

	errorLine := AddLine(fileOpenAndWrite, newItems)
	if errorLine != nil {
		return errorLine
	}

	return nil
}
