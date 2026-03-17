package main

import (
	"log"
	"strings"

	"go.bug.st/serial"
)

type PicoBridgeConfig struct {
	PortName   string
	BaudRate   int
	BufferSize int
}

type PicoBridge struct {
	port serial.Port
	cfg  PicoBridgeConfig
	ch   chan []byte
}

func NewPicoBridge(cfg PicoBridgeConfig) (*PicoBridge, error) {
	mode := &serial.Mode{
		BaudRate: cfg.BaudRate,
	}

	port, err := serial.Open(cfg.PortName, mode)

	if err != nil {
		return nil, err
	}

	return &PicoBridge{
		port: port,
		cfg:  cfg,
	}, nil
}

func (pb *PicoBridge) readLine() ([]byte, error) {
	buff := make([]byte, pb.cfg.BufferSize)
	totalLen := 0
	for {
		n, err := pb.port.Read(buff)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			break
		}
		if strings.Contains(string(buff[:n]), "\r\n") {
			totalLen += n - 2
			break
		}
		totalLen += n
	}
	return buff[:totalLen], nil
}

func (pb *PicoBridge) StartReader() chan []byte {
	pb.ch = make(chan []byte, pb.cfg.BufferSize)
	go func() {
		for {
			line, err := pb.readLine()
			if err != nil {
				log.Printf("failed to read from pico: %v", err)
				continue
			}
			pb.ch <- line
		}
	}()
	return pb.ch
}

func (pb *PicoBridge) GetChannel() chan []byte {
	return pb.ch
}

func (pb *PicoBridge) Write(data []byte) error {
	_, err := pb.port.Write(data)
	return err
}

func (pb *PicoBridge) Close() error {
	return pb.port.Close()
}
