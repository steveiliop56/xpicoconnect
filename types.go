package xpicoconnect

import (
	xphttpbridgego "github.com/steveiliop56/xphttpbridge-go"
	"go.bug.st/serial"
)

type SerialConfig struct {
	Baudrate   int
	Port       string
	BufferSize int
	Timeout    int
}

type XPHTTPBridgeConfig struct {
	Address string
	Port    int
}

type XPicoConnectorConfig struct {
	SerialConfig       SerialConfig
	XPHTTPBridgeConfig XPHTTPBridgeConfig
	PollTime           int
}

type XPicoConnectorState struct {
	isCommandPending   bool
	picoCommandBinders map[string]PicoCommandBind
	bridgeRefBinders   map[string]BridgeRefBind
}

type XPicoConnector struct {
	config     XPicoConnectorConfig
	state      XPicoConnectorState
	port       serial.Port
	xpbridge   *xphttpbridgego.Client
	readerChan chan []byte
}

type PicoCommandBind struct {
	Command  string
	Callback func(value []byte) ([]byte, error)
}

type BridgeRefBind struct {
	Ref      string
	IsSlice  bool
	Callback func(value any)
}
