package daemon

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestBuildSubscriptions(t *testing.T) {
	// Test basic subscriptions
	data := structs.Data{
		Clients: []interface{}{
			map[string]interface{}{"subscriptions": []string{"foo", "bar"}},
			map[string]interface{}{"subscriptions": []string{"foo", "qux"}},
		},
	}
	d := Daemon{Data: &data}
	d.BuildSubscriptions()
	assert.Equal(t, 3, len(d.Data.Subscriptions))

	// Test per-client subscriptions
	data = structs.Data{
		Clients: []interface{}{
			map[string]interface{}{"subscriptions": []string{"foo", "client:foobar"}},
			map[string]interface{}{"subscriptions": []string{"CLIENT:BAZ", "qux"}},
		},
	}
	d = Daemon{Data: &data}
	d.BuildSubscriptions()
	assert.Equal(t, 2, len(d.Data.Subscriptions))
}

func TestIsSubscriptionInSubscriptions(t *testing.T) {
	subscriptions := []structs.Subscription{
		structs.Subscription{Dc: "us-east-1", Name: "foo"},
		structs.Subscription{Dc: "us-east-1", Name: "bar"},
	}

	subscription := structs.Subscription{Dc: "us-east-2", Name: "foo"}
	result := isSubscriptionInSubscriptions(subscription, subscriptions)
	assert.False(t, result)

	subscription = structs.Subscription{Dc: "us-east-1", Name: "foo"}
	result = isSubscriptionInSubscriptions(subscription, subscriptions)
	assert.True(t, result)
}
