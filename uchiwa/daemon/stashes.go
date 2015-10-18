package daemon

import "fmt"

// buildStashes constructs stashes objects for frontend consumption
func (d *Daemon) buildStashes() {
	for _, c := range d.Data.Stashes {
		stash, ok := c.(map[string]interface{})
		if !ok {
			continue
		}

		dc, ok := stash["dc"].(string)
		if !ok {
			continue
		}

		path, ok := stash["path"].(string)
		if !ok {
			continue
		}

		stash["_id"] = fmt.Sprintf("%s/%s", dc, path)
	}
}
