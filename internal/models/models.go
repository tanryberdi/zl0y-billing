package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the postgresql.
type User struct {
	ID           int       `json:"id" db:"id"`
	Login        string    `json:"login" db:"login"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Balance      int       `json:"balance" db:"balance"` // Balance in cents
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// The Report represents a report in the MongoDB.
type Report struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ReportID          string             `json:"report_id" bson:"report_id"`
	UserID            *int               `json:"user_id,omitempty" bson:"user_id,omitempty"`
	ClientGeneratedID string             `json:"client_generated_id" bson:"client_generated_id"`
	IsPurchased       bool               `json:"is_purchased" bson:"is_purchased"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
}

// Auth request/response models
type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

// User request/response models
type LinkAnonymousRequest struct {
	ClientGeneratedID string `json:"client_generated_id" binding:"required"`
}

type ReportsResponse struct {
	Reports []Report `json:"reports"`
	Total   int64    `json:"total"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
}

// Mock request models
type CreateReportRequest struct {
	ClientGeneratedID string `json:"client_generated_id" binding:"required"`
}

// Error response model
type ErrorResponse struct {
	Error string `json:"error"`
}
