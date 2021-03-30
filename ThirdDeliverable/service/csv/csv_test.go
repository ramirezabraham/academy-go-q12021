package csvservice

import (
	"main/model"
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *CsvService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadCSV(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "File read",
			args: args{path: "../result.csv"},
			want: false,
		},
		{
			name: "File not read",
			args: args{path: "../result2.csv"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadCSV(tt.args.path)
			if (err != nil) != tt.want && tt.name == "File read" {
				t.Errorf("ReadCSV() got = %v, want %v", err, tt.want)
			}
			if (err != nil) != tt.want && tt.name == "File not read" {
				t.Errorf("ReadCSV() got1 = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
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
			got, got1 := GetConfig(tt.args.key)
			if got != tt.want {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetConfig() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
