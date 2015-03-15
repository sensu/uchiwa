package sensu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChecks(t *testing.T) {
	sensu := getSensuTester()
	assert := assert.New(t)

	res, err := sensu.GetChecks()
	assert.Nil(err, fmt.Sprintf("GetChecks returned an error: %s", err))
	assert.NotNil(res, "GetChecks returned nil!")
}

func TestGetOneCheck(t *testing.T) {
	sensu := getSensuTester()
	assert := assert.New(t)

	res, err := sensu.GetCheck("check_success")
	assert.Nil(err, fmt.Sprintf("GetCheck(check_sucess) returned an error: %s", err))
	assert.NotNil(res, "GetCheck returned nil!")
}

func TestRequestCheck(t *testing.T) {
	sensu := getSensuTester()
	assert := assert.New(t)

	t.Skip("Skipping TestRequestCheck, code not ready yet.")

	res, err := sensu.RequestCheck("chef_success")
	assert.Nil(err, fmt.Sprintf("RequestCheck(check_sucess) returned an error: %s", err))
	assert.NotNil(res, "RequestCheck returned nil!")
}
