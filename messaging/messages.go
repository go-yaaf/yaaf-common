// Common messaging messages
//

package messaging

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

// BaseMessage basie implementation of IMessage interface
type BaseMessage struct {
	MsgTopic     string `json:"topic"`     // Message topic (channel)
	MsgOpCode    int    `json:"opCode"`    // Message op code
	MsgAddressee string `json:"addressee"` // Message final addressee
	MsgSessionId string `json:"sessionId"` // Session id shared across all messages related to the same session
	MsgPayload   any    `json:"payload"`   // Message payload
}

func (m *BaseMessage) Topic() string     { return m.MsgTopic }
func (m *BaseMessage) OpCode() int       { return m.MsgOpCode }
func (m *BaseMessage) Addressee() string { return m.MsgAddressee }
func (m *BaseMessage) SessionId() string { return m.MsgSessionId }
func (m *BaseMessage) Payload() any      { return m.MsgPayload }

// MessageFactory is a factory method of any message
type MessageFactory func() IMessage

// SubscriptionCallback Message subscription callback function, return true for ack
type SubscriptionCallback func(msg IMessage) bool

// endregion
