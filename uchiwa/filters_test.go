package uchiwa

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/sensu/uchiwa/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestFilterGetRequest(t *testing.T) {

	unauthorized := filterGetRequest("", nil)
	assert.Equal(t, true, unauthorized, "a request with an empty datacenter should not be authorized")

	unauthorized = filterGetRequest("foo", nil)
	assert.Equal(t, false, unauthorized, "a request with a nil token should be authorized")

	// mock a JWT
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["Role"] = auth.Role{}
	unauthorized = filterGetRequest("foo", token)
	assert.Equal(t, false, unauthorized, "a request with a token that contains an empty Role attribute should be authorized")

	// set a datacenters filter
	token.Claims["Role"] = auth.Role{
		Datacenters: []string{"foo", "bar"},
	}

	// request an unauthorized datacenter
	unauthorized = filterGetRequest("qux", token)
	assert.Equal(t, true, unauthorized, "a request with an unauthorized datacenter should not be authorized")

	// request an authorized datacenter
	unauthorized = filterGetRequest("bar", token)
	assert.Equal(t, false, unauthorized, "a request with an authorized datacenter should be authorized")
}

func TestFilterPostRequest(t *testing.T) {

	var data interface{} = map[string]interface{}{"dc": "foo"}

	unauthorized := filterPostRequest(nil, &data)
	assert.Equal(t, false, unauthorized, "a request with a nil token should be authorized")

	// mock a JWT
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["Role"] = auth.Role{}

	unauthorized = filterPostRequest(token, &data)
	assert.Equal(t, false, unauthorized, "both empty datacenters and subscriptions filters within Role should be authorized")

	// set an unknown datacenter
	token.Claims["Role"] = auth.Role{
		Datacenters: []string{"qux"},
	}

	unauthorized = filterPostRequest(token, &data)
	assert.Equal(t, true, unauthorized, "only an element with the datacenter 'qux' should be authorized")

	// set a datacenter
	token.Claims["Role"] = auth.Role{
		Datacenters: []string{"foo"},
	}

	unauthorized = filterPostRequest(token, &data)
	assert.Equal(t, false, unauthorized, "an element with the datacenter 'foo' should be authorized")

	// set an unknown subscriptions
	data = interface{}(map[string]interface{}{"dc": "foo", "check": map[string]interface{}{"subscribers": []string{"linux"}}})

	token.Claims["Role"] = auth.Role{
		Subscriptions: []string{"windows"},
	}

	unauthorized = filterPostRequest(token, &data)
	assert.Equal(t, true, unauthorized, "only an element with the subscription 'windows' should be authorized")

	// set a subscription
	token.Claims["Role"] = auth.Role{
		Subscriptions: []string{"linux"},
	}

	unauthorized = filterPostRequest(token, &data)
	assert.Equal(t, false, unauthorized, "an element with the subscription 'linux' should be authorized")

	// set an authorized datacenter but an unauthorized subscription
	token.Claims["Role"] = auth.Role{
		Datacenters:   []string{"foo"},
		Subscriptions: []string{"windows"},
	}

	unauthorized = filterPostRequest(token, &data)
	assert.Equal(t, true, unauthorized, "only an element with the 'foo' datacenter and the subscription 'linux' should be authorized")

	// set both datacenter and subscription
	token.Claims["Role"] = auth.Role{
		Datacenters:   []string{"foo"},
		Subscriptions: []string{"linux"},
	}

	unauthorized = filterPostRequest(token, &data)
	assert.Equal(t, false, unauthorized, "an element with the 'foo' datacenter and the subscription 'linux' should be authorized")

}

