package models

import (
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
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password" bson:"password" validate:"required"`
	Role      Role               `json:"role" bson:"role" validate:"required"`
	RoleCode  string             `json:"role_code" bson:"role_code"`
	Status    string             `json:"status" bson:"status"`

	SuperAdminDetails  *SuperAdminDetails  `json:"super_admin_details" bson:"super_admin_details,omitempty"`
	SchoolAdminDetails *SchoolAdminDetails `json:"school_admin_details" bson:"school_admin_details,omitempty"`
	ParentDetails      *ParentDetails      `json:"parent_details" bson:"parent_details,omitempty"`
	DriverDetails      *DriverDetails      `json:"driver_details" bson:"driver_details,omitempty"`
}

type SuperAdminDetails struct {
	Phone   string `json:"phone" bson:"phone" validate:"required"`
	Address string `json:"address" bson:"address" validate:"required"`
}

type SchoolAdminDetails struct {
	SchoolID primitive.ObjectID `json:"school_id" bson:"school_id" validate:"required"`
	Phone    string             `json:"phone" bson:"phone" validate:"required"`
	Address  string             `json:"address" bson:"address" validate:"required"`
}

type ParentDetails struct {
	Children []primitive.ObjectID `json:"children" bson:"children"`
	Phone    string               `json:"phone" bson:"phone" validate:"required"`
	Address  string               `json:"address" bson:"address" validate:"required"`
}

type DriverDetails struct {
	VehicleID primitive.ObjectID `json:"vehicle_id" bson:"vehicle_id" validate:"required"`
	Phone     string             `json:"phone" bson:"phone" validate:"required"`
	Address   string             `json:"address" bson:"address" validate:"required"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Role      Role               `json:"role" bson:"role" validate:"required"`
	RoleCode  string             `json:"role_code" bson:"role_code"`
	Status    string             `json:"status" bson:"status"`
}