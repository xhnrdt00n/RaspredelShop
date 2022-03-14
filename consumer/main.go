package main

import (
	"consumer/ConsumerTypes/Category"
	"consumer/ConsumerTypes/Products"
	"consumer/config"
	"consumer/db"
	"context"
	"github.com/streadway/amqp"
	"log"
)

const (
	waitingReqMsg = " [*] Waiting for messages."
)

func main() {
	cfg := config.Init()

	conn, err := amqp.Dial(cfg.QueueAdress)
	ctx := context.Background()
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare("Shopper", "direct", false, true, false, false, nil); err != nil {
		log.Fatalln(err)
	}

	//TODO Под каждый тип создать очередь в отдельных методах и там пусть рулят событиями
	forever := make(chan bool)

	dbService, err := db.NewMongoCon(ctx, cfg)
	{
		Category.ListenAddCategory(ch, dbService)
		Category.ListenChangeCategory(ch, dbService)
		Category.ListenDeleteCategory(ch, dbService)
	}

	{
		Products.ListenAddProducts(ch, dbService)
		Products.ListenChangeProducts(ch, dbService)
		Products.ListenDeleteProducts(ch, dbService)
	}

	log.Printf(waitingReqMsg)
	<-forever
}
