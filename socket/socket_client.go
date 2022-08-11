// Copyright 2022. Motty Cohen
//
// Web socket client implementation
//
package socket

import (
	"encoding/json"
	"fmt"
	"github.com/agentvi/innovi-core-commons/config"
	. "github.com/agentvi/innovi-core-commons/facility_error"
	"github.com/agentvi/innovi-core-commons/instrumenting"
	"github.com/agentvi/innovi-core-commons/utils"
	. "github.com/agentvi/innovi-core-commons/utils/logger"
	"github.com/gorilla/websocket"
	"net"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	defaultReadWriteBufferSize = 8 * 1024
	maxMessageSize             = 1024 * 1024 * 5
)

// region Web Socket client callbacks ----------------------------------------------------------------------------------

// On socket disconnected callback signature
type onDisconnectedCb func(IWSClient)

// On socket connected callback signature
type onConnectedCb func(IWSClient)

// On message received callback signature
type onMessageReceivedCb func(IWSClient, IWSMessage, int)

// endregion

// region Web Socket client structure and fluent API configuration -----------------------------------------------------

// Wev socket client
type wsClient struct {
	id                string                    // Web socket client unique ID
	uri               string                    // Web socket client URI
	queryParams       map[string]string         // URI query parameters map
	decoder           IMessageDecoder           // Message decoder
	onMessageReceived onMessageReceivedCb       // Hook to message received callback
	onDisconnected    onDisconnectedCb          // Hook for socket disconnected callback
	onConnected       onConnectedCb             // Hook for socket connected callback
	onPongReceived    PongReceivedCb            // Hook for PONG message callback
	handlers          map[int]IWSMessageHandler // Map of op-code to message handler
	pingPongEnabled   bool                      // Enable PING-PONG
	closed            bool
	closeGuard        sync.RWMutex
	wcg               chan bool
	conn              *websocket.Conn // Pointer to the underlying web socket connection

}

func NewWsClient()

type WsClientConfig struct {
	Id             string
	WsConn         *websocket.Conn
	OnMsgRvd       onMessageReceivedCb
	OnDisconnected onDisconnectedCb
	OnConnected    onConnectedCb
	MessageDecoder IMessageDecoder
	PinPongEnabled bool
	QueryParams    map[string]string
	Handlers       map[int]IWSMessageHandler
}

func NewWsClient(id string,
	conn *websocket.Conn,
	onMsgRvd onMessageReceivedCb,
	onDisconnected onDisconnectedCb,
	decoder IMessageDecoder,
	lp string,
	qParams map[string]string) IWSClient {

	cfg := WsClientConfig{
		WsConn:         conn,
		OnMsgRvd:       onMsgRvd,
		OnDisconnected: onDisconnected,
		MessageDecoder: decoder,
		PinPongEnabled: false,
		QueryParams:    qParams,
		Handlers:       nil,
	}

	return NewWsClient2(cfg)
}

func NewWsClient2(cfg WsClientConfig) IWSClient {

	ws := &wsClient{
		conn:              cfg.WsConn,
		onMessageReceived: cfg.OnMsgRvd,
		onDisconnected:    cfg.OnDisconnected,
		onConnected:       cfg.OnConnected,
		id:                cfg.Id,
		decoder:           cfg.MessageDecoder,
		pingPongEnabled:   cfg.PinPongEnabled,
		queryParams:       cfg.QueryParams,
		closed:            false,
		wcg:               make(chan bool, 1),
		handlers:          cfg.Handlers}

	if cfg.MessageDecoder == nil {
		ws.decoder = NewJsonDecoder()
	}

	if ws.conn != nil {
		ws.run()
	}
	return ws
}

func (c *wsClient) ID() string {
	return c.id
}

func (c *wsClient) QueryParams() map[string]string {
	if c.queryParams == nil {
		c.queryParams = make(map[string]string)
	}
	return c.queryParams
}

func (c *wsClient) PongReceivedHandler(cb PongReceivedCb) {
	c.onPongReceived = cb
}

