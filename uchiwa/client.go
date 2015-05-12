package uchiwa

import (
	"fmt"

	"github.com/palourde/logger"
)

func buildClientHistory(id *string, history *[]interface{}, dc *string) {
	for _, h := range *history {
		m, ok := h.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert history interface %+v", h)
			continue
		}

		// last_status comes in as a float64, so needs 0.0
		if m["last_status"] == 0.0 {
			last_result := m["last_result"]
			lr, _ := last_result.(map[string]interface{})
			m["output"] = lr["output"]
		} else {
			m["output"] = findOutput(id, m, dc)
		}
		m["acknowledged"] = isAcknowledged(*id, m["check"].(string), *dc)
		m["model"] = findModel(m["check"].(string), *dc)
		m["client"] = id
		m["dc"] = dc
	}
}

// DeleteClient send a DELETE request to the /clients/*client* endpoint in order to delete a client
func DeleteClient(id string, dc string) error {
	api, err := getAPI(dc)
	if err != nil {
		logger.Warning(err)
		return err
	}

	err = api.DeleteClient(id)
	if err != nil {
		logger.Warning(err)
		return err
	}

	return nil
}

func findClientInClients(id *string, dc *string) (map[string]interface{}, error) {
	for _, c := range Results.Clients {
		m, ok := c.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert client interface %+v", c)
			continue
		}
		if m["name"] == *id && m["dc"] == *dc {
			return m, nil
		}
	}
	return nil, fmt.Errorf("Could not find client %s", *id)
}

func findOutput(id *string, h map[string]interface{}, dc *string) string {
	if h["last_status"] == 0 {
		return ""
	}

	for _, e := range Results.Events {
		// does the dc match?
		m, ok := e.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert event interface %+v", e)
			continue
		}
		if m["dc"] != *dc {
			continue
		}

		// does the client match?
		c, ok := m["client"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert event's client interface: %+v", c)
			continue
		}

		if c["name"] != *id {
			continue
		}

		// does the check match?
		k := m["check"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert event's check interface: %+v", k)
			continue
		}
		if k["name"] != h["check"] {
			continue
		}
		return k["output"].(string)
	}

	return ""
}

// GetClient retrieves client history from specified DC
func GetClient(id string, dc string) (map[string]interface{}, error) {
	api, err := getAPI(dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	// lock Results structure while we gather client info
	mutex.Lock()
	defer mutex.Unlock()

	c, err := findClientInClients(&id, &dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	h, err := api.GetClientHistory(id)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	buildClientHistory(&id, &h, &dc)

	// add client history to client map for easy frontend consumption
	c["history"] = h

	return c, nil
}
