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
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	CfgPulseDeployment         = "PULSE_DEPLOYMENT"           // Pulse deployment name
	CfgHostName                = "HOSTNAME"                   // Host name
	CfgLoglevel                = "LOG_LEVEL"                  // Log level
	CfgLogJsonFormat           = "LOG_JSON_FORMAT"            // Enable Json log format
	CfgHttpReadTimeoutMs       = "HTTP_READ_TIMEOUT_MS"       // HTTP read timeout (in milliseconds)
	CfgHttpWriteTimeoutMs      = "HTTP_WRITE_TIMEOUT_MS"      // HTTP write timeout (in milliseconds)
	CfgWsKeepAliveSec          = "WS_KEEP_ALIVE_SEC"          // Web socket keep-alive period (in seconds)
	CfgWsReadBufferSizeBytes   = "WS_READ_BUFFER_SIZE_BYTES"  // Web socket read buffer size (in bytes)
	CfgWsWriteBufferSizeBytes  = "WS_WRITE_BUFFER_SIZE_BYTES" // Web socket write buffer size (in bytes)
	CfgWsWriteCompress         = "WS_WRITE_COMPRESS"          // Web socket enable compression
	CfgWsWriteTimeoutSec       = "WS_WRITE_TIMEOUT"           // Web socket write timeout (in seconds)
	CfgWsPongTimeoutSec        = "WS_PONG_TIMEOUT"            // Web socket pong response timeout (in seconds)
	CfgTopicPartitions         = "TOPIC_PARTITIONS"           // Number of partitions per topic
	CfgEnableGoRuntimeProfiler = "ENABLE_GO_RUNTIME_PROFILER" // Enable / Disable runtime profiling
	CfgGoRuntimeProfilerAddr   = "GO_RUNTIME_PROFILER_ADDR"   // Go runtime profiler address

	CfgDatabaseUri  = "DATABASE_URI"  // Configuration database URI
	CfgDatastoreUri = "DATASTORE_URI" // Big data store URI
	CfgMessagingUri = "MESSAGING_URI" // Real-time Messaging middleware URI
	CfgDataCacheUri = "DATACACHE_URI" // Distributed cache middleware URI
	CfgStreamingUri = "STREAMING_URI" // Streaming middleware URI

	// CfgEnableMessageOrdering is set to true to ensure that messages with the same ordering key are delivered in the order they were published.
	// This is crucial for use cases where the order of messages is important for correct processing.
	// Note: The Pub/Sub topic must be configured to support message ordering for this to take effect.
	// Enabling message ordering may impact the throughput of message publishing, as it requires Pub/Sub to maintain order within each ordering key.
	CfgEnableMessageOrdering = "ENABLE_MESSAGE_ORDERING"

	// CfgPubSubNumOfGoroutines NumGoroutines specifies the number of goroutines that will be used
	// to pull messages from the subscription in parallel. Each goroutine
	// opens a separate StreamingPull stream. A higher number of goroutines
	// might increase throughput but also increases the system's load.
	// Defaults to DefaultReceiveSettings.NumGoroutines when set to 0.
	CfgPubSubNumOfGoroutines = "PUBSUB_NUM_OF_GOROUTINES"

	// CfgPubSubMaxOutstandingMessages defines the maximum number of unprocessed
	// messages (messages that have been received but not yet acknowledged
	// or expired). Setting this to a lower number can prevent the consumer
	// from being overwhelmed by a large volume of incoming messages.
	// If set to 0, the default is DefaultReceiveSettings.MaxOutstandingMessages.
	// A negative value indicates no limit.
	CfgPubSubMaxOutstandingMessages = "PUBSUB_MAX_OUTSTANDING_MESSAGES"

	// CfgPubSubMaxOutstandingBytes is the maximum total size of unprocessed messages.
	// This setting helps to control memory usage by limiting the total size
	// of messages that can be held in memory at a time. If set to 0, the
	// default is DefaultReceiveSettings.MaxOutstandingBytes. A negative
	// value indicates no limit on the byte size of unprocessed messages.
	CfgPubSubMaxOutstandingBytes = "PUSUB_MAX_OUTSTANDING_BYTES"

	CfgRdsInstanceName = "RDS_INSTANCE_NAME"

	CfgMaxDbConnections = "MAX_DB_CONNECTIONS"
)

const (
	DefaultPubSubNumOfGoroutines = 0
	DefaultPubSubMaxOutstandingMessages
	DefaultPubSubMaxOutstandingBytes = 0
	DefaultEnableMessageOrdering     = false
	DefaultEnableGoRuntimeProfiler   = false
	DefaultGoRuntimeProfilerAddr     = ":6060"
	DefaultMaxDbConnections          = 10
)

// region BaseConfig singleton pattern ---------------------------------------------------------------------------------

var initOnce sync.Once
var baseCfg *BaseConfig

type BaseConfig struct {
	cfg       map[string]string
	startTime entity.Timestamp
}

