// Utility functions for metrics monitoring (using Prometheus)

//

package metrics

/**
import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/utils"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace             = "yaaf"
	subsystem             = "subsystem"
	gaugeConnectedClients = "connectedClients"
)

var (
	gauges          = map[string]prometheus.Gauge{}
	messageCounters = map[uint32]*prometheus.CounterVec{}
)

// AddGauge adds metric gauge
func AddGauge(name, key string) {
	if _, ok := gauges[key]; !ok {
		gauges[key] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      name})
		prometheus.MustRegister(gauges[key])
	}
}

// UpdateGauge updates metric gauge
func UpdateGauge(key string, val float64) {
	if v, ok := gauges[key]; ok {
		v.Set(val)
	}
}

// AddBytesCounterForTopic adds number of bytes for topic KPI
func AddBytesCounterForTopic(topic string) {
	hash := utils.HashUtils().Hash(topic)
	if _, ok := messageCounters[hash]; !ok {
		messageCounters[hash] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      fmt.Sprintf("topic_%s", topic),
			}, []string{"counter_of"})
		prometheus.MustRegister(messageCounters[hash])
	}
}

// UpdateBytesCounterForTopic updates number of bytes for topic KPI
func UpdateBytesCounterForTopic(topic string, bytesCount int) {
	hash := utils.HashUtils().Hash(topic)
	if c, ok := messageCounters[hash]; ok {
		c.WithLabelValues("bytes_written_total").Add(float64(bytesCount))
	}
}

// AddMessageCounterForOpCode adds number of messages for a specific op-code KPI
func AddMessageCounterForOpCode(opCode int) {
	if _, ok := messageCounters[uint32(opCode)]; !ok {
		messageCounters[uint32(opCode)] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      fmt.Sprintf("op_code_%d", opCode),
			}, []string{"counter_of"})

		prometheus.MustRegister(messageCounters[uint32(opCode)])
	}
}

// UpdateMessageCounterForOpCode updates number of messages for a specific op-code KPI
func UpdateMessageCounterForOpCode(opCode, bytesCount int) {
	if c, ok := messageCounters[uint32(opCode)]; ok {
		c.WithLabelValues("bytes_received_total").Add(float64(bytesCount))
		c.WithLabelValues("messages_received_total").Inc()
	}
}

// AddConnectedClientsGauge adds number of connected web-socket clients KPI
func AddConnectedClientsGauge() {
	if _, ok := gauges[gaugeConnectedClients]; !ok {
		gauges[gaugeConnectedClients] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "connected_clients"})
		prometheus.MustRegister(gauges[gaugeConnectedClients])
	}
}

// UpdateConnectedClientsGauge updates number of connected web-socket clients KPI
func UpdateConnectedClientsGauge(isConnected bool) {
	if isConnected {
		gauges[gaugeConnectedClients].Inc()
	} else {
		gauges[gaugeConnectedClients].Dec()
	}
}
**/
