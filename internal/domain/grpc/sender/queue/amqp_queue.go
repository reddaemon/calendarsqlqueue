package queue

import (
	"encoding/json"

	"github.com/labstack/gommon/log"

	"github.com/reddaemon/calendarsqlqueue/config"
	"go.uber.org/zap"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/reddaemon/calendarsqlqueue/models/models"
)

type EventQueue struct {
	Config *config.Config
	Logger *zap.Logger
	Conn   *amqp.Connection
}

func NewEventQueue(config *config.Config, logger *zap.Logger, conn *amqp.Connection) *EventQueue {
	return &EventQueue{Config: config, Logger: logger, Conn: conn}
}

func (a EventQueue) GetEvent(event *models.Event) error {
	ch, err := a.Conn.Channel()
	if err != nil {
		return err
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Fatalf("unable to close channel")
		}
	}(ch)
	body, err := json.Marshal(event)
	if err != nil {
		log.Fatalf(err.Error())
		return err
	}

	err = ch.Publish(
		a.Config.Broker["exchange"],
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		log.Fatalf("unable to publish to exchange")
		return err
	}
	log.Printf("Sent %s", body)

	return nil
}
