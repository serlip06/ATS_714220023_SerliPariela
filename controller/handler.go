package controller

import (
	"errors"
	"fmt"
	"log"
	"context"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/serlip06/ATS_714220023_SerliPariela/config"
	inimodel "github.com/serlip06/pointsalesofkantin/model"
	cek "github.com/serlip06/pointsalesofkantin/module"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	
)

var db = cek.MongoConnectdb("kantin")

// function untuk registrasi user 
func RegisterHandler(c *fiber.Ctx) error {
	if c.Method() != "POST" {
		return c.Status(fiber.StatusMethodNotAllowed).SendString("Invalid request method")
	}

	var registration inimodel.PendingRegistration
	if err := c.BodyParser(&registration); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to hash password")
	}
	registration.Password = string(hashedPassword)
	registration.SubmittedAt = time.Now()

	// Simpan data ke pending users
	if err := cek.SavePendingRegistration(registration, db); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save registration")
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]string{"message": "Registration successful"})
}

// function untuk login user 
func LoginHandler(c *fiber.Ctx) error {
	if c.Method() != "POST" {
		return c.Status(fiber.StatusMethodNotAllowed).SendString("Invalid request method")
	}

	var loginData inimodel.LoginRequest
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	collection := db.Collection("users")
	var user inimodel.User
	err := collection.FindOne(context.Background(), bson.M{"username": loginData.Username}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid username or password")
	}

	if !verifyPassword(user.Password, loginData.Password) {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid username or password")
	}

	return c.Status(fiber.StatusOK).JSON(inimodel.Response{
		Status:  "success",
		Message: "Login successful",
	})
}
// verifikasi paswoord 
func verifyPassword(storedPassword string, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	return err == nil
}

// function untuk acc user baru di admin
func ApproveRegistrationHandler(c *fiber.Ctx) error {
    id := c.Params("id") // Ambil ID dari URL

    if id == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "ID is required",
        })
    }

    // Panggil fungsi ApproveRegistration
    pending, user, err := cek.ApproveRegistration(id, db)
    if err != nil {
        log.Printf("Error in ConfirmRegistration: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to confirm registration",
        })
    }

    // Hash password jika belum di-hash sebelumnya
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to hash password",
        })
    }
    user.Password = string(hashedPassword)

    // Simpan user ke koleksi `users`
    collection := db.Collection("users")
    _, err = collection.InsertOne(context.Background(), user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to save user to users",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message":         "Registration confirmed",
        "pending": pending,
        "user":     user,
    })
}

//untuk men- Get data 
func GetAllUsers(c *fiber.Ctx) error {
	users, err := cek.GetAllUsers(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching users",
		})
	}
	return c.JSON(users)
}
// get user byid
// func GetUserByID(userID string, db *mongo.Database) (*inimodel.User, error) {
// 	collection := db.Collection("users")
// 	var user inimodel.User

// 	// Konversi ID string ke ObjectID
// 	objID, err := primitive.ObjectIDFromHex(userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid user ID format: %v", err)
// 	}

// 	// Mencari dokumen berdasarkan ID
// 	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
// 	if err != nil {
// 		log.Println("Error decoding user:", err) // Menambahkan log error

// 		if err == mongo.ErrNoDocuments {
// 			// Menangani kasus tidak ditemukan dokumen
// 			return nil, fmt.Errorf("user dengan ID %s tidak ditemukan", userID)
// 		}
// 		// Menangani error lain yang terjadi selama pencarian atau dekode
// 		return nil, fmt.Errorf("error retrieving user: %v", err)
// 	}

// 	// Mengembalikan objek user jika ditemukan
// 	return &user, nil
// }
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Wrong parameter",
		})
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid id parameter",
		})
	}
	ps, err := cek.GetUserFromID(objID, config.Ulbimongoconn, "users")
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"status":  http.StatusNotFound,
				"message": fmt.Sprintf("No data found for id %s", id),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error retrieving data for id %s", id),
		})
	}
	return c.JSON(ps)
}




func GetAllPendingRegistrations(c *fiber.Ctx) error {
	registrations, err := cek.GetAllPendingRegistrations(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching pending registrations",
		})
	}
	return c.JSON(registrations)
}

// handler untuk mnotifikasi 

// var db = cek.MongoConnectdb("kantin")

