package services

import (
	"context"
	
	"log"
	"shuttle/databases"
	"shuttle/models"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllSchools() ([]models.School, error) {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var schools []models.School

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("schools")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var school models.School
		if err := cursor.Decode(&school); err != nil {
			log.Print(err)
			return nil, err
		}
		schools = append(schools, school)
	}

	if err := cursor.Err(); err != nil {
		log.Print(err)
		return nil, err
	}

	return schools, nil
}

func GetSpecSchool(id string) (models.School, error) {
	var school models.School
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return school, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("schools")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return school, err
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&school)
	if err != nil {
		return school, err
	}

	return school, nil
}

func AddSchool(school models.School) error {
	client, err := database.MongoConnection()
	if err != nil {
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("schools")

	_, err = collection.InsertOne(context.Background(), school)
	if err != nil {
		return err
	}

	return nil
}

func UpdateSchool(id string, school models.School) error {
	client, err := database.MongoConnection()
	if err != nil {
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("schools")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": school})
	if err != nil {
		return err
	}

	return nil
}

func DeleteSchool(id string) error {
	client, err := database.MongoConnection()
	if err != nil {
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("schools")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}