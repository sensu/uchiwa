package daemon

// buildAggregates constructs aggregate elements compatible with
// the new Sensu 0.24.0 aggregates
func (d *Daemon) buildAggregates() {
	for _, e := range d.Data.Aggregates {
		m := e.(map[string]interface{})

		if m["name"] == nil && m["check"] != nil {
			m["name"] = m["check"]
			delete(m, "check")
		}
	}
}
