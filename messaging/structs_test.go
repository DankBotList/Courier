package messaging

import (
	"testing"

	"encoding/json"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var msgToMarshal = Message{
	uuid.Must(uuid.FromString("c26dbfba-b18e-4619-8039-b82e91eb3143")),
	"Source",
	TypeMessageUnknown,
	[]byte{0x7b, 0x22, 0x6b, 0x65, 0x79, 0x22, 0x3a, 0x22, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x7d},
	&map[string]interface{}{"key": "value"},
	nil,
}

var inputJson = `{"id":"c26dbfba-b18e-4619-8039-b82e91eb3143","source_id":"Source","type":0,"data":"eyJrZXkiOiJ2YWx1ZSJ9"}`

// Test marshalling json.
func TestMessage_MarshalJSON(t *testing.T) {
	data, err := json.Marshal(msgToMarshal)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, inputJson, string(data))
}

// Test unmarshalling json.
func TestMessage_UnmarshalJSON(t *testing.T) {
	m := &Message{}
	err := json.Unmarshal([]byte(inputJson), m)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, &msgToMarshal, m)
}
