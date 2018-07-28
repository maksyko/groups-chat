package protocol

type RpcMethod int

const (
	RPC_MESSAGE_GET       RpcMethod = 20
	RPC_MESSAGE_SEND      RpcMethod = 40
	RPC_MESSAGE_DELIVERED RpcMethod = 41
	RPC_MESSAGE_READ      RpcMethod = 42
	//RPC_MESSAGE_UPDATED   RpcMethod = 43
	//RPC_MESSAGE_DELETED   RpcMethod = 44
	RPC_TYPING_START      RpcMethod = 60
	RPC_TYPING_END        RpcMethod = 61
)

type RPC struct {
	Method RpcMethod   `json:"method"`
	Body   interface{} `json:"body,omitempty"`
}

type RpcMessageGet struct {
	Timestamp int64 `json:"timestamp,omitempty"`
}

type RpcMessageSend struct {
	GroupID string `json:"group_id,omitempty"`
	Data    string `json:"data,omitempty"`
}

type RpcMessageDelivered struct {
	MessageID string `json:"message_id,omitempty"`
}

type RpcMessageRead struct {
	MessageID string `json:"message_id,omitempty"`
}

type RpcTypingStart struct {
	GroupID string `json:"group_id,omitempty"`
}

type RpcTypingEnd struct {
	GroupID string `json:"group_id,omitempty"`
}

// =====================================================================================================================

type EventType int

const (
	EVENT_MESSAGE           EventType = 20
	EVENT_MESSAGE_SENT      EventType = 21
	EVENT_MESSAGE_DELIVERED EventType = 22
	EVENT_MESSAGE_READ      EventType = 23
	//EVENT_MESSAGE_UPDATED   EventType = 24
	//EVENT_MESSAGE_DELETED   EventType = 25
	EVENT_TYPING_START      EventType = 40
	EVENT_TYPING_END        EventType = 41
	EVENT_GROUP             EventType = 70
	//EVENT_GROUP_UPDATED     EventType = 71
	EVENT_GROUP_JOINED      EventType = 72
	EVENT_GROUP_LEFT        EventType = 73
)

type Event struct {
	Type      EventType   `json:"type"`
	Timestamp int64       `json:"timestamp,omitempty"`
	Body      interface{} `json:"body,omitempty"`
}

type EventMessage struct {
	MessageID string `json:"message_id,omitempty"`
	Data      string `json:"data,omitempty"`
}

type EventMessageSent struct {
	MessageID string `json:"message_id,omitempty"`
}

type EventMessageDelivered struct {
	MessageID string `json:"message_id,omitempty"`
}

type EventTypingStart struct {
	GroupID string `json:"group_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}

type EventTypingEnd struct {
	GroupID string `json:"group_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}

type EventGroup struct {
	GroupID string   `json:"group_id,omitempty"`
	Name    string   `json:"name,omitempty"`
	UserIDs []string `json:"user_ids,omitempty"`
}

type EventGroupJoined struct {
	GroupID string `json:"group_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}

type EventGroupLeft struct {
	GroupID string `json:"group_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}
