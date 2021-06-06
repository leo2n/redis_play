package controllers

import (
	"go-seckill/internal/logconf"
	"go-seckill/internal/mysql/shop/structure"
	"go-seckill/internal/rabbitmq/common"

	"github.com/sirupsen/logrus"
)

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "controllers"})

var ch = common.Ch

// 全局变量, 存储purchase_limits
var purchaseLimit = make(map[string]*structure.PurchaseLimits)