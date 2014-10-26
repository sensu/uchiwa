package uchiwa

type health struct {
	Uchiwa string                 `json:"uchiwa"`
	Sensu  map[string]interface{} `json:"sensu"`
}

// Health contains the health of Uchiwa and every Sensu API
var Health = health{Uchiwa: "ok", Sensu: make(map[string]interface{})}
