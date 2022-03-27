package config

import (
	"fmt"
	"os"
)

const (
	defaultHTTPPort   = "8081"
	defaultHost       = ""
	Prod              = "prod"
	Dev               = "dev"
	defaultPostgreSQL = "localhost"
	defaultDBPort     = "5431"
	defaultDSSLMode   = "disable"
	defaultDBuser     = "user"
	defaultDBpass     = "pass"
	defaultDBname     = "db"
	defaultSecretKey  = "secretkey"
)

type Config struct {
	Host       string
	Port       string
	ServerMode string
	DbAddress  string
	Secretkey  []byte
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
	if DbUser, ok = os.LookupEnv("DBUSER"); !ok {
		DbUser = defaultDBuser
	}

	var DbPassword string
	if DbPassword, ok = os.LookupEnv("DBPASSWORD"); !ok {
		DbPassword = defaultDBpass
	}

	var DbName string
	if DbName, ok = os.LookupEnv("DBNAME"); !ok {
		DbName = defaultDBname
	}

	var DbHost string
	if DbHost, ok = os.LookupEnv("DBHOST"); !ok {
		DbHost = defaultPostgreSQL
	}

	var DbPort string
	if DbPort, ok = os.LookupEnv("DBPORT"); !ok {
		DbPort = defaultDBPort
	}

	var DbSslMode string
	if DbSslMode, ok = os.LookupEnv("DBSSLMODE"); !ok {
		DbSslMode = defaultDSSLMode
	}

	var SecretKey string
	if SecretKey, ok = os.LookupEnv("SECRET_KEY"); !ok {
		SecretKey = (defaultSecretKey)
	}

	postgresAddr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DbHost, DbPort, DbUser, DbPassword, DbName, DbSslMode)

	return &Config{host, port, serverMode, postgresAddr, []byte(SecretKey)}
}
