// Copyright 2022. Motty Cohen
//
// Web socket listener (server) implementation
//
package socket

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/mottyc/yaaf-common/logger"
	"github.com/mottyc/yaaf-common/metrics"
	"github.com/mottyc/yaaf-common/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsListener struct {
	registry          IWSClientRegistry
	decoder           IMessageDecoder
	handlers          map[int]IWSMessageHandler
	enablePinPong     bool
	enableMetrics     bool
	keepAliveInterval int
}

func NewListener(cfg WsEndpointConfig) (wsh *wsListener) {
	wsh = &wsListener{
		registry: func() IWSClientRegistry {
			if cfg.WsRegistry == nil {
				return &DefaultClientRegistry{Connections: map[string]IWSClient{}}
			} else {
				return cfg.WsRegistry
			}
		}(),
		handlers:          make(map[int]IWSMessageHandler, len(cfg.WsHandlers)),
		enablePinPong:     cfg.EnablePingPong,
		enableMetrics:     cfg.EnableMetrics,
		keepAliveInterval: cfg.KeepAliveInterval,
	}

	// Set message decoder
	if cfg.CustomDecoder != nil {
		wsh.decoder = cfg.CustomDecoder
	} else {
		wsh.decoder = NewJsonDecoder()
	}

	for _, handlerEntry := range cfg.WsHandlers {
		wsh.handlers[handlerEntry.OpCode] = handlerEntry.Handler
		if handlerEntry.Metrics && cfg.EnableMetrics {
			metrics.AddMessageCounterForOpCode(handlerEntry.OpCode)
		}
	}

	if cfg.EnableMetrics {
		metrics.AddConnectedClientsGauge()
	}
	return
}

func (h *wsListener) ListenForWSConnections(w http.ResponseWriter, r *http.Request) {

	connectedClients := h.registry.ConnectedClients()
	maxConnectedClients := 10000

	if maxConnectedClients > 0 && connectedClients == maxConnectedClients {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("error upgrading connection from %s to Web Socket: %s", r.RemoteAddr, err.Error())
		return
	}

	// Fetch value of request's ctx for "client_id" key/value pair.
	// if found, will be used as client connection id. otherwise, generated id will be used.
	clientId := uuid.New().String()
	qParams := make(map[string]string)

	// Get client id from context
	if r.Context().Value("clientId") != nil {
		clientId = r.Context().Value("clientId").(string)
	}

	// Get extra query params from context
	if r.Context().Value("params") != nil {
		qParams = r.Context().Value("params").(map[string]string)
	}

	// Get query params
	for k, v := range r.URL.Query() {
		qParams[k] = v[0]
	}

	// Inject HTTP headers to the params
	for k, v := range r.Header {
		qParams[k] = fmt.Sprintf("%v", v)
	}

	conn.EnableWriteCompression(true)

	tcpConn := conn.UnderlyingConn().(*net.TCPConn)
	_ = tcpConn.SetLinger(0)
	_ = tcpConn.SetNoDelay(true)
	_ = tcpConn.SetWriteBuffer(1048576)
	_ = tcpConn.SetReadBuffer(1048576)

	if h.keepAliveInterval != -1 {
		_ = tcpConn.SetKeepAlive(true)
		_ = tcpConn.SetKeepAlivePeriod(time.Second * time.Duration(h.keepAliveInterval))
	}

	wsCgf := WsClientConfig{
		Id:             clientId,
		WsConn:         conn,
		OnMsgRvd:       h.onMessageReceived,
		OnDisconnected: h.resolveOnDisconnectedCb(),
		MessageDecoder: h.decoder,
		QueryParams:    qParams,
		PinPongEnabled: h.enablePinPong,
	}
	wsClient := NewWsClient2(wsCgf)
	h.registry.RegisterClient(wsClient)
	if h.enableMetrics {
		metrics.UpdateConnectedClientsGauge(true)
	}
	return
}

func (h *wsListener) onMessageReceived(ws IWSClient, m IWSMessage, bytes int) {

	var (
		ok      bool
		handler IWSMessageHandler
	)

	defer utils.RecoverAll(func(err interface{}) {
		logger.Error("wsListener::onMessageReceived error: %s", err)
	})

	//Debug("%s onMessageReceived invoked ", ws.Id())

	client := ws.(*wsClient)

	if client.isClosed() {
		return
	}

	remoteAddress := client.RemoteAddress()

	if handler, ok = h.handlers[m.GetOpCode()]; !ok {
		logger.Debug("handler for opcode %d not found, remote address: %s", m.GetOpCode(), remoteAddress)
		return
	}

	start := time.Now()

	defer func() {
		timeSpentInHandler := time.Since(start).Seconds()
		if timeSpentInHandler > 3.0 {
			logger.Warn("%s handling time exceeded %f secs for msg_id %d op-code: %d", remoteAddress, timeSpentInHandler, m.GetMessageID(), m.GetOpCode())
		}
		metrics.UpdateMessageCounterForOpCode(m.GetOpCode(), bytes)
	}()

	if fe := handler.Handle(m, ws); fe != nil {
		logger.Debug("error handling message opcode: %d, remote address %s: %s", m.GetOpCode(), remoteAddress, fe.Error())
	}
}

func (h *wsListener) resolveOnDisconnectedCb() (cb onDisconnectedCb) {
	cb = func(ws IWSClient) {
		h.registry.UnregisterClient(ws)
	}
	if h.enableMetrics {
		cb = func(ws IWSClient) {
			h.registry.UnregisterClient(ws)
			metrics.UpdateConnectedClientsGauge(false)
		}
	}
	return
}
