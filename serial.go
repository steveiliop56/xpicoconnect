package xpicoconnect

import (
	"strings"

	"go.bug.st/serial"
)

func readSerialLine(bufferSize int, port serial.Port) ([]byte, error) {
	buff := make([]byte, bufferSize)
	var result []byte
	for {
		n, err := port.Read(buff)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			break
		}
		chunk := buff[:n]
		if idx := strings.Index(string(chunk), "\r\n"); idx != -1 {
			result = append(result, chunk[:idx]...)
			break
		}
		result = append(result, chunk...)
	}
	return result, nil
}
