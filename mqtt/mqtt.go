package mqtt

import (
	"fmt"

	"github.com/surgemq/message"
	"github.com/yosssi/gmq/mqtt/client"

	"github.com/think-free/axihome4/protocols/mqtt/client"
	"github.com/think-free/other/picamwrapper/config"
)

// Mqtt client
type Mqtt struct {
	conf     *config.Config
	auto     bool
	cli      *client.Client
	internal *config.Internal
}

// New Create a new mqtt client
func New(conf *config.Config, internal *config.Internal) *Mqtt {

	return &Mqtt{conf: conf, auto: false, internal: internal}
}

// Run start the mqtt client
func (m *Mqtt) Run() {

	if m.conf.DisableMqtt {
		fmt.Println("Mqtt is disabled")
		return
	}

	cli := mqttclient.NewMqttClient(m.conf.MqttClientName, m.conf.MqttServerAdd)
	cli.Connect()

	cli.SubscribeTopic(m.conf.TopicAuto, func(msg *message.PublishMessage) error {

		if string(msg.Payload()) == "true" {
			m.internal.AutoMode <- true
		} else {
			m.internal.AutoMode <- false
		}

		return nil
	})

	cli.SubscribeTopic(m.conf.TopicStop, func(msg *message.PublishMessage) error {

		if string(msg.Payload()) == m.conf.ValueStop {
			m.internal.Chwriteauto <- false
		} else if m.conf.TopicStop == m.conf.TopicStart && string(msg.Payload()) == m.conf.ValueStart {
			m.internal.Chwriteauto <- true
		}

		return nil
	})

	if m.conf.TopicStop != m.conf.TopicStart {

		cli.SubscribeTopic(m.conf.TopicStart, func(msg *message.PublishMessage) error {

			if string(msg.Payload()) == m.conf.ValueStart {
				m.internal.Chwriteauto <- true
			}

			return nil
		})
	}

	for {
		select {
		case st := <-m.internal.Chgetst:

			cli.PublishMessage(m.conf.TopicStatus, st)
		case auto := <-m.internal.Chgetauto:
			cli.PublishMessage(m.conf.TopicAutoStatus, auto)
		}
	}
}
