package Category

import (
	"consumer/db"
	"consumer/model"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func ListenAddCategory(ch *amqp.Channel, service db.DBService) {
	q, err := ch.QueueDeclare(
		"addCategoryQu", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalln(err)
	}

	if err = ch.QueueBind(
		"addCategoryQu", // name of the queue
		"addCategory",   // bindingKey
		"Shopper",       // sourceExchange
		false,           // noWait
		nil,             // arguments
	); err != nil {
		log.Fatalf("Queue Bind: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for d := range msgs {
			fmt.Println(string(d.Body))
			var category model.Category
			err := json.Unmarshal(d.Body, &category)
			if err != nil {
				log.Println(err)
				continue
			}

			if err := service.AddCategories(category); err != nil {
				log.Println(err)
				continue
			}

			fmt.Printf("Категория успешно добавлено из очереди")
		}
	}()
}
