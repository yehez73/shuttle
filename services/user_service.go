package services

import (
	"context"

	"net/mail"
	"shuttle/databases"
	"shuttle/errors"
	"shuttle/models"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetSpecUser(id string) (models.User, error) {
	client, err := database.MongoConnection()
	if err != nil {
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
		return models.User{}, err
	}

	return user, nil
}

func GetAllSuperAdmin() ([]models.UserResponse, error) {
	client, err := database.MongoConnection()
	if err != nil {
		return nil, err
	}

	var users []models.UserResponse

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{"role": models.SuperAdmin})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.UserResponse
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetSpecSuperAdmin(id string) (models.UserResponse, error) {
	client, err := database.MongoConnection()
	if err != nil {
		return models.UserResponse{}, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	var user models.UserResponse
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return models.UserResponse{}, err
	}

	return user, nil
}

func GetAllSchoolAdmin() ([]models.UserResponse, error) {
	client, err := database.MongoConnection()
	if err != nil {
		return nil, err
	}

	var users []models.UserResponse

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "role", Value: models.SchoolAdmin},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "schools"},
				{Key: "localField", Value: "details.school_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "school_details"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$school_details"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			{Key: "$project", Value: bson.D{
				{Key: "id", Value: "$_id"},
				{Key: "picture", Value: 1},
				{Key: "first_name", Value: 1},
				{Key: "last_name", Value: 1},
				{Key: "gender", Value: 1},
				{Key: "email", Value: 1},
				{Key: "role", Value: 1},
				{Key: "role_code", Value: 1},
				{Key: "phone", Value: 1},
				{Key: "address", Value: 1},
				{Key: "status", Value: 1},
				{Key: "last_active", Value: 1},
				{Key: "details", Value: bson.D{
					{Key: "school_name", Value: bson.M{"$ifNull": []interface{}{"$school_details.name", ""}}},
				}},
				{Key: "created_at", Value: 1},
				{Key: "created_by", Value: 1},
				{Key: "updated_at", Value: 1},
				{Key: "updated_by", Value: 1},
			}},
		},
	}	

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.UserResponse
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		if details, ok := user.Details.(primitive.D); ok {
			detailsMap := make(map[string]interface{})
			for _, elem := range details {
				detailsMap[elem.Key] = elem.Value
			}
			user.Details = detailsMap
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetSpecSchoolAdmin(id string) (models.UserResponse, error) {
	client, err := database.MongoConnection()
	if err != nil {
		return models.UserResponse{}, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	var user models.UserResponse
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: objectID},
				{Key: "role", Value: models.SchoolAdmin},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "schools"},
				{Key: "localField", Value: "details.school_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "school_details"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$school_details"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			{Key: "$project", Value: bson.D{
				{Key: "id", Value: "$_id"},
				{Key: "picture", Value: 1},
				{Key: "first_name", Value: 1},
				{Key: "last_name", Value: 1},
				{Key: "gender", Value: 1},
				{Key: "email", Value: 1},
				{Key: "role", Value: 1},
				{Key: "role_code", Value: 1},
				{Key: "phone", Value: 1},
				{Key: "address", Value: 1},
				{Key: "status", Value: 1},
				{Key: "last_active", Value: 1},
				{Key: "details", Value: bson.D{
					{Key: "school_id", Value: "$details.school_id"},
					{Key: "school_name", Value: bson.M{"$ifNull": []interface{}{"$school_details.name", ""}}},
				}},
				{Key: "created_at", Value: 1},
				{Key: "created_by", Value: 1},
				{Key: "updated_at", Value: 1},
				{Key: "updated_by", Value: 1},
			}},
		},
	}
	
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return models.UserResponse{}, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		if err := cursor.Decode(&user); err != nil {
			return models.UserResponse{}, err
		}
		if details, ok := user.Details.(primitive.D); ok {
			detailsMap := make(map[string]interface{})
			for _, elem := range details {
				detailsMap[elem.Key] = elem.Value
			}
			user.Details = detailsMap
		}
	}

	if err := cursor.Err(); err != nil {
		return models.UserResponse{}, err
	}

	return user, nil
}

