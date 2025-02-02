package controller

import (
	//"errors"
	"fmt"
	"context"
	"time"
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

// InsertTransaksi - Menambahkan transaksi baru
func InsertTransaksi(c *fiber.Ctx) error {
	var input struct {
		IDUser           primitive.ObjectID   `json:"id_user"`
		IDCartItem       []primitive.ObjectID `json:"id_cartitem"`
		MetodePembayaran string               `json:"metode_pembayaran"`
		BuktiPembayaran  string               `json:"bukti_pembayaran"`
		Status           string               `json:"status"`
		Alamat           string               `json:"alamat,omitempty"`
	}

	// Parse body JSON
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	// Validasi CartItem tidak boleh kosong
	if len(input.IDCartItem) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Cart items cannot be empty",
		})
	}

	// Hitung total harga berdasarkan CartItem
	totalHarga, err := calculateTotalHarga(input.IDCartItem) // ✅ Pastikan menangkap 2 nilai
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to calculate total price",
			"error":   err.Error(),
		})
	}

	// Buat objek transaksi
	transaksi := inimodel.Transaksi{
		IDTransaksi:      primitive.NewObjectID(),
		IDUser:           input.IDUser,
		IDCartItem:       input.IDCartItem,
		TotalHarga:       totalHarga, // ✅ Menggunakan total harga yang dihitung
		MetodePembayaran: input.MetodePembayaran,
		CreatedAt:        time.Now(),
		BuktiPembayaran:  input.BuktiPembayaran,
		Status:           input.Status,
		Alamat:           input.Alamat,
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

func calculateTotalHarga(idCartItems []primitive.ObjectID) (int, error) {
	collection := config.Ulbimongoconn.Collection("cart_items")
	total := 0

	for _, cartItemID := range idCartItems {
		var cartItem inimodel.CartItem
		err := collection.FindOne(context.TODO(), bson.M{"_id": cartItemID}).Decode(&cartItem)
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve cart item: %w", err)
		}
		total += cartItem.SubTotal // ✅ Pastikan model.CartItem memiliki field SubTotal
	}

	return total, nil
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

	var input struct {
		IDUser           primitive.ObjectID   `json:"id_user"`
		IDCartItem       []primitive.ObjectID `json:"id_cartitem"`
		MetodePembayaran string               `json:"metode_pembayaran"`
		BuktiPembayaran  string               `json:"bukti_pembayaran"`
		Status           string               `json:"status"`
		Alamat           string               `json:"alamat,omitempty"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid input",
		})
	}

	// Hitung total harga berdasarkan CartItem
	calculatedTotal, err := calculateTotalHarga(input.IDCartItem)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to calculate total price",
			"error":   err.Error(),
		})
	}

	collection := config.Ulbimongoconn.Collection("kantin_transaksi")
	filter := bson.M{"_id": objID}

	// Update transaksi dengan data terbaru
	update := bson.M{"$set": bson.M{
		"id_user":           input.IDUser,
		"id_cartitem":       input.IDCartItem,
		"total_harga":       calculatedTotal,
		"metode_pembayaran": input.MetodePembayaran,
		"bukti_pembayaran":  input.BuktiPembayaran,
		"status":            input.Status,
		"alamat":            input.Alamat,
	}}

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

// get transaksi dari id user
func GetTransaksiByUserID(c *fiber.Ctx) error {
    userID := c.Params("id") // Ambil userID dari URL parameter

    // Ambil transaksi berdasarkan ID user
    transaksis, err := cek.GetAllTransaksiByIDUser(userID, config.Ulbimongoconn) // Memanggil fungsi backend yang sudah ada
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  fiber.StatusBadRequest,
            "message": "Invalid user ID format",
            "error":   err.Error(),
        })
    }

    // Jika tidak ada transaksi ditemukan
    if len(transaksis) == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "status":  fiber.StatusNotFound,
            "message": "No transactions found for this user",
        })
    }

    // Kembalikan transaksi yang ditemukan
    return c.JSON(transaksis)
}


//cadangan 
// func calculateTotalHarga(items []inimodel.CartItem) int {
// 	total := 0
// 	for _, item := range items {
// 		total += item.SubTotal
// 	}
// 	return total
// }

// Helper function untuk memeriksa apakah string hanya mengandung karakter hexadecimal
// func isValidHex(s string) bool {
// 	for _, c := range s {
// 		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
// 			return false
// 		}
// 	}
// 	return true
// }
