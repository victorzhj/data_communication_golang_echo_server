package main

import (
	"errors"
	"strings"
)

const (
	CMD_PUBLISH   = "PUBLISH"
	CMD_SUBSCRIBE = "SUBSCRIBE"
	CMD_BACKUP    = "BACKUP"
	CMD_CLEAR     = "CLEAR"
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
	switch packet.Type {
	case CMD_PUBLISH, CMD_SUBSCRIBE, CMD_BACKUP, CMD_CLEAR:
		return packet, nil
	default:
		return ControlPacket{}, errors.New("invalid packet type " + packet.Type)
	}
}
