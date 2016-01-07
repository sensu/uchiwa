package sensu

import (
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/sensu/uchiwa/uchiwa/logger"
)

// These are the methods directly used by the public methods of the sensu
// package in order to handle the failover and load balancing between the APIs of a datacenter

func (s *Sensu) delete(endpoint string) error {
	apis := shuffle(s.APIs)

	var err error
	for i := 0; i < len(apis); i++ {
		logger.Infof("DELETE %s/%s", s.APIs[i].URL, endpoint)
		err = apis[i].delete(endpoint)
		if err == nil {
			return err
		}
		logger.Warningf("DELETE %s/%s returned: %v", s.APIs[i].URL, endpoint, err)
	}

	return err
}

func (s *Sensu) getBytes(endpoint string) ([]byte, *http.Response, error) {
	apis := shuffle(s.APIs)

	for i := 0; i < len(apis); i++ {
		logger.Debugf("GET %s/%s", s.APIs[i].URL, endpoint)
		bytes, res, err := apis[i].getBytes(endpoint)
		if err == nil {
			return bytes, res, err
		}
		logger.Warningf("GET %s/%s returned: %v", s.APIs[i].URL, endpoint, err)
	}

	return nil, nil, errors.New("")
}

func (s *Sensu) getSlice(endpoint string, limit int) ([]interface{}, error) {
	apis := shuffle(s.APIs)

	for i := 0; i < len(apis); i++ {
		logger.Debugf("GET %s/%s", s.APIs[i].URL, endpoint)
		slice, err := apis[i].getSlice(endpoint, limit)
		if err == nil {
			return slice, err
		}
		logger.Warningf("GET %s/%s returned: %v", s.APIs[i].URL, endpoint, err)
	}

	return nil, errors.New("")
}

func (s *Sensu) getMap(endpoint string) (map[string]interface{}, error) {
	apis := shuffle(s.APIs)

	for i := 0; i < len(apis); i++ {
		logger.Debugf("GET %s/%s", s.APIs[i].URL, endpoint)
		m, err := apis[i].getMap(endpoint)
		if err == nil {
			return m, err
		}
		logger.Warningf("GET %s/%s returned: %v", s.APIs[i].URL, endpoint, err)
	}

	return nil, errors.New("")
}

func (s *Sensu) postPayload(endpoint string, payload string) (map[string]interface{}, error) {
	apis := shuffle(s.APIs)

	for i := 0; i < len(apis); i++ {
		logger.Debugf("POST %s/%s", s.APIs[i].URL, endpoint)
		m, err := apis[i].postPayload(endpoint, payload)
		if err == nil {
			return m, err
		}
		logger.Warningf("POST %s/%s returned: %v", s.APIs[i].URL, endpoint, err)
	}

	return nil, errors.New("")
}

// shuffle the provided []API
func shuffle(apis []API) []API {
	rand.Seed(time.Now().UnixNano())
	for i := range apis {
		j := rand.Intn(i + 1)
		apis[i], apis[j] = apis[j], apis[i]
	}
	return apis
}
