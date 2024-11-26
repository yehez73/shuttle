package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	SuperAdmin  Role = "superadmin"
	SchoolAdmin Role = "schooladmin"
	Parent      Role = "parent"
	Driver      Role = "driver"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Picture   string             `json:"picture" bson:"picture"`
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password" bson:"password" validate:"required"`
	Role      Role               `json:"role" bson:"role" validate:"required"`
	RoleCode  string             `json:"role_code" bson:"role_code"`
	Phone     string             `json:"phone" bson:"phone" validate:"required"`
	Address   string             `json:"address" bson:"address" validate:"required"`
	Status    string             `json:"status" bson:"status"`

	Details interface{} `json:"details" bson:"details"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	CreatedBy string    `json:"created_by" bson:"created_by"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy string    `json:"updated_by" bson:"updated_by"`
}

type SchoolAdminDetails struct {
	SchoolID primitive.ObjectID `json:"school_id" bson:"school_id" validate:"required"`
}

type ParentDetails struct {
	Children []primitive.ObjectID `json:"children" bson:"children"`
}

type DriverDetails struct {
	VehicleID primitive.ObjectID `json:"vehicle_id" bson:"vehicle_id" validate:"required"`
	Vehicle   Vehicle            `json:"vehicle" bson:"vehicle" validate:"required"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Role      Role               `json:"role" bson:"role" validate:"required"`
	RoleCode  string             `json:"role_code" bson:"role_code"`
	Phone     string             `json:"phone" bson:"phone" validate:"required"`
	Address   string             `json:"address" bson:"address" validate:"required"`
	Status    string             `json:"status" bson:"status"`

	Details interface{} `json:"details" bson:"details,omitempty"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	CreatedBy string    `json:"created_by" bson:"created_by"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy string    `json:"updated_by" bson:"updated_by"`
}
