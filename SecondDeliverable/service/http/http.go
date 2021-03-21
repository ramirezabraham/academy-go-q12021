package httpservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/model"
	"net/http"
	"net/url"
	"os/exec"

	"github.com/spf13/viper"
)

// GetConfig - Get configuration.
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
			Message: fmt.Sprintf(" key %v not found in config file", key),
		}
		return "", &e
	}
	return keyData, nil
}

// CreateUrlToken - Creates a Token
func CreateUrlToken() (model.Token, *model.Error) {
	//Getting Config
	clientId, errorKey := GetConfig("api.client_id")
	if errorKey != nil {
		return model.Token{}, errorKey
	}
	clientSecret, errorSecret := GetConfig("api.client_secret")
	if errorSecret != nil {
		return model.Token{}, errorSecret
	}
	url := fmt.Sprintf("%v:%v", clientId, clientSecret)
	curl := exec.Command("curl", "-u", url, "-d", "grant_type=client_credentials", "https://us.battle.net/oauth/token")
	out, err := curl.Output()
	if err != nil {
		e := model.Error{
			Code:    http.StatusInternalServerError,
			Message: "Couldn't execute the curl command",
		}
		return model.Token{}, &e
	}
	var token model.Token
	errorJson := json.Unmarshal(out, &token)
	if errorJson != nil {
		e := model.Error{
			Code:    http.StatusInternalServerError,
			Message: "Couldn't execute the curl command",
		}
		return model.Token{}, &e
	}
	return token, nil
}

// GetItemAPI - Get items from external API.
func GetItemAPI(token string, urlParams url.Values) ([]model.ApiItem, *model.Error) {
	//Getting config.
	url, errUrl := GetConfig("api.base_url")
	if errUrl != nil {
		return nil, errUrl
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, &model.Error{
			Code:    http.StatusInternalServerError,
			Message: "Request couldn't worked.",
		}
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	queryParams := req.URL.Query()
	queryParams.Add("region", urlParams["region"][0])
	queryParams.Add("locale", urlParams["locale"][0])
	req.URL.RawQuery = queryParams.Encode()
	// fmt.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		return nil, &model.Error{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong executing your request",
		}
	}
	defer resp.Body.Close()
	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &model.Error{
			Code:    http.StatusInternalServerError,
			Message: "Encountered some issues with the request response",
		}
	}

	var response []model.ApiItem
	json.Unmarshal(bodybytes, &response)

	newItems := response

	return newItems, nil
}
