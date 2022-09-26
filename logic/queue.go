package logic

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Consumer interface {
	Consume(message string) error
}

type Producer interface {
	Produce(message string) error
}

type RabbitConsumer struct {
	Conf Config
}

func (rc *RabbitConsumer) Consume(handle func(c Config, m []byte)) {
	var forever chan struct{}

	c, err := rc.getConnection()
	if err != nil {
		fmt.Printf("Connection couldn't be established: %v", err)
		panic(err)
	}
	defer c.Close()

	channel, err := rc.getChannel(c)
	if err != nil {
		fmt.Printf("Channel couldn't be established: %v", err)
		panic(err)
	}
	defer channel.Close()

	msgs, err := channel.Consume(
		rc.Conf.Queue.QueueName, // queue
		"",                      // consumer
		false,                   // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	go func() {
		for m := range msgs {
			handle(rc.Conf, m.Body)
			m.Ack(false)
		}
	}()

	fmt.Println("Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (rc *RabbitConsumer) getConnection() (*amqp.Connection, error) {
	var c *amqp.Connection
	var err error

	for i := 0; i < rc.Conf.Queue.RetryAmount; i++ {
		c, err = amqp.Dial(rc.Conf.Queue.Url)
		if err != nil {
			fmt.Printf("Couldn't get RMQ connection: %v", err)
			fmt.Printf("Retrying. Retries left: %v", rc.Conf.Queue.RetryAmount-i)
			fmt.Printf("Sleeping : %v seconds", i)
			time.Sleep(time.Duration(i) * time.Second)
		}
	}
	return c, err
}

func (rc *RabbitConsumer) resetConnection(c *amqp.Connection) *amqp.Connection {
	c, err := rc.getConnection()
	if err != nil {
		fmt.Printf("Automatic connection reset has failed: %v", err)
		panic(err)
	}

	return c
}

func (rc *RabbitConsumer) getChannel(c *amqp.Connection) (*amqp.Channel, error) {
	var channel *amqp.Channel
	var err error

	for i := 0; i < rc.Conf.Queue.RetryAmount; i++ {
		channel, err = c.Channel()
		if err != nil {
			fmt.Printf("Couldn't get RMQ connection's channel: %v", err)
			fmt.Printf("Retrying. Retries left: %v", rc.Conf.Queue.RetryAmount-i)
			fmt.Printf("Sleeping : %v seconds", i)
			time.Sleep(time.Duration(i) * time.Second)
		}
	}
	return channel, err
}
