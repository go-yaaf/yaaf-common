// Copyright 2022. Motty Cohen
//
// Utility functions for metrics monitoring (using Prometheus)
//
package metrics

import (
	"fmt"
	"github.com/agentvi/innovi-core-commons/config"
	"github.com/mottyc/yaaf-common/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const (
	namespace             = "yaaf"
	gaugeConnectedClients = "connectedClients"
)

var (
	gauges          = map[string]prometheus.Gauge{}
	messageCounters = map[uint32]*prometheus.CounterVec{}
)

func AddBytesCounterForTopic(topic string) {
	hash := utils.HashUtils().Hash(topic)
	if _, ok := messageCounters[hash]; !ok {
		messageCounters[hash] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: config.GetBaseConfig().InstrumentingSubsystem(),
				Name:      fmt.Sprintf("topic_%s", topic),
			}, []string{"counter_of"})
		prometheus.MustRegister(messageCounters[hash])
	}
}

func UpdateBytesCounterForTopic(topic string, bytesCount int) {
	hash := utils.HashUtils().Hash(topic)
	if c, ok := messageCounters[hash]; ok {
		c.WithLabelValues("bytes_written_total").Add(float64(bytesCount))
	}
}

func AddMessageCounterForOpCode(opCode int) {
	if _, ok := messageCounters[uint32(opCode)]; !ok {
		messageCounters[uint32(opCode)] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: config.GetBaseConfig().InstrumentingSubsystem(),
				Name:      fmt.Sprintf("op_code_%d", opCode),
			}, []string{"counter_of"})

		prometheus.MustRegister(messageCounters[uint32(opCode)])
	}
}

func UpdateMessageCounterForOpCode(opCode, bytesCount int) {
	if c, ok := messageCounters[uint32(opCode)]; ok {
		c.WithLabelValues("bytes_received_total").Add(float64(bytesCount))
		c.WithLabelValues("messages_received_total").Inc()
	}
}

func AddConnectedClientsGauge() {
	if _, ok := gauges[gaugeConnectedClients]; !ok {
		gauges[gaugeConnectedClients] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: config.GetBaseConfig().InstrumentingSubsystem(),
				Name:      "connected_clients"})
		prometheus.MustRegister(gauges[gaugeConnectedClients])
	}
}

func AddGauge(name, key string) {
	if _, ok := gauges[key]; !ok {
		gauges[key] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: config.GetBaseConfig().InstrumentingSubsystem(),
				Name:      name})
		prometheus.MustRegister(gauges[key])
	}
}
func UpdateGauge(key string, val float64) {
	if v, ok := gauges[key]; ok {
		v.Set(val)
	}
}

func UpdateConnectedClientsGauge(isConnected bool) {
	if isConnected {
		gauges[gaugeConnectedClients].Inc()
	} else {
		gauges[gaugeConnectedClients].Dec()
	}
}

func MetricsScrapingHandler() http.Handler {
	return promhttp.Handler()
}
