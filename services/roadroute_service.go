package services

import (
	"context"
	"time"

	"shuttle/errors"
	"shuttle/databases"
	"shuttle/models"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllRoutes(SchoolID primitive.ObjectID) ([]models.RoadRoute, error) {
	client, err := databases.MongoConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("routes")

	cursor, err := collection.Find(context.Background(), bson.M{"school_id": SchoolID})
	if err != nil {
		return nil, err
	}

	var routes []models.RoadRoute
	if err = cursor.All(context.Background(), &routes); err != nil {
		return nil, err
	}

	return routes, nil
}

func GetSpecRoute(RouteID primitive.ObjectID) (models.RoadRoute, error) {
	client, err := databases.MongoConnection()
	if err != nil {
		return models.RoadRoute{}, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("routes")

	var route models.RoadRoute
	err = collection.FindOne(context.Background(), bson.M{"_id": RouteID}).Decode(&route)
	if err != nil {
		return models.RoadRoute{}, err
	}

	return route, nil
}

func AddRoute(route models.RoadRoute, SchoolID primitive.ObjectID, username string) error {
	client, err := databases.MongoConnection()
	if err != nil {
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("routes")

	var existingRoute models.RoadRoute
	err = collection.FindOne(context.Background(), bson.M{"route_name": route.RouteName}).Decode(&existingRoute)
	if err == nil {
		return errors.New("Route with similar name already exists", 409)
	}

	route.CreatedAt = time.Now()
	route.CreatedBy = username
	route.SchoolID = SchoolID

	_, err = collection.InsertOne(context.Background(), route)
	if err != nil {
		return err
	}

	return nil
}