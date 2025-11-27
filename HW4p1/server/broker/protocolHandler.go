package main

import (
	"errors"
	"strings"
)

const (
	CMD_PUBLISH   = "PUBLISH"
	CMD_SUBSCRIBE = "SUBSCRIBE"
)

type ControlPacket struct {
	Type    string
	Topic   string
	Payload string
}

func parsePacket(raw string) (ControlPacket, error) {
	raw = strings.TrimSpace(raw)
	var parts []string = strings.SplitN(raw, "|", 3)
	if len(parts) < 2 {
		return ControlPacket{}, errors.New("invalid format: expected TYPE|TOPIC|PAYLOAD")
	}

	var packet = ControlPacket{
		Type:  parts[0],
		Topic: parts[1],
	}
	if len(parts) == 3 {
		packet.Payload = parts[2]
	}
	if packet.Type != CMD_PUBLISH && packet.Type != CMD_SUBSCRIBE {
		return ControlPacket{}, errors.New("unknown command: " + packet.Type)
	}
	return packet, nil
}
