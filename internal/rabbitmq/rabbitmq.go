package rabbitmq

import (
	"os"

	"github.com/streadway/amqp"
)

// AMQP: Advacned Message Queuing Protocol
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {
	// Dial establishes a connection to the RabbitMQ server.
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, err
	}

	// Create a new channel for communication with RabbitMQ.
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{Conn: conn, Channel: ch}, nil
}

// Publish sends a message to a specified RabbitMQ queue.
// It declares the queue if it doesn't already exist.
func (r *RabbitMQ) Publish(queueName string, message []byte) error {
	q, err := r.Channel.QueueDeclare(
		queueName, // name of the queue
		false,     // durable (whether the queue should survive a broker restart)
		false,     // delete when unused
		false,     // exclusive (used by only one connection and the queue will be deleted when that connection closes)
		false,     // no-wait (the server will not respond to the method)
		nil,       // arguments (optional arguments for the queue)
	)
	if err != nil {
		return err
	}

	return r.Channel.Publish(
		"",     // exchange (default exchange)
		q.Name, // routing key (queue name)
		false,  // mandatory (if true, the message must be routed to a queue or return an error)
		false,  // immediate (if true, the message must be delivered to a consumer immediately or return an error)
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
}

// Consume listens to a specified RabbitMQ queue.
func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	q, err := r.Channel.QueueDeclare(
		queueName, // name of the queue
		false,     // durable (whether the queue should survive a broker restart)
		false,     // delete when unused
		false,     // exclusive (used by only one connection and the queue will be deleted when that connection closes)
		false,     // no-wait (the server will not respond to the method)
		nil,       // arguments (optional arguments for the queue)
	)
	if err != nil {
		return nil, err
	}

	msgs, err := r.Channel.Consume(
		q.Name, // queue name
		"",     // consumer tag (empty string means auto-generated)
		true,   // auto-ack (automatic acknowledgment)
		false,  // exclusive (used by only this consumer)
		false,  // no-local (not supported by RabbitMQ)
		false,  // no-wait (the server will not respond to the method)
		nil,    // arguments (optional arguments for consuming)
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.Conn.Close()
}
