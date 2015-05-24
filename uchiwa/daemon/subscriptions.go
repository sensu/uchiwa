package daemon

func isStringInSlice(s string, slice []interface{}) bool {
	for _, element := range slice {
		if element == s {
			return true
		}
	}
	return false
}

// BuildSubscriptions builds a slice of every client subscriptions
func (d *Daemon) BuildSubscriptions() {
	for _, c := range d.Data.Clients {
		m := c.(map[string]interface{})
		subscriptions := m["subscriptions"].([]interface{})
		for _, element := range subscriptions {
			if !isStringInSlice(element.(string), subscriptions) {
				d.Data.Subscriptions = append(d.Data.Subscriptions, element.(string))
			}
		}
	}
}
