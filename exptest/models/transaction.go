package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Amount      float64            `json:"amount" bson:"amount"`
	Description string             `json:"description" bson:"description"`
	Date        time.Time          `json:"date" bson:"date"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Type        TransactionType    `json:"type" bson:"type"`
}
