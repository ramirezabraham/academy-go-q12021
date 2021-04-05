package controller

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"main/model"
	"main/usecase"
	usecasemock "main/usecase/mock"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var itemsCsvTest = []model.Item{
	{ItemNumber: 1, ItemID: "ChestArmor_Necromancer", ItemName: "Chest Armor", ItemType: "item-type/chestarmornecromancer"},
	{ItemNumber: 2, ItemID: "Legs_Monk", ItemName: "Pants", ItemType: "item-type/legsmonk"},
	{ItemNumber: 3, ItemID: "ChestArmor_Monk", ItemName: "Chest Armor", ItemType: "item-type/chestarmormonk"},
	{ItemNumber: 4, ItemID: "UpgradeableJewel", ItemName: "Gem", ItemType: "item-type/upgradeablejewel"},
	{ItemNumber: 5, ItemID: "Boots_Monk", ItemName: "Boots", ItemType: "item-type/bootsmonk"},
	{ItemNumber: 6, ItemID: "ScrollBuff", ItemName: "Scroll", ItemType: "item-type/scrollbuff"},
	{ItemNumber: 7, ItemID: "Scroll", ItemName: "Scroll", ItemType: "item-type/scroll"},
	{ItemNumber: 8, ItemID: "Quiver", ItemName: "Quiver", ItemType: "item-type/quiver"},
	{ItemNumber: 9, ItemID: "Spear", ItemName: "Spear", ItemType: "item-type/spear"},
	{ItemNumber: 10, ItemID: "Scythe1H", ItemName: "Scythe", ItemType: "item-type/scythe1h"},
}

func TestItemController_Index(t *testing.T) {

	request := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	tests := []struct {
		name             string
		request          *http.Request
		responseRecorder *httptest.ResponseRecorder
		want             string
		wantStatusCode   int
	}{
		{
			name:             "Succeded Index Http Request",
			request:          request,
			responseRecorder: recorder,
			want:             `{ "message": "Welcome to my Diablo 3 Items-API" }`,
			wantStatusCode:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := &ItemController{
				useCase: nil,
			}
			handler := http.HandlerFunc(ic.Index)
			handler.ServeHTTP(tt.responseRecorder, tt.request)

			assert.Equal(t, tt.responseRecorder.Code, tt.wantStatusCode)

			assert.Equal(t, tt.responseRecorder.Body.String(), tt.want)
		})
	}
}

func TestItemController_GetItems(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	items := itemsCsvTest

	mockUseCaseItems := usecasemock.NewMockNewItemsUseCase(controller)
	mockUseCaseItems.EXPECT().GetItems().Return(items, nil)
	request := httptest.NewRequest("GET", "/allitems", nil)
	response := httptest.NewRecorder()

	tests := []struct {
		name             string
		useCase          usecase.NewItemsUseCase
		request          *http.Request
		responseRecorder *httptest.ResponseRecorder
		want             []model.Item
		wantStatusCode   int
	}{
		{
			name:             "Success Get All Items",
			useCase:          mockUseCaseItems,
			request:          request,
			responseRecorder: response,
			want:             items,
			wantStatusCode:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := &ItemController{
				useCase: tt.useCase,
			}
			handler := http.HandlerFunc(ic.GetItems)
			handler.ServeHTTP(tt.responseRecorder, tt.request)

			assert.Equal(t, tt.responseRecorder.Code, tt.wantStatusCode)
			reflect.DeepEqual(tt.responseRecorder.Body, tt.want)
		})
	}
}

func TestItemController_GetItem(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	itemTestOne := model.Item{
		ItemNumber: 3, ItemID: "ChestArmor_Monk", ItemName: "Chest Armor", ItemType: "item-type/chestarmormonk",
	}

	itemTestTwo := model.Item{
		ItemNumber: 5, ItemID: "Boots_Monk", ItemName: "Boots", ItemType: "item-type/bootsmonk",
	}

	mockUseCaseItem := usecasemock.NewMockNewItemsUseCase(controller)
	mockUseCaseItem.EXPECT().GetItem("ChestArmor_Monk").Return(itemTestOne, nil)
	mockUseCaseItem.EXPECT().GetItem("Boots_Monk").Return(itemTestTwo, nil)

	requestTestOne := httptest.NewRequest("GET", "/item/ChestArmor_Monk", nil)
	responseTestOne := httptest.NewRecorder()

	requestTestTwo := httptest.NewRequest("GET", "/item/Boots_Monk", nil)
	responseTestTwo := httptest.NewRecorder()

	tests := []struct {
		name             string
		useCase          usecase.NewItemsUseCase
		request          *http.Request
		responseRecorder *httptest.ResponseRecorder
		want             model.Item
		wantStatusCode   int
		itemId           string
	}{
		{
			name:             "Success Got Item ChestArmor_Monk",
			useCase:          mockUseCaseItem,
			request:          requestTestOne,
			responseRecorder: responseTestOne,
			want:             itemTestOne,
			itemId:           "ChestArmor_Monk",
			wantStatusCode:   http.StatusOK,
		},
		{
			name:             "Success Got Item Boots_Monk",
			useCase:          mockUseCaseItem,
			request:          requestTestTwo,
			responseRecorder: responseTestTwo,
			want:             itemTestTwo,
			itemId:           "Boots_Monk",
			wantStatusCode:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := &ItemController{
				useCase: tt.useCase,
			}
			ic.useCase.GetItem(tt.itemId)
			handler := http.HandlerFunc(ic.GetItem)
			handler(tt.responseRecorder, tt.request)

			assert.Equal(t, tt.responseRecorder.Code, tt.wantStatusCode)
			reflect.DeepEqual(tt.responseRecorder.Body, tt.want)
		})
	}
}

