package mq

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// mqttSession implements pkg.Session over an MQTT client connection.
type mqttSession struct {
	service string
	client  mqtt.Client
}

func (s *mqttSession) Service() string  { return s.service }
func (s *mqttSession) Raw() interface{} { return s.client }

func (s *mqttSession) Close() error {
	if s.client != nil {
		s.client.Disconnect(250)
	}
	return nil
}

// MQTTPlugin is stateless; all connection state lives in mqttSession.
type MQTTPlugin struct{}

func (p *MQTTPlugin) Name() string { return "mqtt" }

func (p *MQTTPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", task.IP, task.Port)).SetUsername(task.Username).SetPassword(task.Password)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &mqttSession{service: task.Service, client: client}, nil
}

func (p *MQTTPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", task.IP, task.Port))
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &mqttSession{service: task.Service, client: client}, nil
}
