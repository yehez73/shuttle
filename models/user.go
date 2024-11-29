package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string
type Gender string

const (
	SuperAdmin  Role = "superadmin"
	SchoolAdmin Role = "schooladmin"
	Parent      Role = "parent"
	Driver      Role = "driver"

	Female Gender = "female"
	Male   Gender = "male"
)

func ParseGender(genderStr string) (Gender, error) {
	switch genderStr {
	case "male":
		return Male, nil
	case "female":
		return Female, nil
	default:
		return "", fmt.Errorf("invalid gender value")
	}
}

func ParseRole(roleStr string) (Role, error) {
	switch roleStr {
	case "superadmin":
		return SuperAdmin, nil
	case "schooladmin":
		return SchoolAdmin, nil
	case "parent":
		return Parent, nil
	case "driver":
		return Driver, nil
	default:
		return "", fmt.Errorf("invalid role value")
	}
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Picture   string             `json:"picture" bson:"picture"`
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required"`
	Gender    Gender             `json:"gender" bson:"gender" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required"`
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
	Children []primitive.ObjectID `json:"children_id" bson:"children_id"`
}

type DriverDetails struct {
	LicenseNumber string             `json:"license_number" bson:"license_number" validate:"required"`
	SchoolID      primitive.ObjectID `json:"school_id" bson:"school_id"`
	VehicleID     primitive.ObjectID `json:"vehicle_id" bson:"vehicle_id"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Picture   string             `json:"picture" bson:"picture"`
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required"`
	Gender    Gender             `json:"gender" bson:"gender" validate:"required"`
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
