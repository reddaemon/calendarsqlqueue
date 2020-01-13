package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/reddaemon/calendarsqlqueue/config"
	"github.com/streadway/amqp"
)

type App struct {
	Config *config.Config
	Logger *log.Logger
	Db     *sqlx.DB
	Amqp   *amqp.Connection
}
