package filters

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// FilterAggregates based on role's datacenters
func FilterAggregates(data *[]interface{}, token *jwt.Token) []interface{} {
	return *data
}

// FilterChecks based on role's datacenters and subscriptions
func FilterChecks(data *[]interface{}, token *jwt.Token) []interface{} {
	return *data
}

// FilterClients based on role's datacenters and subscriptions
func FilterClients(data *[]interface{}, token *jwt.Token) []interface{} {
	return *data
}

// FilterDatacenters based on role's datacenters
func FilterDatacenters(data []*structs.Datacenter, token *jwt.Token) []*structs.Datacenter {
	return data
}

// FilterEvents based on role's datacenters and subscriptions
func FilterEvents(data *[]interface{}, token *jwt.Token) []interface{} {
	return *data
}

// FilterStashes based on role's datacenters
func FilterStashes(data *[]interface{}, token *jwt.Token) []interface{} {
	return *data
}

// FilterSubscriptions based on role's subscriptions
func FilterSubscriptions(data *[]string, token *jwt.Token) []string {
	return *data
}

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
