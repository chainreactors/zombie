package mq

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"github.com/streadway/amqp"
)

type AMQPPlugin struct {
	*pkg.Task
	conn *amqp.Connection
}

func (s *AMQPPlugin) Name() string {
	return s.Service
}

func (s *AMQPPlugin) Unauth() (bool, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", "guest", "guest", s.IP, s.Port))
	if err != nil {
		return false, err
	}
	s.conn = conn
	return true, nil
}

func (s *AMQPPlugin) Login() error {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", s.Username, s.Password, s.IP, s.Port))
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *AMQPPlugin) GetResult() *pkg.Result {
	// todo list queues
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *AMQPPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