func GetAllDriverFromAllSchools() ([]models.UserResponse, error) {
    client, err := database.MongoConnection()
    if err != nil {
        return nil, err
    }

    var users []models.UserResponse
    collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

    pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "role", Value: models.Driver},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "schools"},
				{Key: "localField", Value: "details.school_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "school_details"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$school_details"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "vehicles"},
				{Key: "localField", Value: "details.vehicle_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "vehicle_details"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$vehicle_details"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			{Key: "$project", Value: bson.D{
				{Key: "id", Value: "$_id"},
				{Key: "picture", Value: 1},
				{Key: "first_name", Value: 1},
				{Key: "last_name", Value: 1},
				{Key: "gender", Value: 1},
				{Key: "email", Value: 1},
				{Key: "role", Value: 1},
				{Key: "role_code", Value: 1},
				{Key: "phone", Value: 1},
				{Key: "address", Value: 1},
				{Key: "status", Value: 1},
				{Key: "details", Value: bson.D{
					{Key: "license_number", Value: "$details.license_number"},
					{Key: "school_name", Value: bson.M{"$ifNull": []interface{}{"$school_details.name", ""}}},
					{Key: "vehicle_name", Value: bson.M{"$ifNull": []interface{}{"$vehicle_details.name", ""}}},
				}},
				{Key: "created_at", Value: 1},
				{Key: "created_by", Value: 1},
				{Key: "updated_at", Value: 1},
				{Key: "updated_by", Value: 1},
			}},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.UserResponse
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		if details, ok := user.Details.(primitive.D); ok {
			detailsMap := make(map[string]interface{})
			for _, elem := range details {
				detailsMap[elem.Key] = elem.Value
			}
			user.Details = detailsMap
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetAllDriverForPermittedSchool(schoolID primitive.ObjectID) ([]models.UserResponse, error) {
	client, err := database.MongoConnection()
	if err != nil {
		return nil, err
	}

	var users []models.UserResponse
	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "role", Value: models.Driver},
				{Key: "details.school_id", Value: schoolID},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "schools"},
				{Key: "localField", Value: "details.school_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "school_details"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$school_details"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "vehicles"},
				{Key: "localField", Value: "details.vehicle_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "vehicle_details"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$vehicle_details"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			{Key: "$project", Value: bson.D{
				{Key: "id", Value: "$_id"},
				{Key: "picture", Value: 1},
				{Key: "first_name", Value: 1},
				{Key: "last_name", Value: 1},
				{Key: "gender", Value: 1},
				{Key: "email", Value: 1},
				{Key: "role", Value: 1},
				{Key: "role_code", Value: 1},
				{Key: "phone", Value: 1},
				{Key: "address", Value: 1},
				{Key: "status", Value: 1},
				{Key: "details", Value: bson.D{
					{Key: "license_number", Value: "$details.license_number"},
					{Key: "school_name", Value: bson.M{"$ifNull": []interface{}{"$school_details.name", ""}}},
					{Key: "vehicle_name", Value: bson.M{"$ifNull": []interface{}{"$vehicle_details.name", ""}}},
				}},
				{Key: "created_at", Value: 1},
				{Key: "created_by", Value: 1},
				{Key: "updated_at", Value: 1},
				{Key: "updated_by", Value: 1},
			}},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.UserResponse
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		if details, ok := user.Details.(primitive.D); ok {
			detailsMap := make(map[string]interface{})
			for _, elem := range details {
				detailsMap[elem.Key] = elem.Value
			}
			user.Details = detailsMap
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetSpecDriverFromAllSchools(id string) (models.UserResponse, error) {
	client, err := database.MongoConnection()
	if err != nil {
		return models.UserResponse{}, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	var user models.UserResponse
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: objectID},
				{Key: "role", Value: models.Driver},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "schools"},
				{Key: "localField", Value: "details.school_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "school_details"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$school_details"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "vehicles"},
				{Key: "localField", Value: "details.vehicle_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "vehicle_details"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$vehicle_details"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			{Key: "$project", Value: bson.D{
				{Key: "id", Value: "$_id"},
				{Key: "picture", Value: 1},
				{Key: "first_name", Value: 1},
				{Key: "last_name", Value: 1},
				{Key: "gender", Value: 1},
				{Key: "email", Value: 1},
				{Key: "role", Value: 1},
				{Key: "role_code", Value: 1},
				{Key: "phone", Value: 1},
				{Key: "address", Value: 1},
				{Key: "status", Value: 1},
				{Key: "details", Value: bson.D{
					{Key: "license_number", Value: "$details.license_number"},
					{Key: "school_id", Value: "$details.school_id"},
					{Key: "school_name", Value: bson.M{"$ifNull": []interface{}{"$school_details.name", ""}}},
					{Key: "vehicle_id", Value: "$details.vehicle_id"},
					{Key: "vehicle_name", Value: bson.M{"$ifNull": []interface{}{"$vehicle_details.name", ""}}},
				}},
				{Key: "created_at", Value: 1},
				{Key: "created_by", Value: 1},
				{Key: "updated_at", Value: 1},
				{Key: "updated_by", Value: 1},
			}},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return models.UserResponse{}, err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&user); err != nil {
			return models.UserResponse{}, err
		}

		if details, ok := user.Details.(bson.D); ok {
			detailsMap := make(map[string]interface{})
			for _, elem := range details {
				detailsMap[elem.Key] = elem.Value
			}
			user.Details = detailsMap
		}
	} else {
		return models.UserResponse{}, err
	}

	if err := cursor.Err(); err != nil {
		return models.UserResponse{}, err
	}

	return user, nil
}

func GetSpecDriverForPermittedSchool(id string, schoolID primitive.ObjectID) (models.UserResponse, error) {
	client, err := database.MongoConnection()
	if err != nil {
		return models.UserResponse{}, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	var user models.UserResponse
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: objectID},
				{Key: "role", Value: models.Driver},
				{Key: "details.school_id", Value: schoolID},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "schools"},
				{Key: "localField", Value: "details.school_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "school_details"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$school_details"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "vehicles"},
				{Key: "localField", Value: "details.vehicle_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "vehicle_details"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$vehicle_details"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		{
			{Key: "$project", Value: bson.D{
				{Key: "id", Value: "$_id"},
				{Key: "picture", Value: 1},
				{Key: "first_name", Value: 1},
				{Key: "last_name", Value: 1},
				{Key: "gender", Value: 1},
				{Key: "email", Value: 1},
				{Key: "role", Value: 1},
				{Key: "role_code", Value: 1},
				{Key: "phone", Value: 1},
				{Key: "address", Value: 1},
				{Key: "status", Value: 1},
				{Key: "details", Value: bson.D{
					{Key: "license_number", Value: "$details.license_number"},
					{Key: "school_name", Value: bson.M{"$ifNull": []interface{}{"$school_details.name", ""}}},
					{Key: "vehicle_id", Value: "$details.vehicle_id"},
					{Key: "vehicle_name", Value: bson.M{"$ifNull": []interface{}{"$vehicle_details.name", ""}}},
				}},
				{Key: "created_at", Value: 1},
				{Key: "created_by", Value: 1},
				{Key: "updated_at", Value: 1},
				{Key: "updated_by", Value: 1},
			}},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return models.UserResponse{}, err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&user); err != nil {
			return models.UserResponse{}, err
		}

		if details, ok := user.Details.(bson.D); ok {
			detailsMap := make(map[string]interface{})
			for _, elem := range details {
				detailsMap[elem.Key] = elem.Value
			}
			user.Details = detailsMap
		}
	} else {
		return models.UserResponse{}, err
	}

	if err := cursor.Err(); err != nil {
		return models.UserResponse{}, err
	}

	return user, nil
}

