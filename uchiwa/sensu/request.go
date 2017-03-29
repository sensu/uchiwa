package sensu

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ...
func (api *API) doRequest(req *http.Request) ([]byte, *http.Response, error) {
	if api.User != "" && api.Pass != "" {
		req.SetBasicAuth(api.User, api.Pass)
	}

	res, err := api.Client.Do(req)
	if err != nil {
		status, ok := err.(*url.Error)
		if !ok {
			return nil, nil, fmt.Errorf("Unexpected error, got %T, wanted *url.Error", err)
		}
		return nil, nil, status.Err
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
