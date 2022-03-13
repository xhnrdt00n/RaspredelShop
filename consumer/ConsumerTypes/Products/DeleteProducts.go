package Products

import (
	"consumer/db"
	"github.com/streadway/amqp"
)

func ListenDeleteProducts(ch *amqp.Channel, service db.DBService) {
	//TODO описать метод
}
