package mq

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"github.com/streadway/amqp"
)

// amqpSession implements pkg.Session over an AMQP connection.
type amqpSession struct {
	service string
	conn    *amqp.Connection
}

func (s *amqpSession) Service() string  { return s.service }
func (s *amqpSession) Raw() interface{} { return s.conn }

func (s *amqpSession) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// AMQPPlugin is stateless; all connection state lives in amqpSession.
type AMQPPlugin struct{}

func (p *AMQPPlugin) Name() string { return "amqp" }

func (p *AMQPPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", task.Username, task.Password, task.IP, task.Port))
	if err != nil {
		return nil, err
	}
	return &amqpSession{service: task.Service, conn: conn}, nil
}

func (p *AMQPPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", "guest", "guest", task.IP, task.Port))
	if err != nil {
		return nil, err
	}
	return &amqpSession{service: task.Service, conn: conn}, nil
}
