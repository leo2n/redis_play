package receiver

import (
	"go-seckill/internal/logconf"

	"github.com/sirupsen/logrus"
)

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "rabbitmq-receiver"})
