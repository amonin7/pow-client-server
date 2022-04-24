package model

import (
	"fmt"
	"strconv"
	"strings"
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

// Message - message struct for both server and client
type Message struct {
	Header  int    //type of message
	Payload string //payload, could be json, quote or be empty
}

// DeserializeMessage - parses Message from str, checks header and payload
func DeserializeMessage(str string) (*Message, error) {
	//str = strings.TrimSpace(str)
	var msgType int
	// message has view as 1|payload (payload is optional)
	parts := strings.Split(str, "|")
	if len(parts) < 1 || len(parts) > 2 { //only 1 or 2 parts allowed
		return nil, fmt.Errorf("message doesn't match protocol")
	}
	// try to parse header
	msgType, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("cannot parse header")
	}
	msg := Message{
		Header:  msgType,
		Payload: parts[1],
	}
	return &msg, nil
}

func (m *Message) Serialize() string {
	return fmt.Sprintf("%d|%s", m.Header, m.Payload)
}
