package xpicoconnect

import (
	xphttpbridgego "github.com/steveiliop56/xphttpbridge-go"
	"go.bug.st/serial"
)

type SerialConfig struct {
	Baudrate   int    `ini:"baudrate"`
	Port       string `ini:"port"`
	BufferSize int    `ini:"buffer_size"`
	Timeout    int    `ini:"timeout"`
}

type XPHTTPBridgeConfig struct {
	Address string `ini:"address"`
	Port    int    `ini:"port"`
}

type InitConfig struct {
	RetryInterval int `ini:"retry_interval"`
	MaxRetries    int `ini:"max_retries"`
}

type XPicoConnectorConfig struct {
	SerialConfig       SerialConfig       `ini:"serial"`
	XPHTTPBridgeConfig XPHTTPBridgeConfig `ini:"xphttpbridge"`
	InitConfig         InitConfig         `ini:"init"`
	PollTime           int                `ini:"poll_time"`
}

type XPicoConnectorState struct {
	isCommandPending   bool
	picoCommandBinders map[string]PicoCommandBind
	bridgeRefBinders   map[string]BridgeRefBind
}

type XPicoConnector struct {
	config     XPicoConnectorConfig
	state      XPicoConnectorState
	port       *serial.Port
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
