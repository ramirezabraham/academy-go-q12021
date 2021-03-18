package model

//Item Struct
type Item struct {
	ItemID   string `json:"id"`
	ItemName string `json:"name"`
	ItemType string `json:"type"`
}

//Item Struct
type ApiItem struct {
	ItemID   string `json:"id"`
	ItemName string `json:"name"`
	ItemType string `json:"path"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
