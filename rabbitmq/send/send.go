package send

import (
	"redisplay/rabbitmq/common"

	"github.com/streadway/amqp"
)

func Send(msg []byte, ch *amqp.Channel) {
	err := ch.Publish(
		"logs",
		"logsRecord", // message's routing key
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // msg persistent
			ContentType:  "text/plain",
			Body:         msg,
		},
	)
	common.Errlog(err, "publish msg to queue error")
}
