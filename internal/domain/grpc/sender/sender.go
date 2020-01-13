package sender

import (
	"time"

	"github.com/reddaemon/calendarsqlqueue/models/models"

	"github.com/reddaemon/calendarsqlqueue/app"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetEvent(ctx context.Context, date time.Time) error
}

type Queue interface {
	GetEvent(event *models.Event) error
}

type EventUseCase struct {
	App            *app.App
	eventRepo      Repository
	queue          Queue
	contextTimeout time.Duration
}

type Repository interface {
	GetByDate(ctx context.Context, date time.Time) ([]models.Event, error)
}

func NewEventUseCase(app *app.App, eventRepo Repository, queue Queue, contextTimeout time.Duration) *EventUseCase {
	return &EventUseCase{App: app, eventRepo: eventRepo, queue: queue, contextTimeout: contextTimeout}
}

func (e EventUseCase) GetEvent(ctx context.Context, date time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	events, err := e.eventRepo.GetByDate(ctx, date)
	if err != nil {
		return err
	}
	for _, event := range events {
		err := e.queue.GetEvent(&event)
		if err != nil {
			return err
		}
	}
	return nil
}
