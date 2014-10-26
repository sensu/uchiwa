package uchiwa

func isStringInSlice(a string) bool {
	for _, b := range tmpResults.Subscriptions {
		if b == a {
			return true
		}
	}
	return false
}

// BuildSubscriptions builds a slice of every client subscriptions
func BuildSubscriptions() {
	for _, c := range tmpResults.Clients {
		m := c.(map[string]interface{})
		s := m["subscriptions"].([]interface{})
		//fmt.Println(m["subscriptions"])
		for _, e := range s {
			if !isStringInSlice(e.(string)) {
				tmpResults.Subscriptions = append(tmpResults.Subscriptions, e.(string))
			}
		}
	}
}
