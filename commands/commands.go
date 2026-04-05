package commands

import (
	"fmt"
	"strings"
)

func EncodeCommand(command string, value []byte) []byte {
	var sb strings.Builder
	sb.WriteString(command)
	sb.WriteString(":")
	if len(value) == 0 {
		sb.WriteString("null")
	} else {
		sb.WriteString(string(value))
	}
	sb.WriteString("\n")
	return []byte(sb.String())
}

func DecodeCommand(data []byte) (string, string, error) {
	parts := strings.SplitN(string(data), ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid command, expected command:value, got %v", string(data))
	}
	return parts[0], strings.TrimSpace(parts[1]), nil
}

func DecodeResponse(res []byte) (string, error) {
	parts := strings.SplitN(string(res), ":", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid response, expected command:status:result, got %v", string(res))
	}
	if parts[1] != "ok" {
		return "", fmt.Errorf("command failed with status: %s, result: %s", parts[1], parts[2])
	}
	return strings.TrimSpace(parts[2]), nil
}

func EncodeResponse(command string, status string, result []byte) []byte {
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
	sb.WriteString("\n")
	return []byte(sb.String())
}
