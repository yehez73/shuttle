package services

import (
	"context"
	// "errors"

	"shuttle/databases"
	"shuttle/models"

	"github.com/spf13/viper"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func AddRoadRoute(route models.RoadRoute, SchoolID primitive.ObjectID) error {
	client, err := database.MongoConnection()
	if err != nil {
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("routes")

	_, err = collection.InsertOne(context.Background(), route)
	if err != nil {
		return err
	}

	return nil
}