package sensu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// NoLimit is used as a limit parameter
const NoLimit int = -1

// NoOffset is used as an offset parameter
const NoOffset int = -1

// get returns an array of byte which contains the response body
func (s *Sensu) get(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", s.URL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Parsing error: %q returned: %v", err, err)
	}
	res, err := s.doHTTP(req)
	if err != nil {
		return nil, fmt.Errorf("API call to %q returned: %v", url, err)
	}

	return res, nil
}

// getList Construct an API call and return the list of results
//  LIMITATION: the limit/offset is currently ignored and all results are sent back
func (s *Sensu) getList(endpoint string, limit int, offset int) ([]interface{}, error) {
	url := fmt.Sprintf("%s/%s", s.URL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("URL Parsing error: %q returned: %v", url, err)
	}
	res, err := s.doHTTP(req)
	if err != nil {
		return nil, fmt.Errorf("API call to %q returned: %v", url, err)
	}
	return s.doJSONArray(res)
}

// get returns an array of byte which contains the response body
func (s *Sensu) getMap(endpoint string) (map[string]interface{}, error) {
	body, err := s.get(endpoint)
	if err != nil {
		return nil, err
	}
	return s.doJSON(body)
}

func (s *Sensu) doHTTP(req *http.Request) ([]byte, error) {

	if s.User != "" && s.Pass != "" {
		req.SetBasicAuth(s.User, s.Pass)
	}

	res, err := s.Client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%v", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("Parsing response body returned: %v", err)
	}
	return body, nil
}

// doJsonArray Unmarshall JSON expecting an array
func (s *Sensu) doJSONArray(body []byte) ([]interface{}, error) {
	var results []interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("Parsing JSON-encoded response body: %v", err)
	}
	return results, nil
}

// doJsonArray Unmarshall JSON expecting a map
func (s *Sensu) doJSON(body []byte) (map[string]interface{}, error) {
	var results map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("Parsing JSON-encoded response body: %v", err)
	}
	return results, nil
}

// Post to endpoint
func (s *Sensu) post(endpoint string) (map[string]interface{}, error) {
	// Call a List with magic value of limit 0 and offset 0

	//ERROR GET LIST TODO deal with limit %d and offset %d", limit, offset

	url := fmt.Sprintf("%s/%s", s.URL, endpoint)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Parsing error: %q returned: %v", url, err)
	}

	res, err := s.doHTTP(req)
	if err != nil {
		return nil, fmt.Errorf("API call to %q returned: %v", url, err)
	}
	return s.doJSON(res)
}

// postPayload to endpoint
func (s *Sensu) postPayload(endpoint string, payload string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", s.URL, endpoint)

	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("%s\n\n", payload)))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))

	if err != nil {
		return nil, fmt.Errorf("Parsing error: %q returned: %v", url, err)
	}
	res, err := s.doHTTP(req)
	if err != nil {
		return nil, fmt.Errorf("API call to %q returned: %v", url, err)
	}
	return s.doJSON(res)
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

	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("%v", res.Status)
	}
	return nil
}