func AddUser(user models.User, username string) (primitive.ObjectID, error) {
    client, err := database.MongoConnection()
    if err != nil {
        return primitive.NilObjectID, err
    }

    collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

    if err := validateCommonFields(user); err != nil {
        return primitive.NilObjectID, err
    }

	if len(user.Password) < 8 {
		return primitive.NilObjectID, errors.New("password must be at least 8 characters", 400)
	}

    if user.Password != "" {
        hashedPassword, err := hashPassword(user.Password)
        if err != nil {
            return primitive.NilObjectID, err
        }
        user.Password = hashedPassword
    }

    var existingUser models.User
    err = collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
    if err == nil {
        return primitive.NilObjectID, errors.New("email already exists", 409)
    }

    user.CreatedAt = time.Now()
    user.CreatedBy = username
    user.Status = "offline"

    if err := processRoleDetails(&user); err != nil {
        return primitive.NilObjectID, err
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
        return err
    }

    collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    if err := validateCommonFields(user); err != nil {
        return err
    }

    updateFields := bson.M{
		"picture":    user.Picture,
        "first_name": user.FirstName,
        "last_name":  user.LastName,
		"gender":     user.Gender,
        "email":      user.Email,
        "role":       user.Role,
		"phone":      user.Phone,
		"address":    user.Address,
		"updated_at": time.Now(),
		"updated_by": username,
    }

    if err := processRoleDetails(&user); err != nil {
        return err
    }

    updateFields["details"] = user.Details

    _, err = collection.UpdateOne(
        context.Background(),
        bson.M{"_id": objectID},
        bson.M{"$set": updateFields},
    )

    return err
}

