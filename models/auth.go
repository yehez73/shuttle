package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshToken struct {
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
	ExpiredAt    time.Time `json:"expired_at" bson:"expired_at"`
}