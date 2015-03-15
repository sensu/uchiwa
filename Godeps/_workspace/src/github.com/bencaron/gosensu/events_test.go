package sensu

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetEvents(t *testing.T) {
	sensu := getSensuTester()
	assert := assert.New(t)
	events, err := sensu.GetEvents()
	assert.Nil(err, fmt.Sprintf("GetEvents returned an error: %s", err))
	assert.NotNil(events, "GetEvents returned nil!")
}

// func TestResolveEvents(t *testing.T) {
// 	sensu := getSensuTester()
// 	assert := assert.New(t)
// 	events, err := sensu.ResolveEvent("server-0-13-0", "check_critical")
// 	assert.Nil(err, fmt.Sprintf("ResolveEvents returned an error: %s", err))
// 	assert.NotNil(events, "ResolveEvents returned nil!")
//
// 	ev, err := sensu.GetEventsCheckForClient("server-0-13-0", "check_critical")
// 	assert.NotNil(err, fmt.Sprintf("Sensu Resolve Events should not return an event after it's deletion. Got : %v", ev))
// }

// func TestResolveNonExistingEvents(t *testing.T) {
// 	sensu := getSensuTester()
// 	ev, err := sensu.ResolveEvent("server-0-13-0", "check_not_real")
// 	assert.NotNil(t, err, fmt.Sprintf("Sensu Resolve Events should not returned and error on non-existing checks: %v", ev))
// }
