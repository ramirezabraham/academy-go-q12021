package httpservice

import (
	"main/model"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func Test_httpService_GetConfig(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		h     *httpService
		args  args
		want  string
		want1 *model.Error
	}{
		{
			name:  "File Read",
			args:  args{key: "api.base_url"},
			want:  "https://us.api.blizzard.com/d3/data/item-type?",
			want1: nil,
		},
		{
			name: "File Read",
			args: args{key: "base_url"},
			want: "",
			want1: &model.Error{
				Code:    http.StatusInternalServerError,
				Message: "key base_url not found in config file",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.h.GetConfig(tt.args.key)
			if got != tt.want {
				t.Errorf("httpService.GetConfig() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("httpService.GetConfig() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_httpService_CreateUrlToken(t *testing.T) {
	tests := []struct {
		name  string
		h     *httpService
		want  model.Token
		want1 *model.Error
	}{
		{
			name: "Token retrieved",
			want: model.Token{
				AccessToken: "US0ugivbuJG3eJf5ZUzN3Zg09dWhSrIS7P", TokenType: "bearer", ExpiresIn: 86399,
			},
			want1: nil,
		},
		{
			name: "Token retrieved",
			want: model.Token{
				AccessToken: "USZGnXEqvm4rQnnx3ssmdUM7foBgTwObXH", TokenType: "bearer", ExpiresIn: 86399,
			},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &httpService{}
			got, got1 := h.CreateUrlToken()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpService.CreateUrlToken() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("httpService.CreateUrlToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_httpService_GetItemAPI(t *testing.T) {
	params := make(url.Values, 2)
	params["region"] = append(params["region"], "us")
	params["locale"] = append(params["locale"], "en_US")
	type args struct {
		token     string
		urlParams url.Values
	}
	tests := []struct {
		name  string
		h     *httpService
		args  args
		want  []model.ApiItem
		want1 *model.Error
	}{
		{
			name: "ITems retrieved",
			args: args{
				token:     "US0ugivbuJG3eJf5ZUzN3Zg09dWhSrIS7P",
				urlParams: params,
			},
			want: []model.ApiItem{
				{ItemID: "Shoulders_WitchDoctor", ItemName: "Shoulders", ItemType: "item-type/shoulderswitchdoctor"},
				{ItemID: "Gloves_Barbarian", ItemName: "Gloves", ItemType: "item-type/glovesbarbarian"},
				{ItemID: "Helm_DemonHunter", ItemName: "Helm", ItemType: "item-type/helmdemonhunter"},
				{ItemID: "Helm_Wizard", ItemName: "Helm", ItemType: "item-type/helmwizard"},
				{ItemID: "Shield", ItemName: "Shield", ItemType: "item-type/shield"},
			},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.h.GetItemAPI(tt.args.token, tt.args.urlParams)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpService.GetItemAPI() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("httpService.GetItemAPI() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
