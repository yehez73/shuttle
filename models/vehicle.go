package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vehicle struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VehicleName   string             `json:"vehicle_name" bson:"vehicle_name" validate:"required"`
	VehicleNumber string             `json:"vehicle_number" bson:"vehicle_number" validate:"required"`
	VehicleType   string             `json:"vehicle_type" bson:"vehicle_type" validate:"required"`
	Colour        string             `json:"colour" bson:"colour" validate:"required"`
	Seats         int                `json:"seats" bson:"seats" validate:"required"`
	Status        string             `json:"status" bson:"status"`
}
