package logger

import (
	"github.com/reddaemon/calendarsqlqueue/config"
	"go.uber.org/zap"
)

func GetLogger(cfg *config.Config) (*zap.Logger, error) {
	var l *zap.Logger
	if !cfg.Debug {
		l = zap.NewNop()
		return l, nil
	}
	return l, nil
}
