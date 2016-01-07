package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLevelInt(t *testing.T) {
	integer := getLevelInt("Trace")
	assert.Equal(t, TRACE, integer)

	integer = getLevelInt("DEBUG")
	assert.Equal(t, DEBUG, integer)

	integer = getLevelInt("info")
	assert.Equal(t, INFO, integer)

	integer = getLevelInt("warn")
	assert.Equal(t, WARN, integer)

	integer = getLevelInt("fatal")
	assert.Equal(t, FATAL, integer)

	integer = getLevelInt("none")
	assert.Equal(t, INFO, integer)
}

func TestIsDisabledFor(t *testing.T) {
	configuredLevel = TRACE

	enabled := isDisabledFor("trace")
	assert.Equal(t, false, enabled)
	enabled = isDisabledFor("debug")
	assert.Equal(t, false, enabled)
	enabled = isDisabledFor("info")
	assert.Equal(t, false, enabled)
	enabled = isDisabledFor("warn")
	assert.Equal(t, false, enabled)
	enabled = isDisabledFor("fatal")
	assert.Equal(t, false, enabled)

	configuredLevel = DEBUG

	enabled = isDisabledFor("trace")
	assert.Equal(t, true, enabled)
	enabled = isDisabledFor("debug")
	assert.Equal(t, false, enabled)
	enabled = isDisabledFor("info")
	assert.Equal(t, false, enabled)
	enabled = isDisabledFor("warn")
	assert.Equal(t, false, enabled)
	enabled = isDisabledFor("fatal")
	assert.Equal(t, false, enabled)

	configuredLevel = INFO

	enabled = isDisabledFor("trace")
	assert.Equal(t, true, enabled)
	enabled = isDisabledFor("debug")
	assert.Equal(t, true, enabled)
	enabled = isDisabledFor("info")
	assert.Equal(t, false, enabled)
	enabled = isDisabledFor("warn")
	assert.Equal(t, false, enabled)
	enabled = isDisabledFor("fatal")
	assert.Equal(t, false, enabled)

	configuredLevel = WARN

	enabled = isDisabledFor("trace")
	assert.Equal(t, true, enabled)
	enabled = isDisabledFor("debug")
	assert.Equal(t, true, enabled)
	enabled = isDisabledFor("info")
	assert.Equal(t, true, enabled)
	enabled = isDisabledFor("warn")
	assert.Equal(t, false, enabled)
	enabled = isDisabledFor("fatal")
	assert.Equal(t, false, enabled)

	configuredLevel = FATAL

	enabled = isDisabledFor("trace")
	assert.Equal(t, true, enabled)
	enabled = isDisabledFor("debug")
	assert.Equal(t, true, enabled)
	enabled = isDisabledFor("info")
	assert.Equal(t, true, enabled)
	enabled = isDisabledFor("warn")
	assert.Equal(t, true, enabled)
	enabled = isDisabledFor("fatal")
	assert.Equal(t, false, enabled)
}