func TestFilterSensu(t *testing.T) {

	originalData := &structs.Data{
		Aggregates:    []interface{}{map[string]string{"dc": "foo"}, map[string]string{"dc": "bar"}},
		Checks:        []interface{}{map[string]interface{}{"dc": "foo", "subscribers": []string{"linux"}}, map[string]interface{}{"dc": "foo", "subscribers": []string{"mac"}}, map[string]interface{}{"dc": "bar", "subscribers": []string{"windows"}}},
		Clients:       []interface{}{map[string]interface{}{"dc": "foo", "subscriptions": []string{"linux", "mac"}}, map[string]string{"dc": "bar"}},
		Dc:            []*structs.Datacenter{&structs.Datacenter{Name: "foo"}, &structs.Datacenter{Name: "bar"}},
		Events:        []interface{}{map[string]interface{}{"dc": "foo", "check": map[string]interface{}{"subscribers": []string{"mac"}}}, map[string]string{"dc": "bar"}},
		Stashes:       []interface{}{map[string]string{"dc": "foo"}, map[string]string{"dc": "bar"}},
		Subscriptions: []string{"linux", "mac", "windows"},
	}

	data := filterSensu(nil, originalData)
	assert.Equal(t, originalData, data, "a nil token should return the original data")

	data = filterSensu(jwt.New(jwt.SigningMethodHS256), originalData)
	assert.Equal(t, &structs.Data{}, data, "an invalid token should return an empty Data struct")

	// mock a JWT
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["Role"] = auth.Role{}

	data = filterSensu(token, originalData)
	assert.Equal(t, originalData, data, "both empty datacenters and subscriptions filters should return the original data")

	// set an unknown datacenter
	token.Claims["Role"] = auth.Role{
		Datacenters: []string{"qux"},
	}

	expectedData := &structs.Data{
		Subscriptions: []string{"linux", "mac", "windows"},
	}

	data = filterSensu(token, originalData)
	assert.Equal(t, expectedData, data, "the unknown datacenter 'qux' was not properly filtered")

	// set all datacenters
	token.Claims["Role"] = auth.Role{
		Datacenters: []string{"foo", "bar"},
	}

	data = filterSensu(token, originalData)
	assert.Equal(t, originalData, data, "both datacenters 'foo' & 'bar' were not properly filtered")

	// set datacenters
	token.Claims["Role"] = auth.Role{
		Datacenters: []string{"foo"},
	}

	expectedData = &structs.Data{
		Aggregates:    []interface{}{map[string]string{"dc": "foo"}},
		Checks:        []interface{}{map[string]interface{}{"dc": "foo", "subscribers": []string{"linux"}}, map[string]interface{}{"dc": "foo", "subscribers": []string{"mac"}}},
		Clients:       []interface{}{map[string]interface{}{"dc": "foo", "subscriptions": []string{"linux", "mac"}}},
		Dc:            []*structs.Datacenter{&structs.Datacenter{Name: "foo"}},
		Events:        []interface{}{map[string]interface{}{"dc": "foo", "check": map[string]interface{}{"subscribers": []string{"mac"}}}},
		Stashes:       []interface{}{map[string]string{"dc": "foo"}},
		Subscriptions: []string{"linux", "mac", "windows"},
	}

	data = filterSensu(token, originalData)
	assert.Equal(t, expectedData, data, "the datacenter 'foo' was not properly filtered")

	// set subscriptions
	token.Claims["Role"] = auth.Role{
		Subscriptions: []string{"mac"},
	}

	expectedData = &structs.Data{
		Aggregates:    []interface{}{map[string]string{"dc": "foo"}, map[string]string{"dc": "bar"}},
		Checks:        []interface{}{map[string]interface{}{"dc": "foo", "subscribers": []string{"mac"}}},
		Clients:       []interface{}{map[string]interface{}{"dc": "foo", "subscriptions": []string{"linux", "mac"}}},
		Dc:            []*structs.Datacenter{&structs.Datacenter{Name: "foo"}, &structs.Datacenter{Name: "bar"}},
		Events:        []interface{}{map[string]interface{}{"dc": "foo", "check": map[string]interface{}{"subscribers": []string{"mac"}}}},
		Stashes:       []interface{}{map[string]string{"dc": "foo"}, map[string]string{"dc": "bar"}},
		Subscriptions: []string{"mac"},
	}
	data = filterSensu(token, originalData)
	assert.Equal(t, expectedData, data, "the subscription 'mac' was not properly filtered")

	// set both subscriptions and datacenters
	token.Claims["Role"] = auth.Role{
		Datacenters:   []string{"foo"},
		Subscriptions: []string{"linux"},
	}

	expectedData = &structs.Data{
		Aggregates:    []interface{}{map[string]string{"dc": "foo"}},
		Checks:        []interface{}{map[string]interface{}{"dc": "foo", "subscribers": []string{"linux"}}},
		Clients:       []interface{}{map[string]interface{}{"dc": "foo", "subscriptions": []string{"linux", "mac"}}},
		Dc:            []*structs.Datacenter{&structs.Datacenter{Name: "foo"}},
		Events:        []interface{}(nil),
		Stashes:       []interface{}{map[string]string{"dc": "foo"}},
		Subscriptions: []string{"linux"},
	}
	data = filterSensu(token, originalData)
	assert.Equal(t, expectedData, data, "the datacenter 'foo' & the subscription 'linux' were not properly filtered")
}
