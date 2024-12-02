package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Point struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type RoadRoute struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RouteName string             `json:"route_name" bson:"route_name" validate:"required"`
	Points    []Point            `json:"points" bson:"points" validate:"required"`
	Status    string             `json:"status" bson:"status" validate:"required"`
	SchoolID  primitive.ObjectID `json:"school_id" bson:"school_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	CreatedBy string			 `json:"created_by" bson:"created_by"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	UpdatedBy string			 `json:"updated_by" bson:"updated_by"`
	DeleredAt time.Time          `json:"deleted_at" bson:"deleted_at"`
	DeletedBy string			 `json:"deleted_by" bson:"deleted_by"`
}
