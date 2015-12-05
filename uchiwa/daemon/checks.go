package daemon

import "fmt"

// buildChecks constructs check objects for frontend consumption
func (d *Daemon) buildChecks() {
	for _, c := range d.Data.Checks {
		check, ok := c.(map[string]interface{})
		if !ok {
			continue
		}

		dc, ok := check["dc"].(string)
		if !ok {
			continue
		}

		name, ok := check["name"].(string)
		if !ok {
			continue
		}

		check["_id"] = fmt.Sprintf("%s/%s", dc, name)
	}
}
