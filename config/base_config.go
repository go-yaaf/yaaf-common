// Package config
//
// Base configuration utility
// This utility is used to read configuration parameters from environment variables and expose them to the application
// parts as accessor methods.
//
// The concrete application / service should inherit its own special configuration and extend the base configuration
// with added variables. The base configuration exposes some common configuration parameters used by the middleware
// components

package config

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	CfgLoglevel               = "LOG_LEVEL"
	CfgHttpReadTimeoutMs      = "HTTP_READ_TIMEOUT_MS"
	CfgHttpWriteTimeoutMs     = "HTTP_WRITE_TIMEOUT_MS"
	CfgWsKeepAliveSec         = "WS_KEEP_ALIVE_SEC"
	CfgWsReadBufferSizeBytes  = "WS_READ_BUFFER_SIZE_BYTES"
	CfgWsWriteBufferSizeBytes = "WS_WRITE_BUFFER_SIZE_BYTES"
	CfgWsWriteCompress        = "WS_WRITE_COMPRESS"
	CfgWsWriteTimeoutSec      = "WS_WRITE_TIMEOUT"
	CfgWsPongTimeoutSec       = "WS_PONG_TIMEOUT"
	CfgTopicPartitions        = "TOPIC_PARTITIONS"
)

// region BaseConfig singleton pattern ---------------------------------------------------------------------------------

var initOnce sync.Once
var baseCfg *BaseConfig

type BaseConfig struct {
	cfg map[string]string
}

// Create new
func newBaseConfig() *BaseConfig {
	var bc = BaseConfig{}
	bc.cfg = map[string]string{
		CfgLoglevel:               "INFO",
		CfgHttpReadTimeoutMs:      "3000",
		CfgHttpWriteTimeoutMs:     "3000",
		CfgWsKeepAliveSec:         "-1",
		CfgWsReadBufferSizeBytes:  "1048576",
		CfgWsWriteBufferSizeBytes: "1048576",
		CfgWsWriteCompress:        "true",
		CfgWsPongTimeoutSec:       "5",
		CfgWsWriteTimeoutSec:      "5",
	}
	return &bc
}

// Get singleton instance
func Get() *BaseConfig {
	initOnce.Do(func() {
		baseCfg = newBaseConfig()
		baseCfg.ScanEnvVariables()
	})
	return baseCfg
}

// endregion

// region Helper methods -----------------------------------------------------------------------------------------------

// GetAllVars gets a map of all the configuration variables and values
func (c *BaseConfig) GetAllVars() map[string]string {
	result := make(map[string]string)
	for key, value := range c.cfg {
		result[key] = value
	}
	return result
}

// GetAllKeysSorted gets a list of all the configuration keys
func (c *BaseConfig) GetAllKeysSorted() []string {

	keys := make([]string, 0, len(c.cfg))
	for k := range c.cfg {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// AddConfigVar adds or updates configuration variable
func (c *BaseConfig) AddConfigVar(key, value string) {
	c.cfg[key] = value
}

// ScanEnvVariables scans all environment variables and map their values to existing configuration keys
func (c *BaseConfig) ScanEnvVariables() {
	for key := range c.cfg {
		if tmp := os.Getenv(key); tmp != "" {
			c.cfg[key] = tmp
		}
	}
}

// GetIntParamValueOrDefault gets environment variable as int
func (c *BaseConfig) GetIntParamValueOrDefault(key string, defaultValue int) (val int) {
	val = defaultValue
	if len(c.cfg[key]) > 0 {
		if v, err := strconv.Atoi(c.cfg[key]); err == nil {
			val = v
		}
	}
	return
}

// GetStringParamValueOrDefault gets environment variable as string
func (c *BaseConfig) GetStringParamValueOrDefault(key string, defaultValue string) (val string) {
	val = defaultValue
	if len(c.cfg[key]) > 0 {
		val = c.cfg[key]
	}
	return
}

// GetInt64ParamValueOrDefault gets environment variable as int64
func (c *BaseConfig) GetInt64ParamValueOrDefault(key string, defaultValue int64) (val int64) {
	val = defaultValue
	if len(c.cfg[key]) > 0 {
		val, _ = strconv.ParseInt(c.cfg[key], 10, 64)
	}
	return
}

// GetBoolParamValueOrDefault gets environment variable as bool
func (c *BaseConfig) GetBoolParamValueOrDefault(key string, defaultValue bool) (val bool) {
	val = defaultValue
	if len(c.cfg[key]) > 0 {
		tmp := strings.ToLower(c.cfg[key])
		val = tmp == "true" || tmp == "1"
	}
	return
}

// endregion

// region Configuration accessors methods ------------------------------------------------------------------------------

// LogLevel gets log level
func (c *BaseConfig) LogLevel() string {
	return c.GetStringParamValueOrDefault(CfgLoglevel, "INFO")
}

// HttpReadTimeoutMs gets HTTP read time out in milliseconds
func (c *BaseConfig) HttpReadTimeoutMs() int {
	return c.GetIntParamValueOrDefault(CfgHttpReadTimeoutMs, 3000)
}

// HttpWriteTimeoutMs gets HTTP write time out in milliseconds
func (c *BaseConfig) HttpWriteTimeoutMs() int {
	return c.GetIntParamValueOrDefault(CfgHttpWriteTimeoutMs, 3000)
}

// WsKeepALiveInterval gets web socket keep alive interval (in seconds)
func (c *BaseConfig) WsKeepALiveInterval() int64 {
	return c.GetInt64ParamValueOrDefault(CfgWsKeepAliveSec, -1)
}

// WsReadBufferSizeBytes gets web socket read buffer size
func (c *BaseConfig) WsReadBufferSizeBytes() int {
	return c.GetIntParamValueOrDefault(CfgWsReadBufferSizeBytes, 1048576)
}

// WsWriteBufferSizeBytes gets web socket write buffer size
func (c *BaseConfig) WsWriteBufferSizeBytes() int {
	return c.GetIntParamValueOrDefault(CfgWsWriteBufferSizeBytes, 1048576)
}

// WsWriteCompress gets web socket compression on write flag
func (c *BaseConfig) WsWriteCompress() bool {
	return c.GetBoolParamValueOrDefault(CfgWsWriteCompress, true)
}

// WsPongTimeoutSec gets web socket PONG time out in seconds
func (c *BaseConfig) WsPongTimeoutSec() int {
	return c.GetIntParamValueOrDefault(CfgWsPongTimeoutSec, 5)
}

// WsWriteTimeoutSec gets web socket write time out in seconds
func (c *BaseConfig) WsWriteTimeoutSec() int {
	return c.GetIntParamValueOrDefault(CfgWsWriteTimeoutSec, 5)
}

// TopicPartitions gets default number of partitions per topic
func (c *BaseConfig) TopicPartitions() int {
	return c.GetIntParamValueOrDefault(CfgTopicPartitions, 1)
}

// endregion
