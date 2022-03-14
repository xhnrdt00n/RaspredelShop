package db

import (
	"StoreServer/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const waitTimeout = 10 * time.Second

type MongoCon struct {
	mongoConn    *mongo.Client
	mongoConnCtx context.Context
}

func NewMongoCon(ctx context.Context, config *config.Config) (*MongoCon, error) {
	// Контекст ограниченный по времени ожидания
	instance := &MongoCon{}

	instance.mongoConnCtx = ctx

	var err error
	var count = 0
	for {
		if count < 6 {
			count++
		}
		err = instance.reconnect(config.DbAddress)
		if err != nil {
			log.Printf("connection was lost. Error: %s. Wait %d sec.", err, count*5)
		} else {
			break
		}
		log.Println("Try to reconnect...")
		time.Sleep(time.Duration(count*5) * time.Second)
	}
	return instance, nil
}

func (db *MongoCon) Close() error {
	db.mongoConn.Disconnect(db.mongoConnCtx)
	return nil
}

func (db *MongoCon) reconnect(address string) error {
	connCtx, cancel := context.WithTimeout(db.mongoConnCtx, waitTimeout)
	defer cancel()

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	//TODO заменить на постановку адреса
	clientOptions := options.Client().
		ApplyURI(address).
		SetServerAPIOptions(serverAPIOptions)
	defer cancel()
	client, err := mongo.Connect(connCtx, clientOptions)
	if err != nil {
		return fmt.Errorf("unable to connection to database: %v", err)
	}

	db.mongoConn = client
	return nil
}