func (c *wsClient) Send(r IWSMessage) error {

	var (
		err    error
		buffer []byte
	)

	<-c.wcg
	defer func() {
		c.wcg <- true
	}()

	defer utils.RecoverAll(func(err interface{}) {
		Error("%s wsClient::writeResponse panic: %s", c.id, err)
		Flush()
	})

	if c.isClosed() {
		Debug("%s closed, response discarded", c.id)
		return nil
	}

	deadLine := time.Now().Add(time.Second * time.Duration(config.GetBaseConfig().WsWriteTimeoutSec()))
	if r.GetOpCode() == WsPingOpCode {

		err := c.conn.WriteControl(websocket.PingMessage, []byte(fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))), deadLine)
		if err != nil {
			Error("[%s] error sending ping\n%s", c.id, err)
		} else {
			Trace("[%s] PING sent", c.id)
		}
		return err
	}

	if buffer, err = json.Marshal(r); err == nil {
		c.conn.SetWriteDeadline(deadLine)
		if err = c.conn.WriteMessage(websocket.BinaryMessage, buffer); err != nil {
			Error("[%s] Send to client\n%s", c.id, err)
		}
	} else {
		Error("[%s] error marshalling message when trying to send\n%s", c.id, err)
	}
	return err
}

func (c *wsClient) SendRaw(buffer []byte) error {

	var (
		err error
	)

	<-c.wcg
	defer func() {
		c.wcg <- true
	}()

	defer utils.RecoverAll(func(err interface{}) {
		Error("%s wsClient::writeResponse panic: %s", c.id, err)
		Flush()
	})

	if c.isClosed() {
		Debug("%s closed, response discarded", c.id)
		return nil
	}

	deadLine := time.Now().Add(time.Second * time.Duration(config.GetBaseConfig().WsWriteTimeoutSec()))
	c.conn.SetWriteDeadline(deadLine)
	if err = c.conn.WriteMessage(websocket.BinaryMessage, buffer); err != nil {
		Error("[%s] error send to client %s", c.id, err)
	}
	return err
}

func (c *wsClient) ConnectUrl(url string) (err error) {
	return c.Connect(WSConnectParams{Url: url})
}

func (c *wsClient) Connect(p WSConnectParams) (err error) {

	if c.conn != nil {
		return fmt.Errorf("try to connect an already connected client")
	}

	wsUrl := p.Url

	// If Url parameter is provided use it, otherwise use port and host
	if len(wsUrl) == 0 {
		u := url.URL{Scheme: "ws", Host: p.Host, Path: p.Path}
		wsUrl = u.String()
	}

	dialer := websocket.DefaultDialer
	dialer.EnableCompression = p.CompressionEnabled

	if p.ReadBufferSize == 0 {
		dialer.ReadBufferSize = defaultReadWriteBufferSize
	} else {
		dialer.ReadBufferSize = p.ReadBufferSize
	}

	if p.WriteBufferSize == 0 {
		dialer.WriteBufferSize = defaultReadWriteBufferSize
	} else {
		dialer.WriteBufferSize = p.WriteBufferSize
	}

	if c.conn, _, err = dialer.Dial(wsUrl, p.Header); err == nil {
		c.conn.EnableWriteCompression(p.CompressionEnabled)
		if tcpConn, ok := c.conn.UnderlyingConn().(*net.TCPConn); ok {
			_ = tcpConn.SetWriteBuffer(config.GetBaseConfig().WsWriteBufferSizeBytes())
			_ = tcpConn.SetReadBuffer(config.GetBaseConfig().WsReadBufferSizeBytes())
		}
		c.run()
		if c.onConnected != nil {
			c.onConnected(c)
		}
		c.closed = false
	}
	return err
}

func (c *wsClient) Close() error {
	c.closeConn()
	return nil
}

func (c *wsClient) RemoteAddress() (ra string) {
	if c.conn != nil {
		ra = c.conn.RemoteAddr().String()
	}
	return
}

func (c *wsClient) readPump() {

	defer utils.RecoverAll(func(err interface{}) {
		Error("wsClient::readPump error: %s", err)
		Flush()
	})

	c.conn.SetReadLimit(maxMessageSize)

	if c.pingPongEnabled {
		_ = c.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(config.GetBaseConfig().WsPongTimeoutSec()*2)))
		go c.pinPong()
	} else {
		_ = c.conn.SetReadDeadline(time.Time{})
	}

