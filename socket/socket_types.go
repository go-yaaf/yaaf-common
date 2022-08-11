// Copyright 2022. Motty Cohen
//
// package types
//
package socket

import (
	"encoding/json"
	"net/http"
)

const (
	WsPingOpCode = 0
)

// IWSMessage is a Web socket message header interface:
type IWSMessage interface {
	GetOpCode() int       // Get message op-code
	GetMessageID() uint64 // Get message unique Id
	GetSessionID() string // Get session Id
	GetPayload() any      // Get arbitrary message payload
}

// region Web Socket message header ------------------------------------------------------------------------------------

type WSMessageHeader struct {
	OpCode    int
	MessageID uint64
	SessionID string
}

// Get op-code
func (mb WSMessageHeader) GetOpCode() int { return mb.OpCode }

// Get message ID
func (mb WSMessageHeader) GetMessageID() uint64 { return mb.MessageID }

// Get session ID
func (mb WSMessageHeader) GetSessionID() string { return mb.SessionID }

// endregion

// region Web Socket Ping Pong messages --------------------------------------------------------------------------------

// PONG messaged handler received from client
type PongReceivedCb func(sessionId, pongMessage string, latencyMs int64)

// PING message
type WSPingMessage struct {
	WSMessageHeader
}

// PING message has no payload
func (mp WSPingMessage) GetPayload() interface{} { return nil }

func NewWsPingMessage() IWSMessage {
	return WSPingMessage{WSMessageHeader: WSMessageHeader{OpCode: WsPingOpCode}}
}

var pingMessage = NewWsPingMessage()

// endregion

// region Web Socket Raw message ---------------------------------------------------------------------------------------

type WSRawMessage struct {
	WSMessageHeader
	Body []byte
}

func (m *WSRawMessage) GetPayload() any { return m.Body }

// endregion

// region Web Socket message decoder -----------------------------------------------------------------------------------

type IMessageDecoder interface {
	Encode(message IWSMessage) ([]byte, error)
	Decode(buffer []byte) (IWSMessage, error)
}

// endregion

// region Web Socket client --------------------------------------------------------------------------------------------

type WSConnectParams struct {
	Url                string      // Full url (int is case path and host are ignored)
	Path               string      // URL path segment
	Host               string      // url host + port
	WriteBufferSize    int         // Write buffer size (if not provided use the default 8K buffer)
	ReadBufferSize     int         // Read buffer size (if not provided use the default 8K buffer)
	CompressionEnabled bool        // Tru to enable compression
	Header             http.Header // List of HTTP headers
}

// Web socket client interface
type IWSClient interface {
	ID() string                            // Socket client unique ID
	QueryParams() map[string]string        // Query parameters from the REST call before protocol upgrade
	Connect(p WSConnectParams) error       // Connect to server
	Send(m IWSMessage) error               // Send message through the socket
	SendRaw(m []byte) error                // Send arbitrary data through the socket
	Close() error                          // Close connection
	PongReceivedHandler(cb PongReceivedCb) // PONG message receive handler function
}

// endregion

// region Message factory and default message decoder (JSON) -----------------------------------------------------------

type MessageFactoryFunc func() IWSMessage

var messageFactories = map[int]MessageFactoryFunc{}

func AddMessageFactory(opcode int, f MessageFactoryFunc) {
	messageFactories[opcode] = f
}

func GetMessageFactoryFunc(opcode int) MessageFactoryFunc {
	return messageFactories[opcode]
}

type JsonDecoder struct{}

func NewJsonDecoder() IMessageDecoder {
	return &JsonDecoder{}
}

func (_ JsonDecoder) Encode(m IWSMessage) (result []byte, err error) {
	return json.Marshal(m)
}

func (_ JsonDecoder) Decode(buffer []byte) (msg IWSMessage, err error) {

	bm := &WSMessageHeader{}

	if err = json.Unmarshal(buffer, bm); err != nil {
		return nil, err
	}

	if mf, ok := messageFactories[bm.GetOpCode()]; ok {
		msg = mf()
		if err = json.Unmarshal(buffer, msg); err != nil {
			return nil, err
		}
	} else {
		msg = &WSRawMessage{
			WSMessageHeader: WSMessageHeader{
				OpCode:    bm.GetOpCode(),
				MessageID: bm.GetMessageID(),
				SessionID: bm.GetSessionID(),
			},
			Body: buffer,
		}
	}
	return
}

// endregion

// Web socket message handler
type IWSMessageHandler interface {
	Handle(m IWSMessage, rw IWSClient) error
}

// Web socket message handler entry
type WSHandlerEntry struct {
	OpCode  int
	Handler IWSMessageHandler
	Metrics bool
}

// Web socket client registry
type IWSClientRegistry interface {
	Start()
	RegisterClient(c IWSClient)
	UnregisterClient(c IWSClient)
	ConnectedClients() int
}

// Web socket endpoint configuration
type WSEndpointConfig struct {
	Path              string            // Web socket endpoint path
	Handlers          []WSHandlerEntry  // Web socket handlers
	Registry          IWSClientRegistry // Web socket client registry
	CustomDecoder     IMessageDecoder   // Custom message decoder
	EnablePingPong    bool              // Flag to enable socket PING-PONG
	EnableMetrics     bool              // Flag to enable metrics
	KeepAliveInterval int               // Keep alive interval in seconds (-1 for no keep-alive)
}

// region Web Socket endpoints config ----------------------------------------------------------------------------------

type IWebSocketEndpoint interface {
	Entries() []WSEndpointConfig
}

type WebSocketEndpoint struct {
	entries []WSEndpointConfig
}

func (r WebSocketEndpoint) Entries() []WSEndpointConfig { return r.entries }

func NewWebSocketEndpoint(entries []WSEndpointConfig) WebSocketEndpoint {
	return WebSocketEndpoint{entries: entries}
}

// endregion
