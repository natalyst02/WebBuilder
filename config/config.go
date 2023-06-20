package config

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

type config struct {
	databaseName string
	mongoURI     string
	storePath    string
	port         string
}

var (
	cfg  = &config{}
	once sync.Once
)

func GetConfig() (err error) {
	once.Do(func() {
		err = godotenv.Load()
		if err != nil {
			log.Error(err)
			return
		}

		cfg = &config{
			databaseName: os.Getenv("DATABASE_NAME"),
			mongoURI:     os.Getenv("MONGO_URI"),
			storePath:    os.Getenv("STORE_PATH"),
			port:         os.Getenv("PORT"),
		}
		log.Info("configurations loaded successfully.")
	})

	return nil
}

func GetDatabaseName() string {
	dbName := cfg.databaseName
	if dbName == "" {
		log.Error("DATABASE_NAME is an empty string")
	}

	return dbName
}

func GetMongoURI() string {
	uri := cfg.mongoURI
	if uri == "" {
		log.Error("MONGO_URI is an empty string")
	}

	return uri
}

func GetStorePath() string {
	path := cfg.storePath
	if path == "" {
		log.Error("STORE_PATH is an empty string")
	}

	return path
}

func GetPort() string {
	port := cfg.port
	if port == "" {
		port = "9090"
		log.Error("PORT is an empty string, default port will be 9090")
	}

	return port
}
