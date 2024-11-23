package services

import (
	"context"
	"errors"
	"log"
	"shuttle/databases"
	"shuttle/models"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllUser() ([]models.User, error) {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var users []models.User

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Print(err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Print(err)
		return nil, err
	}

	return users, nil
}

func GetSpecUser(id string) (models.User, error) {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return models.User{}, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	var user models.User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		log.Print(err)
		return models.User{}, err
	}

	return user, nil
}

func AddUser(user models.User) (primitive.ObjectID, error) {
    client, err := database.MongoConnection()
    if err != nil {
        log.Print(err)
        return primitive.NilObjectID, err
    }

    collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

    if err := validateCommonFields(user); err != nil {
        return primitive.NilObjectID, err
    }

    if user.Password != "" {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			return primitive.NilObjectID, err
		}
		user.Password = hashedPassword
	}	

	// Different roles require different details
    switch user.Role {
    case models.SuperAdmin:
        if user.SuperAdminDetails == nil {
            return primitive.NilObjectID, errors.New("super admin details are required for super admin role")
        }
		user.RoleCode = "SA"
    case models.SchoolAdmin:
        if user.SchoolAdminDetails == nil {
            return primitive.NilObjectID, errors.New("school admin details are required for school admin role")
        }
        SchoolID := user.SchoolAdminDetails.SchoolID
        _, err := GetSpecSchool(SchoolID.Hex())
        if err != nil {
            return primitive.NilObjectID, errors.New("school not found")
        }
		user.RoleCode = "AS"
    case models.Parent:
        if user.ParentDetails == nil {
            return primitive.NilObjectID, errors.New("parent details are required for parent role")
        }
		user.RoleCode = "P"
    case models.Driver:
        if user.DriverDetails == nil {
            return primitive.NilObjectID, errors.New("driver details are required for driver role")
        }
		user.RoleCode = "D"
    }

    result, err := collection.InsertOne(context.Background(), user)
    if err != nil {
        return primitive.NilObjectID, err
    }

    return result.InsertedID.(primitive.ObjectID), nil
}

func UpdateUser(id string, user models.User) error {
    client, err := database.MongoConnection()
    if err != nil {
        log.Print(err)
        return err
    }

    collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    var existingUser models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&existingUser)
	if err != nil {
		return errors.New("user not found")
	}

    if err := validateCommonFields(user); err != nil {
        return err
    }

    if user.Password != "" {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

    updateFields := bson.M{
        "first_name": user.FirstName,
        "last_name":  user.LastName,
        "email":      user.Email,
        "password":   user.Password,
        "role":       user.Role,
    }

	// Different roles require different details
    switch user.Role {
    case models.SuperAdmin:
        if user.SuperAdminDetails == nil {
            return errors.New("super admin details are required for super admin role")
        }
        updateFields["super_admin_details"] = user.SuperAdminDetails
    case models.SchoolAdmin:
        if user.SchoolAdminDetails == nil {
            return errors.New("school admin details are required for school admin role")
        }
		SchoolID := user.SchoolAdminDetails.SchoolID
		_, err := GetSpecSchool(SchoolID.Hex())
		if err != nil {
			return errors.New("school not found")
		}
        updateFields["school_admin_details"] = user.SchoolAdminDetails
    case models.Parent:
        if user.ParentDetails == nil {
            return errors.New("parent details are required for parent role")
        }
        updateFields["parent_details"] = user.ParentDetails
    case models.Driver:
        if user.DriverDetails == nil {
            return errors.New("driver details are required for driver role")
        }
        updateFields["driver_details"] = user.DriverDetails
    }
    _, err = collection.UpdateOne(context.Background(),bson.M{"_id": objectID},bson.M{"$set": updateFields})
	return err
}

func validateCommonFields(user models.User) error {
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	validRoles := map[models.Role]bool{
		models.SuperAdmin:  true,
		models.SchoolAdmin: true,
		models.Parent:      true,
		models.Driver:      true,
	}
	if !validRoles[user.Role] {
		return errors.New("invalid role")
	}
	return nil
}

func DeleteUser(id string) error {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func UpdateUserStatus(userID primitive.ObjectID, status string) error {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}