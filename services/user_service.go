package services

import (
	"context"
	"errors"
	"log"
	"shuttle/databases"
	"shuttle/models"
	"time"

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

func AddUser(user models.User, username string) (primitive.ObjectID, error) {
    client, err := database.MongoConnection()
    if err != nil {
        log.Print(err)
        return primitive.NilObjectID, err
    }

    collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

    // Validate common fields
    if err := validateCommonFields(user); err != nil {
        return primitive.NilObjectID, err
    }

    // Hash password if provided
    if user.Password != "" {
        hashedPassword, err := hashPassword(user.Password)
        if err != nil {
            return primitive.NilObjectID, err
        }
        user.Password = hashedPassword
    }

	user.CreatedAt = time.Now()
	user.CreatedBy = username
	user.Status = "offline"

    // Handle role-specific logic
    switch user.Role {
    case models.SuperAdmin:
        user.RoleCode = "SA"
    case models.SchoolAdmin:
        if user.Details == nil {
            return primitive.NilObjectID, errors.New("SchoolAdmin details are required")
        }

        schoolAdminDetails, ok := user.Details.(map[string]interface{})
        if !ok {
            return primitive.NilObjectID, errors.New("invalid details format for SchoolAdmin")
        }

        schoolID, ok := schoolAdminDetails["school_id"].(string)
        if !ok {
            return primitive.NilObjectID, errors.New("school_id is required for SchoolAdmin")
        }

        schoolObjectID, err := primitive.ObjectIDFromHex(schoolID)
        if err != nil {
            return primitive.NilObjectID, errors.New("invalid school_id format")
        }

        user.Details = models.SchoolAdminDetails{SchoolID: schoolObjectID}

        _, err = GetSpecSchool(schoolObjectID.Hex())
        if err != nil {
            return primitive.NilObjectID, errors.New("school not found")
        }
        user.RoleCode = "AS"
    case models.Parent:
        user.RoleCode = "P"
    case models.Driver:
		if user.Details == nil {
			return primitive.NilObjectID, errors.New("driver details are required")
		}

		driverDetails, ok := user.Details.(map[string]interface{})
		if !ok {
			return primitive.NilObjectID, errors.New("invalid details format for Driver")
		}

		vehicleID, ok := driverDetails["vehicle_id"].(string)
		if !ok {
			return primitive.NilObjectID, errors.New("vehicle_id is required for Driver")
		}

		vehicleObjectID, err := primitive.ObjectIDFromHex(vehicleID)
		if err != nil {
			return primitive.NilObjectID, errors.New("invalid vehicle_id format")
		}

		user.Details = models.DriverDetails{VehicleID: vehicleObjectID}
        user.RoleCode = "D"
    default:
        return primitive.NilObjectID, errors.New("invalid role specified")
    }

    result, err := collection.InsertOne(context.Background(), user)
    if err != nil {
        return primitive.NilObjectID, err
    }

    return result.InsertedID.(primitive.ObjectID), nil
}

func UpdateUser(id string, user models.User, username string, file []byte) error {
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
		"phone":      user.Phone,
		"address":    user.Address,
		"updated_at": time.Now(),
		"updated_by": username,
    }

    // Handle details field based on user role
    switch user.Role {
	case models.SuperAdmin:
        user.RoleCode = "SA"
    case models.SchoolAdmin:
		if user.Details == nil {
			return errors.New("SchoolAdmin details are required")
		}
	
		schoolAdminDetails, ok := user.Details.(map[string]interface{})
		if !ok {
			return errors.New("invalid details format for SchoolAdmin")
		}
	
		schoolID, ok := schoolAdminDetails["school_id"].(string)
		if !ok {
			return errors.New("school_id is required for SchoolAdmin")
		}
	
		schoolObjectID, err := primitive.ObjectIDFromHex(schoolID)
		if err != nil {
			return errors.New("invalid school_id format")
		}
	
		user.Details = models.SchoolAdminDetails{SchoolID: schoolObjectID}
	
		_, err = GetSpecSchool(schoolObjectID.Hex())
		if err != nil {
			return errors.New("school not found")
		}
	
		updateFields["details"] = models.SchoolAdminDetails{SchoolID: schoolObjectID}	
    case models.Parent:
        if user.Details == nil {
            return errors.New("parent details are required")
        }
        parentDetails, ok := user.Details.(models.ParentDetails)
        if !ok {
            return errors.New("invalid details format for Parent")
        }
        updateFields["details"] = parentDetails
    case models.Driver:
        if user.Details == nil {
            return errors.New("driver details are required")
        }
        driverDetails, ok := user.Details.(models.DriverDetails)
        if !ok {
            return errors.New("invalid details format for Driver")
        }

        // Validate the vehicle ID in DriverDetails
        vehicleID := driverDetails.VehicleID
        if vehicleID.IsZero() {
            return errors.New("invalid vehicle_id format for Driver")
        }

        // Update the 'details' field directly
        updateFields["details"] = driverDetails
    default:
        return errors.New("invalid role specified")
    }

    // Perform the update in the database
    _, err = collection.UpdateOne(
        context.Background(),
        bson.M{"_id": objectID},
        bson.M{"$set": updateFields},
    )

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

func UpdateUserStatus(userID primitive.ObjectID, status string, lastActive time.Time) error {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	update := bson.M{
		"$set": bson.M{
			"status":      status,
			"last_active": lastActive,
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": userID}, update)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}