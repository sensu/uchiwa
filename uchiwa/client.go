package uchiwa

import (
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/daemon"
	"github.com/sensu/uchiwa/uchiwa/logger"
)

func (u *Uchiwa) buildClientHistory(id *string, history *[]interface{}, dc *string) {
	for _, h := range *history {
		m, ok := h.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert this client history to an interface: %+v", h)
			continue
		}

		// last_status comes in as a float64, so needs 0.0
		if m["last_status"] == 0.0 {
			lastResult := m["last_result"]
			lr, _ := lastResult.(map[string]interface{})
			m["output"] = lr["output"]
		} else {
			m["output"] = u.findOutput(id, m, dc)
		}

		check, ok := m["check"].(string)
		if !ok {
			logger.Warningf("Could not assert this check name to a string: %+v", m["check"])
			continue
		}

		m["model"] = findModel(check, *dc, u.Data.Checks)
		m["client"] = id
		m["dc"] = dc
		m["acknowledged"] = daemon.IsAcknowledged(*id, check, *dc, u.Data.Stashes)
	}
}

// DeleteClient send a DELETE request to the /clients/*client* endpoint in order to delete a client
func (u *Uchiwa) DeleteClient(id string, dc string) error {
	api, err := getAPI(u.Datacenters, dc)
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

func (u *Uchiwa) findClientInClients(id *string, dc *string) (map[string]interface{}, error) {
	for _, c := range u.Data.Clients {
		m, ok := c.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert this client to an interface %+v", c)
			continue
		}
		if m["name"] == *id && m["dc"] == *dc {
			return m, nil
		}
	}
	return nil, fmt.Errorf("Could not find client %s", *id)
}

func (u *Uchiwa) findOutput(id *string, h map[string]interface{}, dc *string) string {
	if h["last_status"] == 0 {
		return ""
	}

	for _, e := range u.Data.Events {
		// does the dc match?
		m, ok := e.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert this event to an interface %+v", e)
			continue
		}
		if m["dc"] != *dc {
			continue
		}

		// does the client match?
		c, ok := m["client"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert this client to an interface: %+v", c)
			continue
		}

		if c["name"] != *id {
			continue
		}

		// does the check match?
		k := m["check"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert this check to an interface: %+v", k)
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
func (u *Uchiwa) GetClient(id string, dc string) (map[string]interface{}, error) {
	api, err := getAPI(u.Datacenters, dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	// lock results while we gather client info
	u.Mu.Lock()
	defer u.Mu.Unlock()

	c, err := u.findClientInClients(&id, &dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	h, err := api.GetClientHistory(id)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	u.buildClientHistory(&id, &h, &dc)

	// add client history to client map for easy frontend consumption
	c["history"] = h

	return c, nil
}
