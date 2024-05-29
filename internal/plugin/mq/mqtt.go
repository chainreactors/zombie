package mq

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTPlugin struct {
	*pkg.Task
	client mqtt.Client
}

func (s *MQTTPlugin) Name() string {
	return s.Service
}

func (s *MQTTPlugin) Unauth() (bool, error) {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", s.IP, s.Port))
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return false, token.Error()
	}
	s.client = client
	return true, nil
}

func (s *MQTTPlugin) Login() error {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", s.IP, s.Port)).SetUsername(s.Username).SetPassword(s.Password)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	s.client = client
	return nil
}

func (s *MQTTPlugin) GetResult() *pkg.Result {
	// todo list topics
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *MQTTPlugin) Close() error {
	if s.client != nil {
		s.client.Disconnect(250)
	}
	return nil
}
