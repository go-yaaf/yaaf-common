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

	// Version identifies a message version
	Version() string

	// Payload is the message body
	Payload() any
}

// BaseMessage base implementation of IMessage interface
type BaseMessage struct {
	MsgTopic     string `json:"topic"`     // Message topic (channel)
	MsgOpCode    int    `json:"opCode"`    // Message op code
	MsgVersion   string `json:"version"`   // Message op code
	MsgAddressee string `json:"addressee"` // Message final addressee
	MsgSessionId string `json:"sessionId"` // Session id shared across all messages related to the same session
}

func (m *BaseMessage) Topic() string     { return m.MsgTopic }
func (m *BaseMessage) OpCode() int       { return m.MsgOpCode }
func (m *BaseMessage) Version() string   { return m.MsgVersion }
func (m *BaseMessage) Addressee() string { return m.MsgAddressee }
func (m *BaseMessage) SessionId() string { return m.MsgSessionId }
func (m *BaseMessage) Payload() any      { return nil }

// MessageFactory is a factory method of any message
type MessageFactory func() IMessage

// SubscriptionCallback Message subscription callback function, return true for ack
type SubscriptionCallback func(msg IMessage) bool

// endregion

// Message is a generic message structure
type Message[T any] struct {
	BaseMessage
	MsgPayload T `json:"payload"` // Data model
}

func NewMessage[T any]() IMessage {
	return &Message[T]{}
}

func GetMessage[T any](topic string, payload T) IMessage {
	message := &Message[T]{
		BaseMessage: BaseMessage{
			MsgTopic:     topic,
			MsgOpCode:    0,
			MsgSessionId: entity.NanoID(),
		},
		MsgPayload: payload,
	}
	return message
}

// EntityMessage general entity message implementation of IMessage interface
type EntityMessage struct {
	BaseMessage
	MsgPayload any `json:"payload"` // Payload
}

func (m *EntityMessage) Payload() any { return m.MsgPayload }

// NewEntityMessage is a message factory
func NewEntityMessage() IMessage {
	return &EntityMessage{}
}

const EntityMessageTopic = "ENTITY"
