package main

import (
	"flag"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"

	"github.com/reddaemon/calendarsqlqueue/internal/database/postgres"

	"context"
	"github.com/reddaemon/calendarsqlqueue/config"
	dbc "github.com/reddaemon/calendarsqlqueue/db"
	"github.com/reddaemon/calendarsqlqueue/logger"
	qc "github.com/reddaemon/calendarsqlqueue/queue"

	"github.com/reddaemon/calendarsqlqueue/internal/domain/grpc/sender"
	"github.com/reddaemon/calendarsqlqueue/internal/domain/grpc/sender/queue"
)

func main() {
	SchedConfigPath := flag.String("config", "config.yml", "path to config file")
	flag.Parse()
	c, err := config.GetConfig(*SchedConfigPath)
	if err != nil {
		log.Fatal("unable to get config")
	}
	l, err := logger.GetLogger(c)
	if err != nil {
		log.Fatalf("unable to get logger")
	}
	db, err := dbc.GetDb(c)
	if err != nil {
		log.Fatalf("unable to get db")
	}
	amqp := qc.GetConnection(c)
	defer amqp.Close()
	contextTimeout := time.Millisecond * 500
	q := queue.NewEventQueue(c, l, amqp)
	repo := postgres.NewPsqlRepository(db)
	ucs := sender.NewEventUseCase(repo, q, contextTimeout)
	ctx := context.Background()

	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	diff := end.Sub(now)
	if diff < 0 {
		end = end.Add(24 * time.Hour)
		diff = end.Sub(now)
	}

	for {
		err := ucs.GetEvent(ctx, time.Now())
		if err != nil {
			l.Error(err.Error())
		}
		time.Sleep(diff)
	}
}
