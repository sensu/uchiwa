package sensu

import "fmt"

// Health The health endpoint checks to see if the api can connect to redis and rabbitmq. It takes parameters for minimum consumers and maximum messages and checks rabbitmq.
func (s *Sensu) Health(consumers int, messages int) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("health/%d/%d", consumers, messages))
}
