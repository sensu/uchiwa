package uchiwa

import (
	"fmt"

	"github.com/palourde/logger"
)

func getRedisStatus(info map[string]interface{}) string {
	redis, ok := info["redis"].(map[string]interface{})
	if !ok {
		logger.Warningf("Could not assert info's redis interface %+v", info["redis"])
		return "unknown"
	}
	return fmt.Sprintf("%t", redis["connected"])
}

func getSensuVersion(info map[string]interface{}) string {
	sensu, ok := info["sensu"].(map[string]interface{})
	if !ok {
		logger.Warningf("Could not assert info's sensu interface %+v", info["sensu"])
		return "?"
	}
	return sensu["version"].(string)
}

func getTransportStatus(info map[string]interface{}) string {
	transport, ok := info["transport"].(map[string]interface{})
	if !ok {
		transport, ok := info["rabbitmq"].(map[string]interface{})
		if !ok {
			logger.Warning("Could not assert info's transport interface")
			return "unknown"
		}
		return fmt.Sprintf("%t", transport["connected"])
	}
	return fmt.Sprintf("%t", transport["connected"])
}

// Status build each DC status
func Status(info map[string]interface{}, name string) map[string]string {
	s := make(map[string]string)

	s["name"] = name
	s["version"] = getSensuVersion(info)
	s["transport"] = getTransportStatus(info)
	s["redis"] = getRedisStatus(info)

	return s
}
