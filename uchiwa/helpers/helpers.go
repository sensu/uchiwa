package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// return true if String in Slice, false otherwise
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// BuildClientsMetrics builds the metrics for the events
func BuildClientsMetrics(clients *[]interface{}) *structs.StatusMetrics {
	metrics := structs.StatusMetrics{}

	metrics.Total = len(*clients)
	for _, c := range *clients {
		client := c.(map[string]interface{})

		silenced, ok := client["silenced"].(bool)
		if ok { // Do not ignore the client if we don't have a silenced attribute
			if silenced {
				metrics.Silenced++
				continue
			}
		}

		status, ok := client["status"].(int)
		if !ok {
			logger.Warningf("Could not assert the status to an int: %+v", client["status"])
			continue
		}

		if status == 2.0 {
			metrics.Critical++
			continue
		} else if status == 1.0 {
			metrics.Warning++
			continue
		} else if status == 0.0 {
			metrics.Healthy++
			continue
		}
		metrics.Unknown++
	}

	return &metrics
}

// BuildEventsMetrics builds the metrics for the events
func BuildEventsMetrics(events *[]interface{}) *structs.StatusMetrics {
	metrics := structs.StatusMetrics{}

	metrics.Total = len(*events)

	for _, e := range *events {
		event := e.(map[string]interface{})

		silenced, ok := event["silenced"].(bool)
		if ok { // Do not ignore the event if we don't have a silenced attribute
			if silenced {
				metrics.Silenced++
				continue
			}
		}

		check, ok := event["check"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert this check to an interface: %+v", event["check"])
			continue
		}

		status, ok := check["status"].(float64)
		if !ok {
			logger.Warningf("Could not assert this status to a flot64: %+v", check["status"])
			continue
		}

		if status == 2.0 {
			metrics.Critical++
			continue
		} else if status == 1.0 {
			metrics.Warning++
			continue
		}
		metrics.Unknown++
	}

	return &metrics
}

// GetBoolFromInterface ...
func GetBoolFromInterface(i interface{}) (bool, error) {
	if i == nil {
		logger.Debug("The interface is nil")
		return false, errors.New("The interface is nil")
	}

	b, ok := i.(bool)
	if !ok {
		logger.Debugf("Could not assert to a boolean the interface: %+v", i)
		return false, errors.New("Could not assert to a boolean the interface")
	}

	return b, nil
}

// GetEvent returns an event associated to a specific check
func GetEvent(check, client, dc string, events *[]interface{}) (map[string]interface{}, error) {
	if check == "" || client == "" || dc == "" || len(*events) == 0 {
		return nil, errors.New("No parameters should be empty")
	}

	for _, e := range *events {
		event, ok := e.(map[string]interface{})
		if !ok {
			continue
		}

		if event["dc"] == nil || event["dc"] != dc {
			continue
		}

		c, ok := event["client"].(map[string]interface{})
		if !ok {
			if event["client"] == nil || event["client"] != client {
				continue
			}
		} else if c["name"] == nil || c["name"] != client {
			continue
		}

		k, ok := event["check"].(map[string]interface{})
		if !ok {
			if event["check"] == nil || event["check"] != check {
				continue
			} else {
				return map[string]interface{}{"check": event["check"], "client": event["client"], "occurrences": event["occurrences"], "output": event["output"], "status": event["status"]}, nil
			}
		} else if k["name"] == nil || k["name"] != check {
			continue
		}

		if event["action"] != nil {
			k["action"] = event["action"]
		}
		if event["occurrences"] != nil {
			k["occurrences"] = event["occurrences"]
		}

		return k, nil
	}

	return nil, errors.New("No event found")
}

// GetInterfacesFromBytes returns a slice of interfaces from a slice of byte
func GetInterfacesFromBytes(bytes []byte) ([]interface{}, error) {
	var interfaces []interface{}
	if err := json.Unmarshal(bytes, &interfaces); err != nil {
		return nil, err
	}
	return interfaces, nil
}

