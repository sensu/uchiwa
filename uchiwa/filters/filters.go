package filters

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// GetRequest is a function that filters GET requests.
func GetRequest(dc string, token *jwt.Token) bool {
	return false
}

// PostRequest is a function that filters POST requests.
func PostRequest(token *jwt.Token, data *interface{}) bool {
	return false
}

// SensuData is a function that filters Sensu Data.
func SensuData(token *jwt.Token, data *structs.Data) *structs.Data {
	return data
}
