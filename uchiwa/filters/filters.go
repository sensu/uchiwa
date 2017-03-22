package filters

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// Filters contains the different filtering methods based on the edition
type Filters interface {
	Aggregates(*[]interface{}, *jwt.Token) []interface{}
	Checks(*[]interface{}, *jwt.Token) []interface{}
	Client(interface{}, *jwt.Token) bool
	Clients(*[]interface{}, *jwt.Token) []interface{}
	Datacenters([]*structs.Datacenter, *jwt.Token) []*structs.Datacenter
	Events(*[]interface{}, *jwt.Token) []interface{}
	// NEED WORK
	GetRequest(string, *jwt.Token) bool
	Silenced(*[]interface{}, *jwt.Token) []interface{}
	Stashes(*[]interface{}, *jwt.Token) []interface{}
	Subscriptions(*[]structs.Subscription, *jwt.Token) []structs.Subscription
}

// Uchiwa represents an instance of the Filters interface for the community filters
type Uchiwa struct{}

// Aggregates filters based on role's datacenters
func (u *Uchiwa) Aggregates(data *[]interface{}, token *jwt.Token) []interface{} {
	aggregates := make([]interface{}, len(*data))
	copy(aggregates, *data)
	return aggregates
}

// Checks filters based on role's datacenters and subscriptions
func (u *Uchiwa) Checks(data *[]interface{}, token *jwt.Token) []interface{} {
	checks := make([]interface{}, len(*data))
	copy(checks, *data)
	return checks
}

// Client is a function that filters GET requests.
func (u *Uchiwa) Client(data interface{}, token *jwt.Token) bool {
	return true
}

// Clients filters based on role's datacenters and subscriptions
func (u *Uchiwa) Clients(data *[]interface{}, token *jwt.Token) []interface{} {
	clients := make([]interface{}, len(*data))
	copy(clients, *data)
	return clients
}

// Datacenters filters based on role's datacenters
func (u *Uchiwa) Datacenters(data []*structs.Datacenter, token *jwt.Token) []*structs.Datacenter {
	return data
}

// Events filters based on role's datacenters and subscriptions
func (u *Uchiwa) Events(data *[]interface{}, token *jwt.Token) []interface{} {
	events := make([]interface{}, len(*data))
	copy(events, *data)
	return events
}

// Silenced filters based on role's datacenters
func (u *Uchiwa) Silenced(data *[]interface{}, token *jwt.Token) []interface{} {
	silenced := make([]interface{}, len(*data))
	copy(silenced, *data)
	return silenced
}

// Stashes filters based on role's datacenters
func (u *Uchiwa) Stashes(data *[]interface{}, token *jwt.Token) []interface{} {
	stashes := make([]interface{}, len(*data))
	copy(stashes, *data)
	return stashes
}

// Subscriptions filters based on role's subscriptions
func (u *Uchiwa) Subscriptions(data *[]structs.Subscription, token *jwt.Token) []structs.Subscription {
	subscriptions := make([]structs.Subscription, len(*data))
	copy(subscriptions, *data)
	return subscriptions
}

// GetRequest is a function that filters GET requests.
func (u *Uchiwa) GetRequest(dc string, token *jwt.Token) bool {
	return false
}
