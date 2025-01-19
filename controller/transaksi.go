package controller

import (
	//"errors"
	"time"
	"context"
	//"fmt"
	//"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
	"github.com/serlip06/ATS_714220023_SerliPariela/config"
	inimodel "github.com/serlip06/pointsalesofkantin/model"
	cek "github.com/serlip06/pointsalesofkantin/module"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// Fungsi untuk menambahkan transaksi baru
func InsertTransaksi(c *fiber.Ctx) error {
	var input struct {
		IDUser          primitive.ObjectID `json:"id_user"`
		Username        string             `json:"username"`
		Items           []inimodel.CartItem   `json:"items"`
		MetodePembayaran string            `json:"metode_pembayaran"`
	}

	// Parse body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	// Validasi Items tidak boleh kosong
	if len(input.Items) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Items cannot be empty",
		})
	}

	// Hitung total harga
	calculatedTotal := calculateTotalHarga(input.Items)

	// Buat objek transaksi
	transaksi := inimodel.Transaksi{
		IDTransaksi:      primitive.NewObjectID(),
		IDUser:           input.IDUser,
		Username:         input.Username,
		Items:            input.Items,
		TotalHarga:       calculatedTotal,
		MetodePembayaran: input.MetodePembayaran,
		CreatedAt:        time.Now(),
	}

	// Simpan transaksi ke database
	result, err := cek.InsertTransaksiToDatabase("kantin", "kantin_transaksi", transaksi)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to insert transaction",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Transaction successfully created",
		"transaction": result,
	})
}


// Fungsi untuk menghitung total harga
func calculateTotalHarga(items []inimodel.CartItem) int {
	total := 0
	for _, item := range items {
		total += item.SubTotal
	}
	return total
}

// Fungsi untuk mendapatkan transaksi berdasarkan ID
func GetTransaksiByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID format",
		})
	}

	collection := config.Ulbimongoconn.Collection("kantin_transaksi")
	var transaksi inimodel.Transaksi
	err = collection.FindOne(c.Context(), bson.M{"_id": objID}).Decode(&transaksi)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": "Transaction not found",
		})
	}

	return c.JSON(transaksi)
}

// Fungsi untuk mendapatkan semua transaksi
func GetAllTransaksi(c *fiber.Ctx) error {
	collection := config.Ulbimongoconn.Collection("kantin_transaksi")
	cursor, err := collection.Find(c.Context(), bson.D{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to fetch transactions",
		})
	}
	defer cursor.Close(context.TODO())

	var transactions []inimodel.Transaksi
	if err := cursor.All(context.TODO(), &transactions); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to decode transactions",
		})
	}

	return c.JSON(transactions)
}

// Fungsi untuk update transaksi
func UpdateTransaksi(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID format",
		})
	}

	var input inimodel.Transaksi
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid input",
		})
	}

	collection := config.Ulbimongoconn.Collection("kantin_transaksi")
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": input}

	result, err := collection.UpdateOne(c.Context(), filter, update)
	if err != nil || result.ModifiedCount == 0 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to update transaction",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Transaction successfully updated",
	})
}

// Fungsi untuk menghapus transaksi
func DeleteTransaksiByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID format",
		})
	}

	collection := config.Ulbimongoconn.Collection("kantin_transaksi")
	result, err := collection.DeleteOne(c.Context(), bson.M{"_id": objID})
	if err != nil || result.DeletedCount == 0 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to delete transaction",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Transaction successfully deleted",
	})
}