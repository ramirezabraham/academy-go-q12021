package usecase

import (
	"main/model"
	csvService "main/service/csv"
	httpService "main/service/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		s csvService.NewCsvService
		h httpService.NewHttpService
	}
	tests := []struct {
		name string
		args args
		want *ItemsUseCase
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.s, tt.args.h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemsPoolSize(t *testing.T) {
	type args struct {
		items          int
		itemsPerWorker int
		totalItems     int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ItemsPoolSize(tt.args.items, tt.args.itemsPerWorker, tt.args.totalItems); got != tt.want {
				t.Errorf("ItemsPoolSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxItems(t *testing.T) {
	type args struct {
		totalItems int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxItems(tt.args.totalItems); got != tt.want {
				t.Errorf("MaxItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemsUseCase_GetItems(t *testing.T) {
	tests := []struct {
		name  string
		us    *ItemsUseCase
		want  []model.Item
		want1 *model.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.us.GetItems()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemsUseCase.GetItems() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ItemsUseCase.GetItems() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestItemsUseCase_GetItem(t *testing.T) {
	type args struct {
		itemID string
	}
	tests := []struct {
		name  string
		us    *ItemsUseCase
		args  args
		want  model.Item
		want1 *model.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.us.GetItem(tt.args.itemID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemsUseCase.GetItem() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ItemsUseCase.GetItem() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestItemsUseCase_GetToken(t *testing.T) {
	tests := []struct {
		name  string
		us    *ItemsUseCase
		want  model.Token
		want1 *model.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.us.GetToken()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemsUseCase.GetToken() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ItemsUseCase.GetToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestItemsUseCase_GetItemsAPI(t *testing.T) {
	type args struct {
		token string
		query url.Values
	}
	tests := []struct {
		name  string
		us    *ItemsUseCase
		args  args
		want  []model.ApiItem
		want1 *model.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.us.GetItemsAPI(tt.args.token, tt.args.query)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemsUseCase.GetItemsAPI() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ItemsUseCase.GetItemsAPI() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestItemsUseCase_GetConcurrentlyItems(t *testing.T) {
	type args struct {
		typeNumber     string
		items          int
		itemsPerWorker int
	}
	tests := []struct {
		name  string
		us    *ItemsUseCase
		args  args
		want  []model.Item
		want1 *model.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.us.GetConcurrentlyItems(tt.args.typeNumber, tt.args.items, tt.args.itemsPerWorker)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemsUseCase.GetConcurrentlyItems() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ItemsUseCase.GetConcurrentlyItems() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
