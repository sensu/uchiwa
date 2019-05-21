package sensu

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/sensu/uchiwa/uchiwa/logger"
)

// These are the methods directly used by the public methods of the sensu
// package in order to handle the failover and load balancing between the APIs of a datacenter

func (s *Sensu) delete(endpoint string) error {
	var err error
	apis, err := s.healthyAPIs()
	if err != nil {
		return err
	}
	shuffledRange := shuffle(makeRange(len(apis)))

	for _, i := range shuffledRange {
		logger.Debugf("DELETE %s/%s (%s)", s.APIs[i].URL, endpoint, s.Name)
		res, err := apis[i].delete(endpoint)
		if err == nil || res.StatusCode < 500 {
			return nil
		}
		s.APIs[i].Healthy = false
		logger.Debugf("DELETE %s/%s (%s) returned: %v", s.APIs[i].URL, endpoint, s.Name, err)
	}

	return err
}

func (s *Sensu) getBytes(endpoint string) ([]byte, *http.Response, error) {
	var bytes []byte
	var err error
	var res *http.Response
	apis, err := s.healthyAPIs()
	if err != nil {
		return nil, nil, err
	}
	shuffledRange := shuffle(makeRange(len(apis)))

	for _, i := range shuffledRange {
		bytes, res, err = s.getBytesFromAPI(apis[i], endpoint)
		if err == nil {
			return bytes, res, nil
		}
	}

	return nil, nil, err
}

func (s *Sensu) getBytesFromAPI(api *API, endpoint string) ([]byte, *http.Response, error) {
	logger.Debugf("GET %s/%s (%s)", api.URL, endpoint, s.Name)
	bytes, res, err := api.getBytes(endpoint)
	if res == nil || res.StatusCode >= 500 {
		api.Healthy = false
	} else {
		api.Healthy = true
	}
	if err == nil && res.StatusCode < 400 {
		return bytes, res, err
	}
	logger.Debugf("GET %s/%s (%s) returned: %v", api.URL, endpoint, s.Name, err)
	return bytes, res, err
}

func (s *Sensu) getSlice(ctx context.Context, endpoint string, limit int) ([]interface{}, error) {
	var err error
	var slice []interface{}
	var res *http.Response
	apis, err := s.healthyAPIs()
	if err != nil {
		return nil, err
	}
	shuffledRange := shuffle(makeRange(len(apis)))

	for _, i := range shuffledRange {
		logger.Debugf("GET %s/%s (%s)", s.APIs[i].URL, endpoint, s.Name)
		slice, res, err = apis[i].getSlice(ctx, endpoint, limit)
		if err == nil || res != nil && res.StatusCode >= 400 && res.StatusCode < 500 {
			return slice, err
		}
		s.APIs[i].Healthy = false
		logger.Debugf("GET %s/%s (%s) returned: %v", s.APIs[i].URL, endpoint, s.Name, err)
	}

	return nil, err
}

func (s *Sensu) getMap(endpoint string) (map[string]interface{}, error) {
	var err error
	var m map[string]interface{}
	var res *http.Response
	apis, err := s.healthyAPIs()
	if err != nil {
		return nil, err
	}
	shuffledRange := shuffle(makeRange(len(apis)))

	for _, i := range shuffledRange {
		logger.Debugf("GET %s/%s (%s)", apis[i].URL, endpoint, s.Name)
		m, res, err = apis[i].getMap(endpoint)
		if res == nil || res.StatusCode >= 500 {
			apis[i].Healthy = false
		} else {
			apis[i].Healthy = true
		}
		if err == nil && res.StatusCode < 400 {
			return m, err
		}
		logger.Debugf("GET %s/%s (%s) returned: %v", apis[i].URL, endpoint, s.Name, err)
	}

	return nil, err
}

func (s *Sensu) getMapFromAPI(api *API, endpoint string) (map[string]interface{}, error) {
	m, res, err := api.getMap(endpoint)
	if err != nil || res != nil && res.StatusCode >= 500 {
		api.Healthy = false
	}
	return m, err
}

func (s *Sensu) postPayload(endpoint string, payload string) (map[string]interface{}, error) {
	var err error
	var m map[string]interface{}
	apis, err := s.healthyAPIs()
	if err != nil {
		return nil, err
	}
	shuffledRange := shuffle(makeRange(len(apis)))

	for _, i := range shuffledRange {
		logger.Debugf("POST %s/%s (%s)", s.APIs[i].URL, endpoint, s.Name)
		m, err = apis[i].postPayload(endpoint, payload)
		if err == nil {
			return m, err
		}
		s.APIs[i].Healthy = false
		logger.Debugf("POST %s/%s (%s) returned: %v", s.APIs[i].URL, endpoint, s.Name, err)
	}

	return nil, err
}

// healthyAPIs returns a list of APIs with Healthy set to true or returns an error when there are
// no healthy APIs
func (s *Sensu) healthyAPIs() ([]*API, error) {
	var healthyAPIs []*API
	for i := range s.APIs {
		api := &s.APIs[i]
		if api.Healthy {
			healthyAPIs = append(healthyAPIs, api)
		}
	}
	if len(healthyAPIs) < 1 {
		return healthyAPIs, fmt.Errorf("No healthy APIs available for datacenter: %s", s.Name)
	}
	return healthyAPIs, nil
}

// makeRange returns an []int range from 0 to length - 1
func makeRange(length int) []int {
	a := make([]int, length)
	for i := range a {
		a[i] = i
	}
	return a
}

// shuffle the provided []int
func shuffle(intRange []int) []int {
	shuffledRange := make([]int, len(intRange))
	copy(shuffledRange, intRange)
	rand.Seed(time.Now().UnixNano())
	for i := range shuffledRange {
		j := rand.Intn(i + 1)
		shuffledRange[i], shuffledRange[j] = shuffledRange[j], shuffledRange[i]
	}
	return shuffledRange
}
