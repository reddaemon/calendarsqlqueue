package queue

import (
	"fmt"

	"github.com/labstack/gommon/log"

	"github.com/reddaemon/calendarsqlqueue/config"
	"github.com/streadway/amqp"
)

func GetConnection(c *config.Config) *amqp.Connection {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		c.Broker["user"],
		c.Broker["pass"],
		c.Broker["host"],
		c.Broker["port"],
	))

	if err != nil {
		log.Fatalf("unable to connect to Rabbitmq: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("unable to get queue channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		c.Broker["exchange"],
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("unable to declare queue exchange")
	}
	return conn
}
