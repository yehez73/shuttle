package database

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var once sync.Once

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func MongoConnection() (*mongo.Client, error) {
	once.Do(func() {
		clientOptions := options.Client().ApplyURI(viper.GetString("MONGO_URI"))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		mongoClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
	})

	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	return mongoClient, nil
}