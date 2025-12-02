// Common messaging messages
//

package messaging

import "github.com/go-yaaf/yaaf-common/entity"

// region Message interface --------------------------------------------------------------------------------------------

// IMessage defines the interface for a standard message used within the messaging system.
// It provides a structured way to define messages with common attributes like topic, operation code, and payload.
type IMessage interface {
	// Topic returns the name of the topic, channel, or queue the message is associated with.
	Topic() string

	// OpCode returns the message's operational code, which can be used to indicate the type or purpose of the message.
	OpCode() int

	// Addressee returns the final recipient of the message. This is an optional field.
	Addressee() string

	// SessionId returns an identifier for a message exchange session. This ID is shared across all messages
	// belonging to the same session, allowing for conversational message patterns.
	SessionId() string

	// Version returns the version of the message, which can be used for compatibility and evolution of message formats.
	Version() string

	// Payload returns the body of the message, which can be of any type.
	Payload() any
}

// BaseMessage provides a basic implementation of the IMessage interface.
// It can be embedded in other message structs to provide default behavior for the common message attributes.
// @Data
type BaseMessage struct {
	MsgTopic     string `json:"topic"`     // The topic (channel or queue) of the message.
	MsgOpCode    int    `json:"opCode"`    // An operational code for the message.
	MsgVersion   string `json:"version"`   // The version of the message.
	MsgAddressee string `json:"addressee"` // The final recipient of the message.
	MsgSessionId string `json:"sessionId"` // A session ID for tracking related messages.
}

// Topic returns the message's topic.
func (m *BaseMessage) Topic() string { return m.MsgTopic }

// OpCode returns the message's operational code.
func (m *BaseMessage) OpCode() int { return m.MsgOpCode }

// Version returns the message's version.
func (m *BaseMessage) Version() string { return m.MsgVersion }

// Addressee returns the message's final addressee.
func (m *BaseMessage) Addressee() string { return m.MsgAddressee }

// SessionId returns the message's session ID.
func (m *BaseMessage) SessionId() string { return m.MsgSessionId }

// Payload returns nil, as the BaseMessage does not carry a payload itself.
// This method is intended to be overridden by embedding structs.
func (m *BaseMessage) Payload() any { return nil }

// MessageFactory is a function type that serves as a factory for creating IMessage instances.
// It is often used in consumers to generate the correct message type for unmarshalling.
type MessageFactory func() IMessage

// SubscriptionCallback is a function type for handling messages received from a subscription.
// It takes an IMessage as input and returns a boolean indicating if the message was successfully processed (acknowledged).
type SubscriptionCallback func(msg IMessage) bool

// endregion

// Message is a generic message structure that embeds BaseMessage and adds a typed payload.
//
// Type Parameters:
//
//	T: The type of the payload.
//
// @Data
type Message[T any] struct {
	BaseMessage
	MsgPayload T `json:"payload"` // The typed data payload of the message.
}

// NewMessage creates a new instance of a generic Message.
//
// Type Parameters:
//
//	T: The type of the payload.
//
// Returns:
//
//	A new IMessage with a typed payload.
func NewMessage[T any]() IMessage {
	return &Message[T]{}
}

// GetMessage creates and initializes a new generic Message with a given topic and payload.
// It automatically assigns a new session ID.
//
// Type Parameters:
//
//	T: The type of the payload.
//
// Parameters:
//
//	topic: The topic for the message.
//	payload: The payload for the message.
//
// Returns:
//
//	A new, initialized IMessage.
func GetMessage[T any](topic string, payload T) IMessage {
	return &Message[T]{
		BaseMessage: BaseMessage{
			MsgTopic:     topic,
			MsgOpCode:    0,
			MsgSessionId: entity.NanoID(),
		},
		MsgPayload: payload,
	}
}

// EntityMessage is a message structure for carrying generic entity data.
// It embeds BaseMessage and uses an `any` type for the payload.
// @Data
type EntityMessage struct {
	BaseMessage
	MsgPayload any `json:"payload"` // The payload of the message, typically an entity.
}

// Payload returns the message's payload.
func (m *EntityMessage) Payload() any { return m.MsgPayload }

// NewEntityMessage is a factory function that creates a new instance of EntityMessage.
//
// Returns:
//
//	A new IMessage of type EntityMessage.
func NewEntityMessage() IMessage {
	return &EntityMessage{}
}

// EntityMessageTopic is a predefined topic name for messages that carry entities.
const EntityMessageTopic = "ENTITY"
