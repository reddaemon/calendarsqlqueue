package queue

import (
	"encoding/json"

	"fmt"

	"github.com/reddaemon/calendarsqlqueue/app"
	"github.com/reddaemon/calendarsqlqueue/models/models"
	"github.com/streadway/amqp"
)

type EventQueue struct {
	app *app.App
}

func NewEventQueue(app *app.App) *EventQueue {
	return &EventQueue{app: app}
}

func (a EventQueue) GetEvent(event *models.Event) error {
	ch, err := a.app.Amqp.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = ch.Publish(
		a.app.Config.Broker.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plaint",
			Body:        body,
		})
	if err != nil {
		return err
	}
	a.app.Logger.Info(fmt.Sprintf("Sent %s", body))

	return nil
}
