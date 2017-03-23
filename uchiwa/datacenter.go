package uchiwa

import (
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/structs"
)

func (u *Uchiwa) Datacenter(name string) (*structs.Datacenter, error) {
	for _, dc := range u.Data.Dc {
		if dc.Name == name {
			return dc, nil
		}
	}

	return nil, fmt.Errorf("")
}
