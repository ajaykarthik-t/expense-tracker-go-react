package handlers

import (
	"time"

	"github.com/ajay/exptest/config"
	"github.com/ajay/exptest/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTransactions(c *fiber.Ctx) error {
	// Get the user ID from the context
	userID := c.Locals("userID").(primitive.ObjectID)
	query := bson.D{{Key: "user_id", Value: userID}}

	cursor, err := config.Db.Collection("transactions").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var transactions []models.Transaction

	if err := cursor.All(c.Context(), &transactions); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(transactions)
}

func CreateTransaction(c *fiber.Ctx) error {
	collection := config.Db.Collection("transactions")

	transaction := new(models.Transaction)

	if err := c.BodyParser(transaction); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Validate transaction type
	if transaction.Type != models.Income && transaction.Type != models.Expense {
		return c.Status(400).SendString("Invalid transaction type. Must be 'income' or 'expense'")
	}

	// Get the user ID from the context
	userID := c.Locals("userID").(primitive.ObjectID)
	transaction.UserID = userID
	transaction.ID = primitive.NewObjectID()
	transaction.Date = time.Now()

	insertionResult, err := collection.InsertOne(c.Context(), transaction)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdTransaction := &models.Transaction{}
	createdRecord.Decode(createdTransaction)

	return c.Status(201).JSON(createdTransaction)
}

func UpdateTransaction(c *fiber.Ctx) error {
	idParam := c.Params("id")

	transactionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).SendString("Invalid transaction ID")
	}

	transaction := new(models.Transaction)

	if err := c.BodyParser(transaction); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Validate transaction type
	if transaction.Type != models.Income && transaction.Type != models.Expense {
		return c.Status(400).SendString("Invalid transaction type. Must be 'income' or 'expense'")
	}

	userID := c.Locals("userID").(primitive.ObjectID)

	query := bson.D{{Key: "_id", Value: transactionID}, {Key: "user_id", Value: userID}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "title", Value: transaction.Title},
				{Key: "amount", Value: transaction.Amount},
				{Key: "description", Value: transaction.Description},
				{Key: "date", Value: time.Now()},
				{Key: "type", Value: transaction.Type}, // Update type
			},
		},
	}

	err = config.Db.Collection("transactions").FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(404)
		}
		return c.Status(500).SendString(err.Error())
	}

	transaction.ID = transactionID
	transaction.UserID = userID

	return c.Status(200).JSON(transaction)
}

func DeleteTransaction(c *fiber.Ctx) error {
	transactionID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid transaction ID")
	}

	userID := c.Locals("userID").(primitive.ObjectID)

	query := bson.D{{Key: "_id", Value: transactionID}, {Key: "user_id", Value: userID}}
	result, err := config.Db.Collection("transactions").DeleteOne(c.Context(), &query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if result.DeletedCount < 1 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON("Transaction deleted")
}
