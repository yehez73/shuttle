package models

import (
	"shuttle/models/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required"`
	Class     string             `json:"class" bson:"class" validate:"required"`
	ParentID  primitive.ObjectID `json:"parent_id" bson:"parent_id"`
	SchoolID  primitive.ObjectID `json:"school_id" bson:"school_id"`
}

type SchoolStudentRequest struct {
	Student `json:"student" validate:"required"`
	Parent  entity.User `json:"parent" validate:"required"`
}

type SchoolStudentParentResponse struct {
	Student `json:"student"`
	Parent  ParentResponse `json:"parent"`
}

type ParentResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Role      entity.User               `json:"role" bson:"role" validate:"required"`
}

type StudentResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required"`
	ParentID  primitive.ObjectID `json:"parent_id" bson:"parent_id"`
	SchoolID  primitive.ObjectID `json:"school_id" bson:"school_id"`
}
