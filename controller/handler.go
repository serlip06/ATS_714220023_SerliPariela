package controller

import (

	"log"
	"context"
	"time"
	"github.com/gofiber/fiber/v2"

	inimodel "github.com/serlip06/pointsalesofkantin/model"
	cek "github.com/serlip06/pointsalesofkantin/module"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
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
