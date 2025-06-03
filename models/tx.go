package models

// import (
// 	"time"
//
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusConfirmed TransactionStatus = "confirmed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusInvalid   TransactionStatus = "invalid"
)
