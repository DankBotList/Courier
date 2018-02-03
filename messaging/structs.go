package messaging

import (
	"time"

	"errors"
	"reflect"

	"encoding/json"

	"github.com/satori/go.uuid"
)

// The message type
type MessageType uint

const (
	TypeMessageUnknown MessageType = iota
	TypeMessageEvent
	TypeMessageRebootEstablish
)

// Event an event identifier.
type Event uint

const (
	// Event called when a node connects
	EventConnect Event = iota

	// Event called when a node disconnects cleanly
	EventDisconnect

	// Event called when a node shuts down cleanly
	EventShutdown

	// Event called when a node is lost due to connection loss/crash
	EventLoss

	// Event called when an update is available
	EventAvailable
)

// Message an abstract, the message to be sent/received/forwarded.
type Message struct {
	// ID an ID to only forward once
	ID uuid.UUID
	// SourceID the ID of the sender so we can traceback.
	SourceID string
	// Type the type of message this is.
	Type MessageType
	// Data the JSON data
	Data []byte
	// Message the message we contain
	Message interface{}
	// Already marshaled json
	Cache []byte
}

func (m Message) MarshalJSON() ([]byte, error) {

	var ret []byte
	if m.Message == nil {
		return ret, errors.New("bad message")
	}

	m.findMessageType()
	data, err := json.Marshal(m.Message)
	if err != nil {
		return ret, err
	} else {
		m.Data = data
	}

	return json.Marshal(struct {
		ID       string      `json:"id"`
		SourceID string      `json:"source_id"`
		Type     MessageType `json:"type"`
		Data     []byte      `json:"data"`
	}{
		m.ID.String(), m.SourceID, m.Type, m.Data,
	})

}

func (m *Message) UnmarshalJSON(b []byte) error {
	msg := struct {
		ID       string      `json:"id"`
		SourceID string      `json:"source_id"`
		Type     MessageType `json:"type"`
		Data     []byte      `json:"data"`
	}{}

	if err := json.Unmarshal(b, &msg); err != nil {
		return err
	}

	if id, err := uuid.FromString(msg.ID); err != nil {
		return err
	} else {
		m.ID = id
	}

	m.SourceID = msg.SourceID
	m.Type = msg.Type
	m.Data = msg.Data

	m.createMessageType()
	return json.Unmarshal(m.Data, m.Message)

}

func (m *Message) findMessageType() {
	msg := reflect.TypeOf(m.Message)
	switch msg {
	case reflect.TypeOf(EventMessage{}):
	case reflect.TypeOf(&EventMessage{}):
		m.Type = TypeMessageEvent
		return
	case reflect.TypeOf(RebootEstablishMessage{}):
	case reflect.TypeOf(&RebootEstablishMessage{}):
		m.Type = TypeMessageRebootEstablish
		return
	default:
		m.Type = TypeMessageUnknown
		return
	}
}

func (m *Message) createMessageType() {
	switch m.Type {
	case TypeMessageEvent:
		m.Message = &EventMessage{}
		return
	case TypeMessageRebootEstablish:
		m.Message = &RebootEstablishMessage{}
		return
	default:
		m.Message = &map[string]interface{}{}
		return
	}
}

// EventMessage a message containing an event and it's data.
type EventMessage struct {
	// Event the event const that this corresponds to
	Event Event `json:"event"`
	// Any data that belongs to the event.
	Data []string `json:"data"`
}

// RebootEstablishMessage a message to establish this instances reboot timings.
type RebootEstablishMessage struct {
	// RebootAfter I will reboot after this time after the previous server has come up.
	RebootAfter time.Duration `json:"reboot_after"`
	// DeadAfter Consider me dead after this duration of not coming back up
	DeadAfter time.Duration `json:"dead_after"`
}