func UpdateProductHandler(c *fiber.Ctx) error {
	// Ambil data produk dari request body
	var produk inimodel.Produk
	if err := c.BodyParser(&produk); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Update produk ke database
	if err := updateProductInDB(produk); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update product"})
	}

	// Periksa apakah stok produk menipis atau habis
	if produk.Stok == 0 {
		// Stok habis, buat notifikasi stok habis
		notification := cek.CreateOutOfStockNotificationFromProduk(produk)
		if err := saveNotificationToDB(notification); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save notification"})
		}
	} else if produk.Stok < 10 {
		// Stok menipis, buat notifikasi stok menipis
		notification := cek.CreateLowStockNotificationFromProduk(produk, 5)
		if err := saveNotificationToDB(notification); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save notification"})
		}
	}

	// Notifikasi sukses
	return c.Status(200).JSON(fiber.Map{"message": "Product updated and notifications created successfully"})
}

func updateProductInDB(produk inimodel.Produk) error {
	// Gunakan db yang sudah dideklarasikan di luar fungsi
	collection := db.Collection("products")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": produk.IDProduk},  // Pastikan IDProduk sesuai dengan field yang ada
		bson.M{"$set": produk},           // Update produk dengan data terbaru
	)
	return err
}

//notif nambah barang
func AddProductHandler(c *fiber.Ctx) error {
	// Ambil data produk dari request body
	var produk inimodel.Produk
	if err := c.BodyParser(&produk); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Tambahkan produk baru ke database
	if err := addProductToDB(produk); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add product"})
	}

	// Buat notifikasi untuk produk baru
	notification := createNewProductNotification(produk)

	// Simpan notifikasi ke database
	if err := saveNotificationToDB(notification); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save notification"})
	}

	// Periksa apakah stok produk menipis atau habis
	if produk.Stok == 0 {
		// Stok habis, buat notifikasi stok habis
		notification := createOutOfStockNotification(produk)
		if err := saveNotificationToDB(notification); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save out of stock notification"})
		}
	} else if produk.Stok < 10 {
		// Stok menipis, buat notifikasi stok menipis
		notification := createLowStockNotification(produk, 5) // Ambang batas 5
		if err := saveNotificationToDB(notification); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save low stock notification"})
		}
	}

	// Notifikasi sukses
	return c.Status(201).JSON(fiber.Map{
		"message": fmt.Sprintf("Product '%s' added successfully and notifications created", produk.Nama_Produk),
	})
}

func createNewProductNotification(produk inimodel.Produk) interface{} {
	// Membuat notifikasi tentang produk baru
	notification := map[string]interface{}{
		"message": fmt.Sprintf("New product added: %s", produk.Nama_Produk),
		"product_id": produk.IDProduk,
		"product_name": produk.Nama_Produk,
		"created_at": time.Now(),
	}
	return notification
}

// Fungsi untuk membuat notifikasi stok habis
func createOutOfStockNotification(produk inimodel.Produk) interface{} {
	notification := map[string]interface{}{
		"message": fmt.Sprintf("Product out of stock: %s", produk.Nama_Produk),
		"product_id": produk.IDProduk,
		"created_at": time.Now(),
	}
	return notification
}

// Fungsi untuk membuat notifikasi stok menipis
func createLowStockNotification(produk inimodel.Produk, threshold int) interface{} {
    // Cek apakah stok produk lebih rendah dari threshold
    if produk.Stok < threshold {
        notification := map[string]interface{}{
            "message": fmt.Sprintf("Low stock alert: %s, only %d items left", produk.Nama_Produk, produk.Stok),
            "product_id": produk.IDProduk,
            "created_at": time.Now(),
        }
        return notification
    }
    return nil // Jika stok tidak rendah, kembalikan nil
}

func addProductToDB(produk inimodel.Produk) error {
	// Gunakan db yang sudah dideklarasikan di luar fungsi
	collection := db.Collection("produk")
	_, err := collection.InsertOne(context.Background(), produk) // Insert produk baru ke database
	return err
}

func saveNotificationToDB(notification interface{}) error {
	// Gunakan db yang sudah dideklarasikan di luar fungsi
	collection := db.Collection("notifications")
	_, err := collection.InsertOne(context.Background(), notification) // Menyimpan notifikasi ke database
	return err
}