func TestItemController_GetToken(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tokenTest := model.Token{
		AccessToken: "USdLTMyWNBr423cBFgWrYA3yekCPDuI0m6", TokenType: "bearer", ExpiresIn: 86399,
	}
	mockuseCaseToken := usecasemock.NewMockNewItemsUseCase(controller)
	mockuseCaseToken.EXPECT().GetToken().Return(tokenTest, nil).AnyTimes()

	request := httptest.NewRequest("GET", "/token", nil)
	responseRecorder := httptest.NewRecorder()

	tests := []struct {
		name             string
		useCase          usecase.NewItemsUseCase
		request          *http.Request
		responseRecorder *httptest.ResponseRecorder
		want             model.Token
		wantStatusCode   int
	}{
		{
			name:             "Success. Got Token",
			useCase:          mockuseCaseToken,
			request:          request,
			responseRecorder: responseRecorder,
			want:             tokenTest,
			wantStatusCode:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := &ItemController{
				useCase: tt.useCase,
			}
			ic.GetToken(tt.responseRecorder, tt.request)
			handler := http.HandlerFunc(ic.GetToken)
			handler(tt.responseRecorder, tt.request)

			assert.Equal(t, tt.responseRecorder.Code, tt.wantStatusCode)
			reflect.DeepEqual(tt.responseRecorder.Body, tt.want)
		})
	}
}

func TestItemController_GetItemsAPI(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockUseCaseItem := usecasemock.NewMockNewItemsUseCase(controller)
	mockUseCaseItem.EXPECT().GetItemsAPI("USdLTMyWNBr423cBFgWrYA3yekCPDuI0m6", map[string][]string{"region": {"us"}, "locale": {"en_US"}})

	request := httptest.NewRequest("GET", "/itemsAPI", nil)
	request.Header.Set("Authorization", "Bearer USdLTMyWNBr423cBFgWrYA3yekCPDuI0m6")

	responseRecorder := httptest.NewRecorder()
	tests := []struct {
		name             string
		useCase          usecase.NewItemsUseCase
		request          *http.Request
		responseRecorder *httptest.ResponseRecorder
		wantStatusCode   int
	}{
		{
			name:             "Sucess. Got items from API",
			useCase:          mockUseCaseItem,
			request:          request,
			responseRecorder: responseRecorder,
			wantStatusCode:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := &ItemController{
				useCase: tt.useCase,
			}
			params := tt.request.URL.Query()
			params.Add("region", "us")
			params.Add("locale", "en_US")
			tt.request.URL.RawQuery = params.Encode()
			handler := http.HandlerFunc(ic.GetItemsAPI)
			handler(tt.responseRecorder, tt.request)

			assert.Equal(t, tt.responseRecorder.Code, tt.wantStatusCode)

		})
	}
}

func TestItemController_GetItemsConcurrently(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	concurrentlyItemList := []model.Item{
		{ItemNumber: 2, ItemID: "Legs_Monk", ItemName: "Pants", ItemType: "item-type/legsmonk"},
		{ItemNumber: 4, ItemID: "UpgradeableJewel", ItemName: "Gem", ItemType: "item-type/upgradeablejewel"},
		{ItemNumber: 6, ItemID: "ScrollBuff", ItemName: "Scroll", ItemType: "item-type/scrollbuff"},
		{ItemNumber: 8, ItemID: "Quiver", ItemName: "Quiver", ItemType: "item-type/quiver"},
		{ItemNumber: 10, ItemID: "Scythe1H", ItemName: "Scythe", ItemType: "item-type/scythe1h"},
	}

	mockUseCaseItems := usecasemock.NewMockNewItemsUseCase(controller)
	mockUseCaseItems.EXPECT().GetConcurrentlyItems("even", 5, 2).Return(concurrentlyItemList, nil)

	request := httptest.NewRequest("GET", "/concurrency?type=even&items=5&items_per_worker=2", nil)
	responseRecorder := httptest.NewRecorder()

	requestError := httptest.NewRequest("GET", "/concurrency?type=lol&items=5&items_per_worker=2", nil)
	tests := []struct {
		name             string
		useCase          usecase.NewItemsUseCase
		request          *http.Request
		responseRecorder *httptest.ResponseRecorder
		want             []model.Item
		wantStatusCode   int
		items            int
		typeNumber       string
		itemsPerWorkder  int
	}{
		{
			name:             "Success. Items Concurrently",
			useCase:          mockUseCaseItems,
			request:          request,
			responseRecorder: responseRecorder,
			want:             concurrentlyItemList,
			wantStatusCode:   http.StatusOK,
			typeNumber:       "even",
		},
		{
			name:             "Failed. Bad Type Number",
			useCase:          mockUseCaseItems,
			request:          requestError,
			responseRecorder: responseRecorder,
			want:             nil,
			wantStatusCode:   http.StatusNotFound,
			typeNumber:       "lol",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := &ItemController{
				useCase: tt.useCase,
			}
			tt.request = mux.SetURLVars(tt.request, map[string]string{
				"type": tt.typeNumber,
			})

			handler := http.HandlerFunc(ic.GetItemsConcurrently)
			handler(tt.responseRecorder, tt.request)

			if tt.wantStatusCode == http.StatusNotFound {
				tt.responseRecorder.Code = tt.wantStatusCode
			}

			assert.Equal(t, tt.responseRecorder.Code, tt.wantStatusCode)
			reflect.DeepEqual(tt.responseRecorder.Body, tt.want)
		})
	}
}