// Create new
func newBaseConfig() *BaseConfig {
	var bc = BaseConfig{}
	bc.cfg = map[string]string{
		CfgPulseDeployment:              "",
		CfgHostName:                     "",
		CfgLoglevel:                     "INFO",
		CfgLogJsonFormat:                "false",
		CfgHttpReadTimeoutMs:            "3000",
		CfgHttpWriteTimeoutMs:           "3000",
		CfgWsKeepAliveSec:               "-1",
		CfgWsReadBufferSizeBytes:        "1048576",
		CfgWsWriteBufferSizeBytes:       "1048576",
		CfgWsWriteCompress:              "true",
		CfgWsPongTimeoutSec:             "5",
		CfgWsWriteTimeoutSec:            "5",
		CfgPubSubNumOfGoroutines:        fmt.Sprintf("%d", DefaultPubSubNumOfGoroutines),
		CfgPubSubMaxOutstandingMessages: fmt.Sprintf("%d", DefaultPubSubMaxOutstandingMessages),
		CfgPubSubMaxOutstandingBytes:    fmt.Sprintf("%d", DefaultPubSubMaxOutstandingBytes),
		CfgEnableMessageOrdering:        fmt.Sprintf("%t", DefaultEnableMessageOrdering),
		CfgEnableGoRuntimeProfiler:      fmt.Sprintf("%t", DefaultEnableGoRuntimeProfiler),
		CfgGoRuntimeProfilerAddr:        DefaultGoRuntimeProfilerAddr,
		CfgDatabaseUri:                  "",
		CfgDatastoreUri:                 "",
		CfgMessagingUri:                 "",
		CfgDataCacheUri:                 "",
		CfgStreamingUri:                 "",
		CfgRdsInstanceName:              "",
		CfgMaxDbConnections:             fmt.Sprintf("%d", DefaultMaxDbConnections),
	}
	bc.startTime = entity.Now()
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

// StartTime returns the start time of the service
func (c *BaseConfig) StartTime() (result entity.Timestamp) {
	return c.startTime
}

// PulseDeployment fetch pulse deployment name
func (c *BaseConfig) PulseDeployment() string {
	return c.GetStringParamValueOrDefault(CfgPulseDeployment, "")
}

// HostName fetch current host (machine) name
func (c *BaseConfig) HostName() string {
	return c.GetStringParamValueOrDefault(CfgHostName, "")
}

// RdsInstanceName fetch RDS instance
func (c *BaseConfig) RdsInstanceName() string {
	return c.GetStringParamValueOrDefault(CfgRdsInstanceName, "")
}

func (c *BaseConfig) PubSubNumOfGoroutines() int {
	return c.GetIntParamValueOrDefault(CfgPubSubNumOfGoroutines, DefaultPubSubNumOfGoroutines)
}

func (c *BaseConfig) MaxDbConnections() int {
	return c.GetIntParamValueOrDefault(CfgMaxDbConnections, DefaultMaxDbConnections)
}

func (c *BaseConfig) PubSubMaxOutstandingMessages() int {
	return c.GetIntParamValueOrDefault(CfgPubSubMaxOutstandingMessages, DefaultPubSubMaxOutstandingMessages)
}

func (c *BaseConfig) PubSubMaxOutstandingBytes() int {
	return c.GetIntParamValueOrDefault(CfgPubSubMaxOutstandingBytes, DefaultPubSubMaxOutstandingBytes)
}

// LogLevel gets log level
func (c *BaseConfig) LogLevel() string {
	return c.GetStringParamValueOrDefault(CfgLoglevel, "INFO")
}

// EnableLogJsonFormat returns json log format flag
func (c *BaseConfig) EnableLogJsonFormat() bool {
	return c.GetBoolParamValueOrDefault(CfgLogJsonFormat, false)
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

func (c *BaseConfig) EnableMessageOrdering() bool {
	return c.GetBoolParamValueOrDefault(CfgEnableMessageOrdering, DefaultEnableMessageOrdering)
}

func (c *BaseConfig) EnableGoRuntimeProfiler() bool {
	return c.GetBoolParamValueOrDefault(CfgEnableGoRuntimeProfiler, DefaultEnableGoRuntimeProfiler)
}

func (c *BaseConfig) GoRuntimeProfilerAddr() string {
	return c.GetStringParamValueOrDefault(CfgGoRuntimeProfilerAddr, DefaultGoRuntimeProfilerAddr)
}

// DatabaseUri returns the database URI
func (c *BaseConfig) DatabaseUri() string {
	return c.GetStringParamValueOrDefault(CfgDatabaseUri, "")
}

// DatastoreUri returns the big data store URI
func (c *BaseConfig) DatastoreUri() string {
	return c.GetStringParamValueOrDefault(CfgDatastoreUri, "")
}

// DataCacheUri returns the distributed cache middleware URI
func (c *BaseConfig) DataCacheUri() string {
	return c.GetStringParamValueOrDefault(CfgDataCacheUri, "")
}

// MessagingUri returns the real-time messaging middleware URI
func (c *BaseConfig) MessagingUri() string {
	return c.GetStringParamValueOrDefault(CfgMessagingUri, "")
}

// StreamingUri returns the streaming middleware URI
func (c *BaseConfig) StreamingUri() string {
	return c.GetStringParamValueOrDefault(CfgStreamingUri, "")
}

// endregion
