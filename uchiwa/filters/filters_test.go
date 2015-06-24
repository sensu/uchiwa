package filters

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestGetRequest(t *testing.T) {

	unauthorized := GetRequest("", nil)
	assert.Equal(t, false, unauthorized)
}

func TestPostRequest(t *testing.T) {

	var data interface{} = map[string]interface{}{"dc": "foo"}

	unauthorized := PostRequest(nil, &data)
	assert.Equal(t, false, unauthorized)
}

func TestSensuData(t *testing.T) {

	originalData := &structs.Data{
		Aggregates:    []interface{}{map[string]string{"dc": "foo"}, map[string]string{"dc": "bar"}},
		Checks:        []interface{}{map[string]interface{}{"dc": "foo", "subscribers": []string{"linux"}}, map[string]interface{}{"dc": "foo", "subscribers": []string{"mac"}}, map[string]interface{}{"dc": "bar", "subscribers": []string{"windows"}}},
		Clients:       []interface{}{map[string]interface{}{"dc": "foo", "subscriptions": []string{"linux", "mac"}}, map[string]string{"dc": "bar"}},
		Dc:            []*structs.Datacenter{&structs.Datacenter{Name: "foo"}, &structs.Datacenter{Name: "bar"}},
		Events:        []interface{}{map[string]interface{}{"dc": "foo", "check": map[string]interface{}{"subscribers": []string{"mac"}}}, map[string]string{"dc": "bar"}},
		Stashes:       []interface{}{map[string]string{"dc": "foo"}, map[string]string{"dc": "bar"}},
		Subscriptions: []string{"linux", "mac", "windows"},
	}

	data := SensuData(nil, originalData)
	assert.Equal(t, originalData, data)
}
