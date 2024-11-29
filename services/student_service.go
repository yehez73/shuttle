package services

import (
	"context"
	"log"

	"shuttle/errors"
	"shuttle/databases"
	"shuttle/models"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllPermitedSchoolStudentsWithParents(schoolID primitive.ObjectID) ([]models.SchoolStudentParentResponse, error) {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("students")

	var students []models.Student
	cursor, err := collection.Find(context.Background(), bson.M{"school_id": schoolID})
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var student models.Student
		if err := cursor.Decode(&student); err != nil {
			log.Print(err)
			return nil, err
		}
		students = append(students, student)
	}

	if err := cursor.Err(); err != nil {
		log.Print(err)
		return nil, err
	}

	var Parents []models.ParentResponse
	parentCollection := client.Database(viper.GetString("MONGO_DB")).Collection("users")
	for _, student := range students {
		var parent models.ParentResponse
		err := parentCollection.FindOne(context.Background(), bson.M{"_id": student.ParentID}, options.FindOne().SetProjection(bson.M{"password": 0})).Decode(&parent)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		Parents = append(Parents, parent)
	}

	var schoolStudents []models.SchoolStudentParentResponse
	for i, student := range students {
		schoolStudents = append(schoolStudents, models.SchoolStudentParentResponse{
			Student: student,
			Parent:  Parents[i],
		})
	}

	return schoolStudents, nil
}

func AddPermittedSchoolStudentWithParents(student models.SchoolStudentRequest, schoolID primitive.ObjectID, username string) error {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return err
	}

	if (models.User{}) == student.Parent {
		return errors.New("parent data is required", 400)
	}

	if student.Parent.Phone == "" || student.Parent.Address == "" || student.Parent.Email == "" {
		return errors.New("parent details and email are required", 400)
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")
	var existingParent models.User
	// Email as unique identifier for parent
	err = collection.FindOne(context.Background(), bson.M{"email": student.Parent.Email, "role": models.Parent}).Decode(&existingParent)

	// If parent is not yet added, add the parent
	var parentID primitive.ObjectID
	if err == nil {
		parentID = existingParent.ID
	} else if err == mongo.ErrNoDocuments {
		parentUser := student.Parent
		parentUser.Role = models.Parent

		parentUser.Details = &models.ParentDetails{
			Children: []primitive.ObjectID{},
		}

		parentID, err = AddUser(parentUser, username)
		if err != nil {
			log.Print(err)
			return err
		}
	} else {
		log.Print(err)
		return err
	}

	studentDocument := bson.D{
		{Key: "first_name", Value: student.Student.FirstName},
		{Key: "last_name", Value: student.Student.LastName},
		{Key: "parent_id", Value: parentID},
		{Key: "school_id", Value: schoolID},
	}

	studentsCollection := client.Database(viper.GetString("MONGO_DB")).Collection("students")
	result, err := studentsCollection.InsertOne(context.Background(), studentDocument)
	if err != nil {
		log.Print(err)
		return err
	}

	studentID := result.InsertedID.(primitive.ObjectID)

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": parentID},
		bson.M{"$push": bson.M{"details.children_id": studentID}},
	)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func UpdatePermittedSchoolStudentWithParents(id string, student models.SchoolStudentRequest, schoolID primitive.ObjectID) error {
    client, err := database.MongoConnection()
    if err != nil {
        log.Print(err)
        return err
    }

    collection := client.Database(viper.GetString("MONGO_DB")).Collection("students")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

	if err := CheckStudentAvailability(objectID, schoolID); err != nil {
        return err
    }

	// Pipeline to get the student and parent details
    var existingStudent models.SchoolStudentRequest
    pipeline := mongo.Pipeline{
        bson.D{{Key: "$match", Value: bson.M{"_id": objectID, "school_id": schoolID}}},
        bson.D{{Key: "$lookup", Value: bson.M{
            "from":         "users",
            "localField":   "parent_id",
            "foreignField": "_id",
            "as":           "parent",
        }}},
        bson.D{{Key: "$unwind", Value: bson.M{"path": "$parent"}}},
    }

	// Aggregate the pipeline
    cursor, err := collection.Aggregate(context.Background(), pipeline)
    if err != nil {
        log.Print(err)
        return err
    }
    defer cursor.Close(context.Background())

    if cursor.Next(context.Background()) {
        if err := cursor.Decode(&existingStudent); err != nil {
            log.Print(err)
            return err
        }
    }

    if (models.User{}) == existingStudent.Parent {
        return errors.New("parent not found", 404)
    }

    updateStudent := bson.M{
        "first_name": student.FirstName,
        "last_name":  student.LastName,
    }

    _, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID, "school_id": schoolID}, bson.M{"$set": updateStudent})
    if err != nil {
        log.Print(err)
        return err
    }

	// If the parent details are changed, update the parent details
    if (models.User{}) != student.Parent && student.Parent != existingStudent.Parent {
        parentCollection := client.Database(viper.GetString("MONGO_DB")).Collection("users")
        parentID := existingStudent.Parent.ID

        updateParent := bson.M{
            "first_name": student.Parent.FirstName,
            "last_name":  student.Parent.LastName,
            "email":      student.Parent.Email,
			"phone":    student.Parent.Phone,
            "address":  student.Parent.Address,
            "details": bson.M{
                "children": []primitive.ObjectID{objectID},
            },
        }

        _, err = parentCollection.UpdateOne(context.Background(), bson.M{"_id": parentID}, bson.M{"$set": updateParent})
        if err != nil {
            log.Print(err)
            return err
        }
    } else { // Else, parent remains the same
        _, err = client.Database(viper.GetString("MONGO_DB")).Collection("users").UpdateOne(
            context.Background(),
            bson.M{"_id": existingStudent.Parent.ID},
            bson.M{"$addToSet": bson.M{"parent_details.children": objectID}},
        )
        if err != nil {
            log.Print(err)
            return err
        }
    }

    return nil
}

