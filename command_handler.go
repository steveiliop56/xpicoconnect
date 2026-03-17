package main

import (
	"fmt"
	"strings"
)

type CommandHandler struct {
	pb             *PicoBridge
	commandPending bool
}

func NewCommandHandler(pb *PicoBridge) *CommandHandler {
	return &CommandHandler{
		pb: pb,
	}
}

func (ch *CommandHandler) EncodeCommand(command string, value []byte) []byte {
	var sb strings.Builder
	sb.WriteString(command)
	sb.WriteString(":")
	if len(value) == 0 {
		sb.WriteString("null")
	} else {
		sb.WriteString(string(value))
	}
	return []byte(sb.String())
}

func (ch *CommandHandler) DecodeCommand(data []byte) (string, string, error) {
	parts := strings.SplitN(string(data), ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid command, expected command:value, got %v", string(data))
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func (ch *CommandHandler) EncodeResponse(command string, status string, result []byte) []byte {
	var sb strings.Builder
	sb.WriteString(command)
	sb.WriteString(":")
	sb.WriteString(status)
	sb.WriteString(":")
	if len(result) == 0 {
		sb.WriteString("null")
	} else {
		sb.WriteString(string(result))
	}
	return []byte(sb.String())
}

func (ch *CommandHandler) DecodeResponse(res []byte) (string, error) {
	parts := strings.SplitN(string(res), ":", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid response, expected command:status:result, got %v", string(res))
	}
	if parts[1] != "ok" {
		return "", fmt.Errorf("command failed with status: %s, result: %s", parts[1], parts[2])
	}
	return strings.TrimSpace(parts[2]), nil
}

func (ch *CommandHandler) SendCommand(command string, value []byte) (string, error) {
	encoded := ch.EncodeCommand(command, value)
	err := ch.pb.Write(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to write command to: %v", err)
	}
	ch.commandPending = true
	res := <-ch.pb.GetChannel()
	ch.commandPending = false
	return ch.DecodeResponse(res)
}

func (ch *CommandHandler) CommandPending() bool {
	return ch.commandPending
}
