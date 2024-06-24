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

func GetExpenses(c *fiber.Ctx) error {
	query := bson.D{{}}

	cursor, err := config.Db.Collection("expenses").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var expenses []models.Transaction

	if err := cursor.All(c.Context(), &expenses); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(expenses)
}

func CreateExpense(c *fiber.Ctx) error {
	collection := config.Db.Collection("expenses")

	expense := new(models.Transaction)

	if err := c.BodyParser(expense); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	expense.ID = primitive.NewObjectID()
	expense.Date = time.Now()

	insertionResult, err := collection.InsertOne(c.Context(), expense)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdExpense := &models.Transaction{}
	createdRecord.Decode(createdExpense)

	return c.Status(201).JSON(createdExpense)
}

func UpdateExpense(c *fiber.Ctx) error {
	idParam := c.Params("id")

	expenseID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).SendString("Invalid expense ID")
	}

	expense := new(models.Transaction)

	if err := c.BodyParser(expense); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	query := bson.D{{Key: "_id", Value: expenseID}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "title", Value: expense.Title},
				{Key: "amount", Value: expense.Amount},
				{Key: "description", Value: expense.Description},
				{Key: "date", Value: time.Now()}, // Update date to current time
			},
		},
	}

	err = config.Db.Collection("expenses").FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(404)
		}
		return c.Status(500).SendString(err.Error())
	}

	expense.ID = expenseID

	return c.Status(200).JSON(expense)
}

func DeleteExpense(c *fiber.Ctx) error {
	expenseID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid expense ID")
	}

	query := bson.D{{Key: "_id", Value: expenseID}}
	result, err := config.Db.Collection("expenses").DeleteOne(c.Context(), &query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if result.DeletedCount < 1 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON("Expense deleted")
}
