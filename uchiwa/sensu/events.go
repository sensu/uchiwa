package sensu

import (
	"context"
	"fmt"
)

// GetEvents returns a slice of all clients
func (s *Sensu) GetEvents(ctx context.Context) ([]interface{}, error) {
	return s.getSlice(ctx, "events", DefaultLimit)
}

// DeleteEvent delete an event
func (s *Sensu) DeleteEvent(check, client string) error {
	return s.delete(fmt.Sprintf("events/%s/%s", client, check))
}
