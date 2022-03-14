package config

import (
	"fmt"
	"log"
	"os"
)

const (
	defaultQuHost = "localhost"
	defaultQuPort = "5672"
	defaultQuUser = "guest"
	defaultQuPass = "guest"
)

type Config struct {
	QueueAdress    string
	PostgresAdress string
}

func Init() *Config {
	var ok bool

	var QuHost string
	if QuHost, ok = os.LookupEnv("QUEUE_HOST"); !ok {
		QuHost = defaultQuHost
	}

	var QuPort string
	if QuPort, ok = os.LookupEnv("QUEUE_PORT"); !ok {
		QuPort = defaultQuPort
	}

	var QuUser string
	if QuUser, ok = os.LookupEnv("QUEUE_USER"); !ok {
		QuUser = defaultQuUser
	}

	var QuPass string
	if QuPass, ok = os.LookupEnv("QUEUE_PASS"); !ok {
		QuPass = defaultQuPass
	}

	var DbUser string
	if DbUser, ok = os.LookupEnv("mongo_user"); !ok {
		log.Fatal("no mongo user")
	}

	var DbPassword string
	if DbPassword, ok = os.LookupEnv("mongo_pass"); !ok {
		log.Fatal("no mongo pass")
	}

	var DbHost string
	if DbHost, ok = os.LookupEnv("mongo_host"); !ok {
		log.Fatal("no mongo host")
	}

	queueAdress := fmt.Sprintf("amqp://%s:%s@%s:%s/", QuUser, QuPass, QuHost, QuPort)
	mongoAddr := fmt.Sprintf("mongodb+srv://%s:%s@%s/myFirstDatabase?retryWrites=true&w=majority", DbUser, DbPassword, DbHost)

	return &Config{queueAdress, mongoAddr}
}
