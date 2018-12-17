package sensu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/sensu/uchiwa/uchiwa/helpers"
	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// These are the methods of the API struct that interface the more basic
// API methods below in order to return specifics data types and
// handle the pagination for example.

// getBytes returns the body of a GET request as []byte
func (api *API) getBytes(endpoint string) ([]byte, *http.Response, error) {
	return api.get(context.Background(), fmt.Sprintf("%s/%s", api.URL, endpoint))
}

// getSlice returns the body of a GET request as []interface{}
func (api *API) getSlice(ctx context.Context, endpoint string, limit int) ([]interface{}, error) {
	var offset int

	u, err := url.Parse(fmt.Sprintf("%s/%s", api.URL, endpoint))
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

	body, res, err := api.get(ctx, u.String())
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

			body, _, err := api.get(ctx, u.String())
			if err != nil {
				return nil, err
			}

			partialList, err := helpers.GetInterfacesFromBytes(body)
			if err != nil {
				return nil, fmt.Errorf("Could not parse the JSON-encoded response body: %v", err)
			}

			if len(partialList) == 0 {
				logger.Debugf("No additional elements found, exiting pagination for %s endpoint", endpoint)
				break
			}

			for _, e := range partialList {
				list = append(list, e)
			}
		}
	}

	return list, nil
}

// getSlice returns the body of a GET request as map[string]inteface{}
func (api *API) getMap(endpoint string) (map[string]interface{}, error) {

	body, _, err := api.get(context.Background(), fmt.Sprintf("%s/%s", api.URL, endpoint))
	if err != nil {
		return nil, err
	}
	return helpers.GetMapFromBytes(body)
}

// postPayload sends a POST request to a provided enpoint with the provided payload
func (api *API) postPayload(endpoint string, payload string) (map[string]interface{}, error) {

	url := fmt.Sprintf("%s/%s", api.URL, endpoint)

	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("%s\n\n", payload)))
	if err != nil {
		return nil, fmt.Errorf("Parsing error: %q returned: %v", url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))

	body, _, err := api.doRequest(req)
	if err != nil {
		return nil, err
	}

	return helpers.GetMapFromBytes(body)
}

// These are the methods of the API struct corresponding to
// their equivalent HTTP method (DELETE, GET and POST).

// delete performs a DELETE HTTP request to the provided endpoint
func (api *API) delete(endpoint string) error {
	url := fmt.Sprintf("%s/%s", api.URL, endpoint)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Parsing error: %q returned: %v", err, err)
	}

	if api.User != "" && api.Pass != "" {
		req.SetBasicAuth(api.User, api.Pass)
	}

	res, err := api.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("%v", res.Status)
	}

	return nil
}

// get returns the body of a GET HTTP request to a provided URL as []byte
func (api *API) get(ctx context.Context, u string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Parsing error: %q returned: %v", err, err)
	}
	req = req.WithContext(ctx)

	body, res, err := api.doRequest(req)
	if err != nil {
		return nil, nil, err
	}

	return body, res, nil
}

// post performs a POST HTTP request to a provided endpoint
func (api *API) post(endpoint string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", api.URL, endpoint)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	body, _, err := api.doRequest(req)
	if err != nil {
		return nil, err
	}

	return helpers.GetMapFromBytes(body)
}
