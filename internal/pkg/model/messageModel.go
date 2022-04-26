package model

import (
	"encoding/json"
	"fmt"
)

// Represents the messageType
const (
	// ChallengeReq - the request for the challenge
	ChallengeReq = 1
	// ChallengeRes - the response for the provided challenge
	ChallengeRes = 2
	// ResourceReq -  the request for some resource
	ResourceReq = 3
	// ResourceRes -  the response with the requested resource
	ResourceRes = 4
)

// Message - structure of the message, which client and server are exchanged with.
type Message struct {
	//Represents the messageType, described above
	Header  int    `json:"header"`
	Payload string `json:"payload"`
}

// DeserializeMessage - deserializes Message from the array of bytes
func DeserializeMessage(serializedMessage string) (*Message, error) {
	var message Message
	err := json.Unmarshal([]byte(serializedMessage), &message)
	if err != nil {
		return nil, fmt.Errorf("cannot parse message")
	}
	return &message, nil
}

// Serialize - serializes Message to the array of bytes
func (m *Message) Serialize() ([]byte, error) {
	str, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize message")
	}
	return str, nil
}