func processRoleDetails(user *models.User) error {
    switch user.Role {
	case models.SuperAdmin:
		user.RoleCode = "SA"
    case models.SchoolAdmin:
        if user.Details == nil {
            return errors.New("SchoolAdmin details are required", 400)
        }

        schoolAdminDetails, ok := user.Details.(map[string]interface{})
        if !ok {
            return errors.New("invalid details format for SchoolAdmin", 400)
        }

        schoolID, ok := schoolAdminDetails["school_id"].(string)
        if !ok {
            return errors.New("school_id is required for SchoolAdmin", 400)
        }

        schoolObjectID, err := primitive.ObjectIDFromHex(schoolID)
        if err != nil {
            return errors.New("invalid school_id format", 400)
        }

        user.Details = models.SchoolAdminDetails{SchoolID: schoolObjectID}
        _, err = GetSpecSchool(schoolObjectID.Hex())
        if err != nil {
            return errors.New("school not found", 400)
        }
		user.RoleCode = "AS"
    case models.Parent:
        if user.Details == nil {
            return errors.New("parent details are required", 400)
        }
		user.RoleCode = "P"
	case models.Driver:
		if user.Details == nil {
			return errors.New("driver details are required", 400)
		}
	
		driverDetails, ok := user.Details.(map[string]interface{})
		if !ok {
			return errors.New("invalid details format for Driver", 400)
		}
	
		licenseNumber, ok := driverDetails["license_number"].(string)
		if !ok || len(licenseNumber) == 0 {
			return errors.New("license number is required for Driver", 400)
		}
	
		vehicleID, ok := driverDetails["vehicle_id"].(string)
		if !ok {
			vehicleID = ""
		}
	
		if vehicleID != "" {
			vehicleObjectID, err := primitive.ObjectIDFromHex(vehicleID)
			if err != nil {
				return errors.New("invalid vehicle_id format", 400)
			}
			_, err = GetSpecVehicle(vehicleObjectID.Hex())
			if err != nil {
				return errors.New("vehicle not found", 400)
			}
			driverDetails["vehicle_id"] = vehicleObjectID
		} else {
			driverDetails["vehicle_id"] = nil
		}
	
		schoolID, ok := driverDetails["school_id"].(string)
		if !ok {
			schoolID = ""
		}
	
		if schoolID != "" {
			schoolObjectID, err := primitive.ObjectIDFromHex(schoolID)
			if err != nil {
				return errors.New("invalid school_id format", 400)
			}
			_, err = GetSpecSchool(schoolObjectID.Hex())
			if err != nil {
				return errors.New("school not found", 400)
			}
			driverDetails["school_id"] = schoolObjectID
		} else {
			driverDetails["school_id"] = nil
		}
	
		var vehicleObjectID, schoolObjectID primitive.ObjectID
		if driverDetails["vehicle_id"] != nil {
			vehicleObjectID = driverDetails["vehicle_id"].(primitive.ObjectID)
		}
		if driverDetails["school_id"] != nil {
			schoolObjectID = driverDetails["school_id"].(primitive.ObjectID)
		}
	
		user.Details = models.DriverDetails{
			LicenseNumber: licenseNumber,
			SchoolID:      schoolObjectID,
			VehicleID:     vehicleObjectID,
		}
		user.RoleCode = "D"	
    default:
        return errors.New("invalid role specified", 400)
    }
    return nil
}

func DeleteUser(id string) error {
	client, err := database.MongoConnection()
	if err != nil {
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
		return err
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}

func DeleteSchoolDriver(id string, schoolID primitive.ObjectID) error {
	client, err := database.MongoConnection()
	if err != nil {
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
		return err
	}

	if details, ok := user.Details.(map[string]interface{}); ok {
		if details["school_id"].(primitive.ObjectID).Hex() != schoolID.Hex() {
			return errors.New("driver does not belong to this school", 400)
		}
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserStatus(userID primitive.ObjectID, status string, lastActive time.Time) error {
	client, err := database.MongoConnection()
	if err != nil {
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
		return err
	}

	return nil
}

func validateCommonFields(user models.User) error {
	validRoles := map[models.Role]bool{
		models.SuperAdmin:  true,
		models.SchoolAdmin: true,
		models.Parent:      true,
		models.Driver:      true,
	}
	if !validRoles[models.Role(strings.ToLower(string(user.Role)))] {
		return errors.New("invalid role", 400)
	}

	validGender := map[models.Gender]bool{
		models.Female: true,
		models.Male:   true,
	}
	if !validGender[models.Gender(strings.ToLower(string(user.Gender)))] {
		return errors.New("invalid gender", 400)
	}

	if len(user.Phone) < 12 || len(user.Phone) > 15 {
		return errors.New("phone number must be between 12 and 15 characters", 400)
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return errors.New("invalid email format", 400)
	}
	return nil
}