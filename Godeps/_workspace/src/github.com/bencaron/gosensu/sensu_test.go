package sensu

import (
	"fmt"
	"os"
	"testing"
)

var showedConfig bool = false

func getSensuTester() *Sensu {
	if !showedConfig {
		fmt.Printf("Reading sensu server config from ENV SENSU_SERVER_URL: %s\n", os.Getenv("SENSU_SERVER_URL"))
		showedConfig = true
	}
	sensu := New("Sensu test API", "", os.Getenv("SENSU_SERVER_URL"), 15, "admin", "secret")
	return sensu
}

func TestSensuTester(t *testing.T) {
	sensu := getSensuTester()
	//sensu.NewSensu()
	if sensu == nil {
		t.Error("Sensu object is nil")
	}
}

func TestSensuInfo(t *testing.T) {
	sensu := getSensuTester()
	if sensu == nil {
		t.Error("Sensu object is nil")
	} /* meh. Need to figure out correct type assertions :/
	info, err := sensu.Info()
	switch ty := info.(type) {
	default:
		t.Error("Sensu object is nil")
	case map[string]interface{}:
		fmt.Printf("sensu.GetInfo() is the correct type")
	}
	*/
}
