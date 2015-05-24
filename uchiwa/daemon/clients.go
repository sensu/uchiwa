package daemon

// buildClients constructs clients objects for frontend consumption
func (d *Daemon) buildClients() {
	for _, c := range d.Data.Clients {
		m := c.(map[string]interface{})

		if m["version"] == nil {
			m["version"] = "0.12.x"
		}

		d.findStatus(m)

		m["acknowledged"] = IsAcknowledged(m["name"].(string), "", m["dc"].(string), d.Data.Stashes)
	}
}
