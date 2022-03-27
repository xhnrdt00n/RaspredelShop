package config

import (
	"fmt"
	"log"
	"os"
)

const (
	defaultHTTPPort  = "8080"
	defaultHost      = ""
	Prod             = "prod"
	Dev              = "dev"
	defaultSecretKey = "secretkey"
)

type Config struct {
	Host       string
	Port       string
	ServerMode string
	DbAddress  string
	Secretkey  []byte
}

type AuthType struct {
	VKconfig      AuthConfig
	DiscordConfig AuthConfig
	GoogleConfig  AuthConfig
}

type AuthConfig struct {
	ClientID     string
	ClientSecret string
}

func Init() *Config {
	serverMode, ok := os.LookupEnv("SERVER_MODE")
	if !ok {
		serverMode = Dev
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultHTTPPort
	}

	var DbUser string
	if DbUser, ok = os.LookupEnv("mongo_user"); !ok {
		log.Fatal("db user not set up in env")
	}

	var DbPassword string
	if DbPassword, ok = os.LookupEnv("mongo_pass"); !ok {
		log.Fatal("db pass not set up in env")
	}

	var DbHost string
	if DbHost, ok = os.LookupEnv("mongo_host"); !ok {
		log.Fatal("db host not set up in env")
	}

	var SecretKey string
	if SecretKey, ok = os.LookupEnv("SECRET_KEY"); !ok {
		SecretKey = (defaultSecretKey)
	}

	mongoAddr := fmt.Sprintf("mongodb+srv://%s:%s@%s/myFirstDatabase?retryWrites=true&w=majority", DbUser, DbPassword, DbHost)

	return &Config{host, port, serverMode, mongoAddr, []byte(SecretKey)}
}
