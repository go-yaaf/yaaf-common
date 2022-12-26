// Copyright 2022. Shield-IoT Ltd.
//
// Common messaging messages
//

package messaging

// region Queue Message ------------------------------------------------------------------------------------------------

// IQueueMessage queue message interface
type IQueueMessage interface {
	// Queue is the name of the message queue (or List)
	Queue() string

	// Payload is the message body
	Payload() any
}

// QueueMessage is a basic implementation of IQueueMessage
type QueueMessage struct {
	MsgQueue   string `json:"queue"`   // Message destination queue
	MsgPayload any    `json:"payload"` // Message payload
}

func (m *QueueMessage) Queue() string { return m.MsgQueue }
func (m *QueueMessage) Payload() any  { return m.MsgPayload }

// QueueMessageFactory is a factory method of queue message
type QueueMessageFactory func() IQueueMessage

// endregion

// region PubSub Message -----------------------------------------------------------------------------------------------

// IPubSubMessage Publish-Subscribe message interface
type IPubSubMessage interface {
	// Topic name (channel)
	Topic() string

	// OpCode message operational code
	OpCode() int

	// Addressee message final addressee (recipient) - optional field
	Addressee() string

	// SessionId identifies a message exchange session which is shared across all messages related to the same session
	SessionId() string
}

// PubSubMessage basic implementation of IPubSubMessage message interface
type PubSubMessage struct {
	MsgTopic     string `json:"topic"`     // Message topic (channel)
	MsgOpCode    int    `json:"opCode"`    // Message op code
	MsgAddressee string `json:"addressee"` // Message final addressee
	MsgSessionId string `json:"sessionId"` // Session id shared across all messages related to the same session
}

func (m *PubSubMessage) Topic() string     { return m.MsgTopic }
func (m *PubSubMessage) OpCode() int       { return m.MsgOpCode }
func (m *PubSubMessage) Addressee() string { return m.MsgAddressee }
func (m *PubSubMessage) SessionId() string { return m.MsgSessionId }

// PubSubMessageFactory is a factory method of pub-sub message
type PubSubMessageFactory func() IPubSubMessage

// endregion

// SubscriptionCallback Message subscription callback function
type SubscriptionCallback func(msg IPubSubMessage)
