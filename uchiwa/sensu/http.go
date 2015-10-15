package sensu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/sensu/uchiwa/uchiwa/helpers"
	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// get returns an array of byte which contains the response body
func (s *Sensu) get(u string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Parsing error: %q returned: %v", err, err)
	}

	body, res, err := s.doHTTP(req)
	if err != nil {
		return nil, nil, fmt.Errorf("API call to %q returned: %v", u, err)
	}

	return body, res, nil
}

// getList constructs an API call and returns the list of results
func (s *Sensu) getList(endpoint string, limit int) ([]interface{}, error) {
	var offset int

	u, err := url.Parse(fmt.Sprintf("%s/%s", s.URL, endpoint))
	if err != nil {
		return nil, fmt.Errorf("Could not parse the URL '%s': %v", u.String(), err)
	}

	// Add limit & offset parameters when required
	if limit != -1 {
		params := u.Query()
		params.Add("limit", strconv.Itoa(limit))
		params.Add("offset", strconv.Itoa(offset))
		u.RawQuery = params.Encode()
	}

	body, res, err := s.get(u.String())
	if err != nil {
		return nil, err
	}

	list, err := helpers.GetInterfacesFromBytes(body)
	if err != nil {
		return nil, fmt.Errorf("Could not parse the JSON-encoded response body: %v", err)
	}

	// Verify if the endpoint supports pagination
	if limit != -1 && res.Header.Get("X-Pagination") != "" {
		var xPagination structs.XPagination

		err = json.Unmarshal([]byte(res.Header.Get("X-Pagination")), &xPagination)
		if err != nil {
			logger.Warning(err)
		}

		for len(list) < xPagination.Total {
			offset += limit
			params := u.Query()
			params.Set("offset", strconv.Itoa(offset))
			u.RawQuery = params.Encode()

			body, _, err := s.get(u.String())
			if err != nil {
				return nil, err
			}

			partialList, err := helpers.GetInterfacesFromBytes(body)
			if err != nil {
				return nil, fmt.Errorf("Could not parse the JSON-encoded response body: %v", err)
			}

			for _, e := range partialList {
				list = append(list, e)
			}
		}
	}

	return list, nil
}

// get returns an array of byte which contains the response body
func (s *Sensu) getMap(endpoint string) (map[string]interface{}, error) {
	body, _, err := s.get(fmt.Sprintf("%s/%s", s.URL, endpoint))
	if err != nil {
		return nil, err
	}
	return helpers.GetMapFromBytes(body)
}

func (s *Sensu) doHTTP(req *http.Request) ([]byte, *http.Response, error) {
	if s.User != "" && s.Pass != "" {
		req.SetBasicAuth(s.User, s.Pass)
	}

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, nil, fmt.Errorf("%v", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("Parsing response body returned: %v", err)
	}

	return body, res, nil
}

// Post to endpoint
func (s *Sensu) post(endpoint string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", s.URL, endpoint)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Parsing error: %q returned: %v", url, err)
	}

	body, _, err := s.doHTTP(req)
	if err != nil {
		return nil, fmt.Errorf("API call to %q returned: %v", url, err)
	}

	return helpers.GetMapFromBytes(body)
}

// postPayload to endpoint
func (s *Sensu) postPayload(endpoint string, payload string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", s.URL, endpoint)

	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("%s\n\n", payload)))
	if err != nil {
		return nil, fmt.Errorf("Parsing error: %q returned: %v", url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))

	body, _, err := s.doHTTP(req)
	if err != nil {
		return nil, fmt.Errorf("API call to %q returned: %v", url, err)
	}

	return helpers.GetMapFromBytes(body)
}

// Delete resource
func (s *Sensu) delete(endpoint string) error {
	url := fmt.Sprintf("%s/%s", s.URL, endpoint)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Parsing error: %q returned: %v", err, err)
	}

	if s.User != "" && s.Pass != "" {
		req.SetBasicAuth(s.User, s.Pass)
	}

	res, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("API call to %q returned: %v", url, err)
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("%v", res.Status)
	}

	return nil
}
