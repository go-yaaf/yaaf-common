// Copyright 2022. Motty Cohen
//
// Base configuration utility tests

package test

import (
	"os"
	"testing"

	"github.com/go-yaaf/yaaf-common/config"
	"github.com/stretchr/testify/assert"
)

func TestBaseConfig_ReadConfig(t *testing.T) {

	if err := os.Setenv("LOG_LEVEL", "ERROR"); err != nil {
		assert.FailNow(t, err.Error())
	}
	if err := os.Setenv("KEY_2", "12"); err != nil {
		assert.FailNow(t, err.Error())
	}

	config.Get().AddConfigVar("LOG_LEVEL", "ERROR")
	config.Get().AddConfigVar("KEY_2", "true")
	config.Get().AddConfigVar("KEY_3", "456")
	mp := config.Get().GetAllVars()

	if val, ok := mp["LOG_LEVEL"]; !ok {
		assert.False(t, ok, "key not found")
	} else {
		assert.Equal(t, val, "ERROR")
	}

	assert.Equal(t, "ERROR", config.Get().GetStringParamValueOrDefault("LOG_LEVEL", ""))
	assert.Equal(t, true, config.Get().GetBoolParamValueOrDefault("KEY_2", false))
	assert.Equal(t, int64(456), config.Get().GetInt64ParamValueOrDefault("KEY_3", 100))
}
