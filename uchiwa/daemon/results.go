package daemon

func (d *Daemon) buildResults() {
	for _, r := range d.Data.Results {

		// cast "result" interface to map[string]interface{}
		result := GetMapFromInterface(r)
		if result == nil {
			continue
		}

		clientName, ok := result["client"].(string)
		if !ok {
			continue
		}

		dcName, ok := result["dc"].(string)
		if !ok {
			continue
		}

		// cast interface to map[string]interface{}
		check := GetMapFromInterface(result["check"])
		if check == nil {
			continue
		}

		checkNameInterface := check["name"]
		if checkNameInterface == nil {
			continue
		}

		checkName, ok := checkNameInterface.(string)
		if !ok {
			continue
		}

		// Determine if the check is silenced
		check["acknowledged"] = IsAcknowledged(clientName, checkName, dcName, d.Data.Stashes)
	}
}
