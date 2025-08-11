package rabbitmq

import (
	"email_microservice/internal/config"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

type QueueManager struct {
	config config.RabbitMQConfig
	conn   *amqp.Connection
	ch     *amqp.Channel
	queue  *amqp.Queue
}

func NewQueueManager(config config.RabbitMQConfig) (*QueueManager, error) {
	var queueManager QueueManager
	queueManager = QueueManager{config: config}
	return &queueManager, nil
}

func (m *QueueManager) Connect() error {
	var err error
	if m.conn == nil {
		var virtualHost string
		if strings.HasPrefix(m.config.VirtualHost, "/") {
			virtualHost = m.config.VirtualHost[1:]
		} else {
			virtualHost = m.config.VirtualHost
		}

		rabbitMQUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
			m.config.User,
			m.config.Password,
			m.config.Host,
			m.config.Port,
			virtualHost,
		)

		var conn *amqp.Connection
		conn, err = amqp.Dial(rabbitMQUrl)
		if err != nil {
			return fmt.Errorf("failed to connect to RabbitMQ: %s", err)
		}
		m.conn = conn
	}

	if m.ch == nil {
		var ch *amqp.Channel
		ch, err = m.conn.Channel()
		if err != nil {
			return fmt.Errorf("failed to open a RabbitMQ channel: %s", err)
		}
		m.ch = ch
	}

	if m.queue == nil {
		var queue amqp.Queue
		err = m.ch.ExchangeDeclare(
			m.config.ExchangeName,
			amqp.ExchangeDirect, false, false, false, false, nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare an exchange with name %s: %s", m.config.ExchangeName, err)
		}

		queue, err = m.ch.QueueDeclare(
			m.config.QueueName,
			false, false, false, false, nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare a queue with name %s: %s", m.config.QueueName, err)
		}

		err = m.ch.QueueBind(queue.Name, m.config.RoutingKey, m.config.QueueName, false, nil)
		if err != nil {
			return fmt.Errorf("failed to bind a queue with name %s: %s", m.config.QueueName, err)
		}

		m.queue = &queue
	}
	return nil
}

func (m *QueueManager) Close() error {
	if m.conn != nil {
		err := m.conn.Close()
		if err != nil {
			return fmt.Errorf("failed to close RabbitMQ connection: %s", err)
		}
	}
	return nil
}

func (m *QueueManager) Consume() (<-chan amqp.Delivery, error) {
	err := m.Connect()
	if err != nil {
		fmt.Printf("failed to connect to RabbitMQ before comsuming: %s", err)
	}

	return m.ch.Consume(m.queue.Name, "", true, false, false, false, nil)
}
