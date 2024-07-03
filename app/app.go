package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/reddaemon/calendarsqlqueue/config"
)

type App struct {
	Config *config.Config
	Logger *log.Logger
	Db     *sqlx.DB
	Amqp   *amqp.Connection
}
