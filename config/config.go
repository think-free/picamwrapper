package config

import "sync"

// Config is the configuration of the application
type Config struct {
	PiCamPath   string `usage:"The path to picam (the folder containing picam)"`
	PiCamParams string `usage:"Picam parameters"`

	HTTPListen string `usage:"Listen ip:port for the http server"`

	DisableMqtt     bool   `usage:"Disable the mqtt client, automode will not work (default to false)"`
	MqttServerAdd   string `usage:"The ip of the mqtt server"`
	MqttClientName  string `usage:"The name of the mqtt client"`
	TopicStop       string `usage:"Mqtt topic to stop recording when in auto mode"`
	ValueStop       string `usage:"Mqtt value for to stop recording"`
	TopicStart      string `usage:"Mqtt topic to start recording when in auto mode"`
	ValueStart      string `usage:"Mqtt value for to start recording"`
	TopicAuto       string `usage:"Mqtt topic to activate/desactivate auto mode"`
	TopicAutoStatus string `usage:"Mqtt topic to set the state of the auto mode"`
	TopicStatus     string `usage:"Mqtt topic to set the recording status"`
}

// Internal is a shared object for state and channels
type Internal struct {
	AutoMode    chan bool
	Chwritest   chan bool
	Chwriteauto chan bool
	Chgetst     chan bool
	Chgetauto   chan bool

	Auto  bool
	State bool
	sync.Mutex
}