LOOP:
	for {
		msgType, rawMessage, err := c.conn.ReadMessage()

		if err != nil {
			Debug("[%s] READ Client: %s. last message(or part of it: %s)", c.ID(), err, string(rawMessage))
			break LOOP
		} else {
			Trace("received  message type: %d from %s:\n%s", msgType, c.ID(), string(rawMessage))

			msg, fe := c.decoder.Decode(rawMessage)
			if fe != nil && fe.ErrorCode() == JsonErr {
				Error("error decoding received message from: [%s]: error: %s message dump: %s", c.Id(), fe.ErrorMessage(), string(rawMessage))
			} else if fe != nil && fe.ErrorCode() == WsNoMessageFactoryFound {
				if msg != nil {
					Trace("no message factory found for op-code %d. handler for the op-code will be invoked with raw message (if handler exists) ", msg.OpCode())
					go c.handleMessageReceived(msg, len(rawMessage))
				}
			} else {
				go c.handleMessageReceived(msg, len(rawMessage))
			}
		}
	}

	c.closeConn()
	Debug("[%s] READ pump exit.", c.Id())
}
func (c *wsClient) pinPong() {
	pingTicker := time.NewTicker(3 * time.Second)
	for !c.isClosed() {
		select {
		case <-pingTicker.C:
			_ = c.Send(pingMessage)
		}
	}
	pingTicker.Stop()
	Debug("[%s] ping ticker stopped", c.Id())
}

func (c *wsClient) run() {
	c.wcg <- true
	if c.pingPongEnabled {
		c.conn.SetPongHandler(func(s string) error {
			ts, _ := strconv.ParseInt(s, 10, 64)
			tsNow := time.Now().UnixNano() / int64(time.Millisecond)
			latencyMs := tsNow - ts
			Trace("[%s] PONG control message received. ws latency: %d ms", c.id, latencyMs)
			_ = c.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(config.GetBaseConfig().WsPongTimeoutSec())))

			if c.onPongReceived != nil {
				go c.onPongReceived(c.id, s, latencyMs)
			}
			return nil
		})
	}
	go c.readPump()
}

func (c *wsClient) closeConn() {

	if c.isClosed() {
		return
	}

	c.setIsClosed()

	Debug("%s closeConn invoked", c.id)

	if err := c.conn.Close(); err != nil {
		Error("%s error when closing:\n%s", err)
	}

	if c.onDisconnected != nil {
		c.onDisconnected(c)
	}

	c.onDisconnected = nil
	c.onMessageReceived = nil
	c.onPongReceived = nil
}

func (c *wsClient) isClosed() (v bool) {
	c.closeGuard.Lock()
	defer c.closeGuard.Unlock()

	v = c.closed
	return
}

func (c *wsClient) setIsClosed() (v bool) {
	c.closeGuard.Lock()
	defer c.closeGuard.Unlock()

	c.closed = true
	return
}

//  if onMessageReceived callback provided, invoke it
func (c *wsClient) handleMessageReceived(m IWSMessage, sizeInBytes int) {

	invokeHandlerWrapper := func(handler IWSMessageHandler) {
		start := time.Now()

		defer func() {
			timeSpentInHandler := time.Since(start).Seconds()
			if timeSpentInHandler > 3.0 {
				Warn("%s handling time exceeded %f secs for msg_id %d op-code: %d", c.RemoteAddress(), timeSpentInHandler, m.GetMessageID(), m.GetOpCode())
			}
			instrumenting.UpdateMessageCounterForOpCode(m.GetOpCode(), sizeInBytes)
		}()

		if fe := handler.Handle(m, c); fe != nil {
			Debug("error handling message opcode: %d, remote address %s\n %s", m.GetOpCode(), c.RemoteAddress(), fe.Error())
		}
	}

	if c.handlers != nil && len(c.handlers) > 0 {
		if v, ok := c.handlers[m.GetOpCode()]; ok {
			invokeHandlerWrapper(v)
			return
		}
	} else {
		if c.onMessageReceived != nil {
			c.onMessageReceived(c, m, sizeInBytes)
		}
	}

}
