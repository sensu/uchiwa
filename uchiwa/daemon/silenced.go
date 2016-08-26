package daemon

import "fmt"

// buildSilenced constructs silenced objects for frontend consumption
func (d *Daemon) buildSilenced() {
	for _, c := range d.Data.Silenced {
		silence, ok := c.(map[string]interface{})
		if !ok {
			continue
		}

		dc, ok := silence["dc"].(string)
		if !ok {
			continue
		}

		id, ok := silence["id"].(string)
		if !ok {
			continue
		}

		silence["_id"] = fmt.Sprintf("%s/%s", dc, id)
	}
}
