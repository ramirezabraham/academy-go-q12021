package model

//Item Struct
type Item struct {
	ItemNumber int    `json:"number"`
	ItemID     string `json:"id"`
	ItemName   string `json:"name"`
	ItemType   string `json:"type"`
}

//ApiItem Struct
type ApiItem struct {
	ItemID   string `json:"id"`
	ItemName string `json:"name"`
	ItemType string `json:"path"`
}

//Token Struct
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

//ConcurrencyResponse Struct
type ConcurrencyResponse struct {
	TypeNumber    string `json:"type"`
	ItemNumber    int    `json:"item_number"`
	ItemPerWorker int    `json:"items_per_worker"`
	Items         []Item `json:"itemsResult"`
}
