package uchiwa

// BuildClients constructs clients objects for frontend consumption
func BuildClients() {
	for _, c := range tmpResults.Clients {
		m := c.(map[string]interface{})

		if m["version"] == nil {
			m["version"] = "0.12.x"
		}

		findStatus(m)

		m["acknowledged"] = isAcknowledged(m["name"].(string), "", m["dc"].(string))
	}
}