// GetMapFromBytes returns a map from a slice of byte
func GetMapFromBytes(bytes []byte) (map[string]interface{}, error) {
	var m map[string]interface{}

	if len(bytes) == 0 {
		return m, nil
	}

	if err := json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetMapFromInterface returns a map from an interface
func GetMapFromInterface(i interface{}) map[string]interface{} {
	m, ok := i.(map[string]interface{})
	if !ok {
		logger.Debugf("Could not assert to a map the interface: %+v", i)
		return nil
	}

	return m
}

// GetIP returns the real user IP address
func GetIP(r *http.Request) string {
	if xForwardedFor := r.Header.Get("X-FORWARDED-FOR"); len(xForwardedFor) > 0 {
		return xForwardedFor
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// InterfaceToSlice takes a slice of type interface{} and returns a slice of interface
func InterfaceToSlice(slice interface{}) ([]interface{}, error) {
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		return nil, fmt.Errorf("The interface provided is not a slice: %+v", slice)
	}

	result := make([]interface{}, value.Len())

	for i := 0; i < value.Len(); i++ {
		result[i] = value.Index(i).Interface()
	}

	return result, nil
}

// InterfaceToString takes a slice of interface{} a slice of string
func InterfaceToString(i []interface{}) []string {
	var result []string

	for _, value := range i {
		v, ok := value.(string)
		if ok {
			result = append(result, v)
		}
	}

	return result
}

// IsCheckSilenced determines whether a check for a particular client is silenced.
// Returns true if the check is silenced and a slice of silence entries IDs
func IsCheckSilenced(check, client map[string]interface{}, dc string, silenced []interface{}) (bool, []string) {
	var isSilenced bool
	var commonSubscriptions, isSilencedBy, subscribers, subscriptions []string

	if dc == "" || len(silenced) == 0 {
		return false, isSilencedBy
	}

	checkName, ok := check["name"].(string)
	if !ok {
		return false, isSilencedBy
	}

	clientName, ok := client["name"].(string)
	if !ok {
		clientName = ""
	}

	// Retrieve the check subscribers
	if check["subscribers"] != nil {
		s, ok := check["subscribers"].([]interface{})
		if !ok {
			return false, isSilencedBy
		}
		subscribers = InterfaceToString(s)
	}

	// Retrieve the client subscriptions
	if client["subscriptions"] != nil {
		s, ok := client["subscriptions"].([]interface{})
		if !ok {
			return false, isSilencedBy
		}
		subscriptions = InterfaceToString(s)
	}

	// Get the subscriptions both check and client have in common
	for _, subscriber := range subscribers {
		if IsStringInArray(subscriber, subscriptions) {
			commonSubscriptions = append(commonSubscriptions, subscriber)
		}
	}

	for _, silence := range silenced {
		m, ok := silence.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert this silence entry to a map: %+v", silence)
			continue
		}

		if m["dc"] != dc {
			continue
		}

		// Ignore silenced entries that have not begun yet
		if m["begin"] != nil {
			begin := time.Unix(int64(m["begin"].(float64)), 0)
			now := time.Now()
			if now.Before(begin) {
				continue
			}
		}

		// Check (e.g. *:check_cpu)
		if m["id"] == fmt.Sprintf("*:%s", checkName) {
			isSilenced = true
			isSilencedBy = append(isSilencedBy, m["id"].(string))
			continue
		}

		// Client subscription (e.g. client:foo:* )
		if m["id"] == fmt.Sprintf("client:%s:*", clientName) {
			isSilenced = true
			isSilencedBy = append(isSilencedBy, m["id"].(string))
			continue
		}

		// Client's check subscription (e.g. client:foo:check_cpu )
		if m["id"] == fmt.Sprintf("client:%s:%s", clientName, checkName) {
			isSilenced = true
			isSilencedBy = append(isSilencedBy, m["id"].(string))
			continue
		}

		for _, subscription := range commonSubscriptions {
			// Subscription (e.g. load-balancer:* )
			if m["id"] == fmt.Sprintf("%s:*", subscription) {
				isSilenced = true
				isSilencedBy = append(isSilencedBy, m["id"].(string))
				continue
			}

			// Subscription' check (e.g. load-balancer:check_cpu)
			if m["id"] == fmt.Sprintf("%s:%s", subscription, checkName) {
				isSilenced = true
				isSilencedBy = append(isSilencedBy, m["id"].(string))
				continue
			}
		}

	}

	return isSilenced, isSilencedBy
}

// IsClientSilenced determines whether a client is silenced.
// Returns true if the client is silenced.
func IsClientSilenced(client, dc string, silenced []interface{}) bool {
	if client == "" || dc == "" || len(silenced) == 0 {
		return false
	}

	for _, silence := range silenced {
		m, ok := silence.(map[string]interface{})
		if !ok {
			continue
		}

		if m["dc"] == dc && m["id"] == fmt.Sprintf("client:%s:*", client) {
			return true
		}
	}

	return false
}

// IsStringInArray searches 'array' for 'item' string
// Returns true 'item' is a value of 'array'
func IsStringInArray(item string, array []string) bool {
	if item == "" || len(array) == 0 {
		return false
	}

	for _, element := range array {
		if element == item {
			return true
		}
	}

	return false
}

// RandomString generates a random string of the provided length
func RandomString(length int) string {
	if length == 0 {
		length = 32
	}

	char := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UTC().UnixNano())
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = char[rand.Intn(len(char)-1)]
	}
	return string(buf)
}
