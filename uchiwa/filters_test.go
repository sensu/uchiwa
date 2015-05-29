package uchiwa

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/sensu/uchiwa/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestFilterGetSensu(t *testing.T) {

	originalData := &structs.Data{
		Aggregates: []interface{}{map[string]string{"dc": "foo"}, map[string]string{"dc": "bar"}},
		Checks:     []interface{}{map[string]interface{}{"dc": "foo", "subscribers": []string{"linux"}}, map[string]interface{}{"dc": "foo", "subscribers": []string{"mac"}}, map[string]interface{}{"dc": "bar", "subscribers": []string{"windows"}}},
		Clients:    []interface{}{map[string]interface{}{"dc": "foo", "subscriptions": []string{"linux", "mac"}}, map[string]string{"dc": "bar"}},
		Dc:         []*structs.Datacenter{&structs.Datacenter{Name: "foo"}, &structs.Datacenter{Name: "bar"}},
		Events:     []interface{}{map[string]interface{}{"dc": "foo", "check": map[string]interface{}{"subscribers": []string{"mac"}}}, map[string]string{"dc": "bar"}},
		Stashes:    []interface{}{map[string]string{"dc": "foo"}, map[string]string{"dc": "bar"}},
	}

	data := filterGetSensu(nil, originalData)
	assert.Equal(t, originalData, data, "a nil token should return the original data")

	data = filterGetSensu(jwt.New(jwt.SigningMethodHS256), originalData)
	assert.Equal(t, &structs.Data{}, data, "an invalid token should return an empty Data struct")

	// mock a JWT
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["Role"] = auth.Role{}

	data = filterGetSensu(token, originalData)
	assert.Equal(t, originalData, data, "empty datacenters and subscriptions attributes should return the original data")

	// set an unknown datacenter
	token.Claims["Role"] = auth.Role{
		Datacenters: []string{"qux"},
	}

	data = filterGetSensu(token, originalData)
	assert.Equal(t, &structs.Data{}, data, "the unknown datacenter 'qux' was not properly filtered")

	// set all datacenters
	token.Claims["Role"] = auth.Role{
		Datacenters: []string{"foo", "bar"},
	}

	data = filterGetSensu(token, originalData)
	assert.Equal(t, originalData, data, "both datacenters 'foo' & 'bar' were not properly filtered")

	// set one datacenter within the Token
	token.Claims["Role"] = auth.Role{
		Datacenters: []string{"foo"},
	}

	expectedData := &structs.Data{
		Aggregates: []interface{}{map[string]string{"dc": "foo"}},
		Checks:     []interface{}{map[string]interface{}{"dc": "foo", "subscribers": []string{"linux"}}, map[string]interface{}{"dc": "foo", "subscribers": []string{"mac"}}},
		Clients:    []interface{}{map[string]interface{}{"dc": "foo", "subscriptions": []string{"linux", "mac"}}},
		Dc:         []*structs.Datacenter{&structs.Datacenter{Name: "foo"}},
		Events:     []interface{}{map[string]interface{}{"dc": "foo", "check": map[string]interface{}{"subscribers": []string{"mac"}}}},
		Stashes:    []interface{}{map[string]string{"dc": "foo"}},
	}

	data = filterGetSensu(token, originalData)
	assert.Equal(t, expectedData, data, "the datacenter 'foo' was not properly filtered")

	// set one subscription within the Token
	token.Claims["Role"] = auth.Role{
		Subscriptions: []string{"mac"},
	}

	expectedData = &structs.Data{
		Aggregates: []interface{}{map[string]string{"dc": "foo"}, map[string]string{"dc": "bar"}},
		Checks:     []interface{}{map[string]interface{}{"dc": "foo", "subscribers": []string{"mac"}}},
		Clients:    []interface{}{map[string]interface{}{"dc": "foo", "subscriptions": []string{"linux", "mac"}}},
		Dc:         []*structs.Datacenter{&structs.Datacenter{Name: "foo"}, &structs.Datacenter{Name: "bar"}},
		Events:     []interface{}{map[string]interface{}{"dc": "foo", "check": map[string]interface{}{"subscribers": []string{"mac"}}}},
		Stashes:    []interface{}{map[string]string{"dc": "foo"}, map[string]string{"dc": "bar"}},
	}
	data = filterGetSensu(token, originalData)
	assert.Equal(t, expectedData, data, "the subscription 'mac' was not properly filtered")

	// set both subscription and datacenter within the Token
	token.Claims["Role"] = auth.Role{
		Datacenters:   []string{"foo"},
		Subscriptions: []string{"linux"}, // <----- SHOULD WE ONLY DISPLAY THE LINUX SUB. OF THE BAR DC?
	}

	expectedData = &structs.Data{
		Aggregates: []interface{}{map[string]string{"dc": "foo"}},
		Checks:     []interface{}{map[string]interface{}{"dc": "foo", "subscribers": []string{"linux"}}},
		Clients:    []interface{}{map[string]interface{}{"dc": "foo", "subscriptions": []string{"linux", "mac"}}},
		Dc:         []*structs.Datacenter{&structs.Datacenter{Name: "foo"}},
		Events:     []interface{}(nil),
		Stashes:    []interface{}{map[string]string{"dc": "foo"}},
	}
	data = filterGetSensu(token, originalData)
	assert.Equal(t, expectedData, data, "the datacenter 'foo' & the subscription 'linux' were not properly filtered")

}
