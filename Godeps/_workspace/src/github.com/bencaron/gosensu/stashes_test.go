package sensu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testPath string = "pouellexxx"

func TestCreateStashes(t *testing.T) {
	assert := assert.New(t)
	sensu := getSensuTester()
	if assert.NotNil(t, sensu) {
		content := make(map[string]interface{})
		content["tesxxxt"] = "allsssso"
    //stash := Stash{testPath, content, 30}
    stash := `{"path": "test", "content":{}, "expire": 10}`
		stashes, err := sensu.CreateStash(stash)
		fmt.Printf("\nstashes: %v err %s\n", stashes, err)
		assert.NotNil(stashes, fmt.Sprintf("CreateStash is nil, error is : %s", err))
	}
}

func TestCreateStashePath(t *testing.T) {
  assert := assert.New(t)
  sensu := getSensuTester()
  if assert.NotNil(t, sensu) {
    content := make(map[string]interface{})
    content["test"] = "allo"
    stashes, err := sensu.CreateStashPath("testPathTest", content)
    fmt.Printf("\nstashes: %v err %s\n", stashes, err)
    assert.NotNil(stashes, fmt.Sprintf("CreateStash is nil, error is : %s", err))
  }
}

func TestGetStashes(t *testing.T) {
	assert := assert.New(t)
	sensu := getSensuTester()
//	var empty []interface{}
	if assert.NotNil(t, sensu) {
		stashes, err := sensu.GetStashes()
		assert.NotNil(stashes, fmt.Sprintf("GetStashes is nil, error is : %s", err))
    // FIXME
		//assert.Equal(stashes, empty)
	}
}
