package services

import (
	"context"
	"errors"
	"shuttle/databases"
	"shuttle/models"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllVehicles() ([]models.Vehicle, error) {
	client, err := database.MongoConnection()
	if err != nil {
		return nil, err
	}

	var vehicles []models.Vehicle

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("vehicles")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var vehicle models.Vehicle
		if err := cursor.Decode(&vehicle); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, vehicle)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return vehicles, nil
}

func GetSpecVehicle(id string) (models.Vehicle, error) {
	var vehicle models.Vehicle
	client, err := database.MongoConnection()
	if err != nil {
		return vehicle, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("vehicles")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return vehicle, err
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&vehicle)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return vehicle, errors.New("vehicle not found")
		}
		return vehicle, err
	}

	return vehicle, nil
}

func AddVehicle(vehicle models.Vehicle) error {
	client, err := database.MongoConnection()
	if err != nil {
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("vehicles")

	_, err = collection.InsertOne(context.Background(), vehicle)
	if err != nil {
		return err
	}

	return nil
}

func UpdateVehicle(vehicle models.Vehicle, id string) error {
	client, err := database.MongoConnection()
	if err != nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("vehicles")

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": vehicle})
	if err != nil {
		return err
	}

	return nil
}

func DeleteVehicle(id string) error {
	client, err := database.MongoConnection()
	if err != nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("vehicles")

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}