func DeletePermittedSchoolStudentWithParents(id string, schoolID primitive.ObjectID) error {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("students")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err := CheckStudentAvailability(objectID, schoolID); err != nil {
		return err
	}

	var student models.Student
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID, "school_id": schoolID}).Decode(&student)
	if err != nil {
		log.Print("Error finding student: ", err)
		return err
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID, "school_id": schoolID})
	if err != nil {
		log.Print("Error deleting student: ", err)
		return err
	}

	// Remove the student from the parent's children list
	parentCollection := client.Database(viper.GetString("MONGO_DB")).Collection("users")
	_, err = parentCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": student.ParentID},
		bson.M{"$pull": bson.M{"parent_details.children": objectID}},
	)
	if err != nil {
		return err
	}

	var parent models.User
	err = parentCollection.FindOne(context.Background(), bson.M{"_id": student.ParentID}).Decode(&parent)
	if err != nil {
		return err
	}

	// If the children array is empty, delete the parent
	if len(parent.Details.(models.ParentDetails).Children) == 0 {
		_, err = parentCollection.DeleteOne(context.Background(), bson.M{"_id": student.ParentID})
		if err != nil {
			log.Print("Error in deleting parent: ", err)
			return err
		}
	}

	return nil
}


func CheckPermittedSchoolAccess(userID string) (primitive.ObjectID, error) {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return primitive.NilObjectID, err
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid user id", 400)
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return primitive.NilObjectID, errors.New("user not found", 404)
	}

	var schoolAdminDetails models.SchoolAdminDetails
	detailsBytes, err := bson.Marshal(user.Details)
	if err != nil {
		return primitive.NilObjectID, err
	}
	err = bson.Unmarshal(detailsBytes, &schoolAdminDetails)
	if err != nil || schoolAdminDetails.SchoolID.IsZero() {
		return primitive.NilObjectID, errors.New("user details are not in the correct format or school_id is missing", 400)
	}

	return schoolAdminDetails.SchoolID, nil
}

func CheckStudentAvailability(studentID primitive.ObjectID, schoolID primitive.ObjectID) error {
    client, err := database.MongoConnection()
    if err != nil {
        log.Print(err)
        return err
    }

    collection := client.Database(viper.GetString("MONGO_DB")).Collection("students")

    var student models.SchoolStudentRequest
    err = collection.FindOne(context.Background(), bson.M{"_id": studentID, "school_id": schoolID}).Decode(&student)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return errors.New("this student is not available in this school", 404)
		}
        log.Print(err)
        return err
    }

    return nil
}