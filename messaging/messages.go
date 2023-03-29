// Common messaging messages
//

package messaging

import "github.com/go-yaaf/yaaf-common/entity"

// region Message interface --------------------------------------------------------------------------------------------

// IMessage General message interface
type IMessage interface {
	// Topic name (also known as channel or queue)
	Topic() string

	// OpCode message operational code
	OpCode() int

	// Addressee message final addressee (recipient) - optional field
	Addressee() string

	// SessionId identifies a message exchange session which is shared across all messages related to the same session
	SessionId() string

	// Payload is the message body
	Payload() any
}

// BaseMessage base implementation of IMessage interface
type BaseMessage struct {
	MsgTopic     string `json:"topic"`     // Message topic (channel)
	MsgOpCode    int    `json:"opCode"`    // Message op code
	MsgAddressee string `json:"addressee"` // Message final addressee
	MsgSessionId string `json:"sessionId"` // Session id shared across all messages related to the same session
}

func (m *BaseMessage) Topic() string     { return m.MsgTopic }
func (m *BaseMessage) OpCode() int       { return m.MsgOpCode }
func (m *BaseMessage) Addressee() string { return m.MsgAddressee }
func (m *BaseMessage) SessionId() string { return m.MsgSessionId }
func (m *BaseMessage) Payload() any      { return nil }

// MessageFactory is a factory method of any message
type MessageFactory func() IMessage

// SubscriptionCallback Message subscription callback function, return true for ack
type SubscriptionCallback func(msg IMessage) bool

// endregion

// EntityMessage generic implementation of IMessage interface
type EntityMessage struct {
	MsgTopic     string        `json:"topic"`     // Message topic (channel)
	MsgOpCode    int           `json:"opCode"`    // Message op code
	MsgAddressee string        `json:"addressee"` // Message final addressee
	MsgSessionId string        `json:"sessionId"` // Session id shared across all messages related to the same session
	MsgPayload   entity.Entity `json:"payload"`   // Payload
}

func (m *EntityMessage) Topic() string     { return m.MsgTopic }
func (m *EntityMessage) OpCode() int       { return m.MsgOpCode }
func (m *EntityMessage) Addressee() string { return m.MsgAddressee }
func (m *EntityMessage) SessionId() string { return m.MsgSessionId }
func (m *EntityMessage) Payload() any      { return m.MsgPayload }
