package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Point struct {
	Name      string  `json:"name" bson:"name" validate:"required"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type RoadRoute struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RouteName         string             `json:"route_name" bson:"route_name" validate:"required"`
	Points            []Point            `json:"points" bson:"points" validate:"required"`
	Status            string             `json:"status" bson:"status" validate:"required"`
	AssignedVehicleID primitive.ObjectID `json:"assigned_vehicle_id" bson:"assigned_vehicle_id"`
}