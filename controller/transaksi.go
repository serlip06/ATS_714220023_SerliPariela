package controller

import (
	//"errors"
	"bytes"
	//"fmt"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"io"
	"mime/multipart"
	//"os"
	"time"
	//"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
	"github.com/serlip06/ATS_714220023_SerliPariela/config"
	inimodel "github.com/serlip06/pointsalesofkantin/model"
	cek "github.com/serlip06/pointsalesofkantin/module"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	//"strings"
)

// Fungsi untuk menambahkan transaksi baru
// func InsertTransaksi(c *fiber.Ctx) error {
// 	var input struct {
// 		IDUser           primitive.ObjectID   `json:"id_user"`
// 		IDCartItem       []primitive.ObjectID `json:"id_cartitem"`
// 		MetodePembayaran string               `json:"metode_pembayaran"`
// 		BuktiPembayaran  string               `json:"bukti_pembayaran"`
// 		Status           string               `json:"status"`
// 		Alamat           string               `json:"alamat,omitempty"`
// 	}

// 	// Parse body JSON
// 	if err := c.BodyParser(&input); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"status":  http.StatusBadRequest,
// 			"message": "Invalid input",
// 			"error":   err.Error(),
// 		})
// 	}

// 	// Validasi CartItem tidak boleh kosong
// 	if len(input.IDCartItem) == 0 {
// 		// Cari item keranjang berdasarkan id_user
// 		var cartItems []struct {
// 			ID primitive.ObjectID `bson:"_id"`
// 		}

// 		cursor, err := config.Ulbimongoconn.Collection("cart_items").Find(
// 			c.Context(),
// 			bson.M{"id_user": input.IDUser},
// 		)
// 		if err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 				"status":  http.StatusInternalServerError,
// 				"message": "Failed to fetch cart items",
// 				"error":   err.Error(),
// 			})
// 		}
// 		defer cursor.Close(c.Context()) // Tutup cursor setelah selesai digunakan

// 		// Decode hasil query ke dalam slice cartItems
// 		if err = cursor.All(c.Context(), &cartItems); err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 				"status":  http.StatusInternalServerError,
// 				"message": "Failed to decode cart items",
// 				"error":   err.Error(),
// 			})
// 		}

// 		// Jika tidak ada item dalam keranjang
// 		if len(cartItems) == 0 {
// 			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 				"status":  http.StatusBadRequest,
// 				"message": "Cart items cannot be empty",
// 			})
// 		}

// 		// Ambil ID dari hasil query
// 		for _, item := range cartItems {
// 			input.IDCartItem = append(input.IDCartItem, item.ID)
// 		}
// 	}

// 	// Hitung total harga berdasarkan CartItem
// 	totalHarga, err := calculateTotalHarga(input.IDCartItem)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": "Failed to calculate total price",
// 			"error":   err.Error(),
// 		})
// 	}
// Fungsi untuk menambahkan transaksi baru
func InsertTransaksi(c *fiber.Ctx) error {
	var input struct {
		IDUser           primitive.ObjectID   `json:"id_user"`           // Menggunakan ObjectID langsung
		IDCartItem       []primitive.ObjectID `json:"id_cartitem"`       // CartItem sebagai array of ObjectID
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
		// Cari item keranjang berdasarkan id_user
		var cartItems []struct {
			ID primitive.ObjectID `bson:"_id"`
		}

		cursor, err := config.Ulbimongoconn.Collection("cart_items").Find(
			c.Context(),
			bson.M{"id_user": input.IDUser}, // Tidak perlu konversi, langsung gunakan ObjectID
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "Failed to fetch cart items",
				"error":   err.Error(),
			})
		}
		defer cursor.Close(c.Context()) // Tutup cursor setelah selesai digunakan

		// Decode hasil query ke dalam slice cartItems
		if err = cursor.All(c.Context(), &cartItems); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "Failed to decode cart items",
				"error":   err.Error(),
			})
		}

		// Jika tidak ada item dalam keranjang
		if len(cartItems) == 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"status":  http.StatusBadRequest,
				"message": "Cart items cannot be empty",
			})
		}

		// Ambil ID dari hasil query
		for _, item := range cartItems {
			input.IDCartItem = append(input.IDCartItem, item.ID)
		}
	}

	// Hitung total harga berdasarkan CartItem
	totalHarga, err := calculateTotalHarga(input.IDCartItem)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to calculate total price",
			"error":   err.Error(),
		})
	}

	// Proses upload gambar bukti pembayaran
	file, err := c.FormFile("bukti_pembayaran") // Mengambil gambar dari form
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Failed to process image file: " + err.Error(),
		})
	}

	// ✅ Pastikan menambahkan c sebagai parameter pertama
	imageURL, err := UploadImageToGitHub(c, file, "bukti_pembayaran")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to upload payment proof: " + err.Error(),
		})
	}

	// Buat objek transaksi dengan URL gambar
	transaksi := inimodel.Transaksi{
		IDTransaksi:      primitive.NewObjectID(),
		IDUser:           input.IDUser, // Tetap menggunakan ObjectID
		IDCartItem:       input.IDCartItem,
		TotalHarga:       totalHarga,
		MetodePembayaran: input.MetodePembayaran,
		CreatedAt:        time.Now(),
		BuktiPembayaran:  imageURL, // Menyimpan URL gambar bukti pembayaran
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

	// Kembalikan respon sukses
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Transaction successfully created",
		"transaction": result,
	})
}

// Fungsi untuk menambahkan transaksi baru
// func InsertTransaksi(c *fiber.Ctx) error {
// 	var input struct {
// 		IDUser           string               `json:"id_user"` // Pastikan ini adalah string yang valid
// 		IDCartItem       []primitive.ObjectID `json:"id_cartitem"`
// 		MetodePembayaran string               `json:"metode_pembayaran"`
// 		BuktiPembayaran  string               `json:"bukti_pembayaran"`
// 		Status           string               `json:"status"`
// 		Alamat           string               `json:"alamat,omitempty"`
// 	}

// 	// Parse body JSON
// 	if err := c.BodyParser(&input); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"status":  http.StatusBadRequest,
// 			"message": "Invalid input",
// 			"error":   err.Error(),
// 		})
// 	}

// 	// Mengonversi ID user yang dikirimkan menjadi ObjectID yang valid
// 	idUser, err := primitive.ObjectIDFromHex(input.IDUser)
// 	if err != nil {
// 		fmt.Println("Invalid ID format:", err) // Debugging
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"status":  http.StatusBadRequest,
// 			"message": "Invalid user ID format", // Pesan jika format ID salah
// 			"error":   err.Error(),
// 		})
// 	}

// 	// Validasi CartItem tidak boleh kosong
// 	if len(input.IDCartItem) == 0 {
// 		// Cari item keranjang berdasarkan id_user
// 		var cartItems []struct {
// 			ID primitive.ObjectID `bson:"_id"`
// 		}

// 		cursor, err := config.Ulbimongoconn.Collection("cart_items").Find(
// 			c.Context(),
// 			bson.M{"id_user": idUser}, // ✅ Gunakan idUser yang sudah dikonversi
// 		)
// 		if err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 				"status":  http.StatusInternalServerError,
// 				"message": "Failed to fetch cart items",
// 				"error":   err.Error(),
// 			})
// 		}
// 		defer cursor.Close(c.Context())

// 		// Decode hasil query ke dalam slice cartItems
// 		if err = cursor.All(c.Context(), &cartItems); err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 				"status":  http.StatusInternalServerError,
// 				"message": "Failed to decode cart items",
// 				"error":   err.Error(),
// 			})
// 		}

// 		// Jika tidak ada item dalam keranjang
// 		if len(cartItems) == 0 {
// 			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 				"status":  http.StatusBadRequest,
// 				"message": "Cart items cannot be empty",
// 			})
// 		}

// 		// Ambil ID dari hasil query
// 		for _, item := range cartItems {
// 			input.IDCartItem = append(input.IDCartItem, item.ID)
// 		}
// 	}

// 	// Hitung total harga berdasarkan CartItem
// 	totalHarga, err := calculateTotalHarga(input.IDCartItem)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": "Failed to calculate total price",
// 			"error":   err.Error(),
// 		})
// 	}

// 	// Proses upload gambar bukti pembayaran
// 	file, err := c.FormFile("bukti_pembayaran") // Mengambil gambar dari form
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"status":  http.StatusBadRequest,
// 			"message": "Failed to process image file: " + err.Error(),
// 		})
// 	}

// 	// ✅ Panggil fungsi UploadImageToGitHub (tidak ada perubahan)
// 	imageURL, err := UploadImageToGitHub(c, file, "bukti_pembayaran")
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": "Failed to upload payment proof: " + err.Error(),
// 		})
// 	}

// 	// Buat objek transaksi dengan URL gambar
// 	transaksi := inimodel.Transaksi{
// 		IDTransaksi:      primitive.NewObjectID(),
// 		IDUser:           idUser, // ✅ Gunakan ObjectID yang sudah dikonversi
// 		IDCartItem:       input.IDCartItem,
// 		TotalHarga:       totalHarga,
// 		MetodePembayaran: input.MetodePembayaran,
// 		CreatedAt:        time.Now(),
// 		BuktiPembayaran:  imageURL, // ✅ Menyimpan URL gambar bukti pembayaran
// 		Status:           input.Status,
// 		Alamat:           input.Alamat,
// 	}

// 	// Simpan transaksi ke database
// 	result, err := cek.InsertTransaksiToDatabase("kantin", "kantin_transaksi", transaksi)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": "Failed to insert transaction",
// 			"error":   err.Error(),
// 		})
// 	}

// 	return c.Status(http.StatusOK).JSON(fiber.Map{
// 		"status":      http.StatusOK,
// 		"message":     "Transaction successfully created",
// 		"transaction": result,
// 	})
// }

// 	// Proses upload gambar bukti pembayaran
// 	file, err := c.FormFile("bukti_pembayaran") // Mengambil gambar dari form
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"status":  http.StatusBadRequest,
// 			"message": "Failed to process image file: " + err.Error(),
// 		})
// 	}

// 	// ✅ Pastikan menambahkan `c` sebagai parameter pertama
// 	imageURL, err := UploadImageToGitHub(c, file, "bukti_pembayaran")
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": "Failed to upload payment proof: " + err.Error(),
// 		})
// 	}

// 	// Buat objek transaksi dengan URL gambar
// 	transaksi := inimodel.Transaksi{
// 		IDTransaksi:      primitive.NewObjectID(),
// 		IDUser:           input.IDUser,
// 		IDCartItem:       input.IDCartItem,
// 		TotalHarga:       totalHarga,
// 		MetodePembayaran: input.MetodePembayaran,
// 		CreatedAt:        time.Now(),
// 		BuktiPembayaran:  imageURL, // ✅ Menyimpan URL gambar bukti pembayaran
// 		Status:           input.Status,
// 		Alamat:           input.Alamat,
// 	}

// 	// Simpan transaksi ke database
// 	result, err := cek.InsertTransaksiToDatabase("kantin", "kantin_transaksi", transaksi)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": "Failed to insert transaction",
// 			"error":   err.Error(),
// 		})
// 	}

// 	return c.Status(http.StatusOK).JSON(fiber.Map{
// 		"status":      http.StatusOK,
// 		"message":     "Transaction successfully created",
// 		"transaction": result,
// 	})
// }

//untuk di foldernya
// func getFileSHA(githubToken, repoOwner, repoName, filePath string) (string, error) {
//     url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, filePath)

//     req, _ := http.NewRequest("GET", url, nil)
//     req.Header.Set("Authorization", "Bearer "+githubToken)
//     req.Header.Set("Accept", "application/vnd.github.v3+json")

//     resp, err := http.DefaultClient.Do(req)
//     if err != nil {
//         return "", fmt.Errorf("failed to fetch file info: %w", err)
//     }
//     defer resp.Body.Close()

//     if resp.StatusCode == 404 {
//         return "", nil // File belum ada, jadi tidak perlu SHA
//     }

//     var result map[string]interface{}
//     if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
//         return "", fmt.Errorf("failed to parse GitHub API response: %w", err)
//     }

//     sha, ok := result["sha"].(string)
//     if !ok {
//         return "", fmt.Errorf("failed to get SHA from GitHub response")
//     }

//     return sha, nil
// }

// upload gitub
// func UploadImageToGitHub(file *multipart.FileHeader, productName string) (string, error) {
//     githubToken := os.Getenv("GH_ACCESS_TOKEN") // Pastikan token GitHub sudah di-set
//     repoOwner := "serlip06" // Ganti dengan username GitHub kamu
//     repoName := "Gambar" // Ganti dengan nama repository kamu
//     filePath := fmt.Sprintf("bukti pembayaran/%d_%s.jpg", time.Now().Unix(), productName)

//     fileContent, err := file.Open()
//     if err != nil {
//         return "", fmt.Errorf("failed to open image file: %w", err)
//     }
//     defer fileContent.Close()

// 	imageData, err := io.ReadAll(fileContent)
//     if err != nil {
//         return "", fmt.Errorf("failed to read image file: %w", err)
//     }

//     encodedImage := base64.StdEncoding.EncodeToString(imageData)
//     payload := map[string]string{
//         "message": fmt.Sprintf("Add image for product %s", productName),
//         "content": encodedImage,
//     }
//     payloadBytes, _ := json.Marshal(payload)

//     req, _ := http.NewRequest("PUT", fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, filePath), bytes.NewReader(payloadBytes))
// 	// req.Header.Set("Authorization", "token "+githubToken)
//     req.Header.Set("Authorization", "Bearer "+githubToken)
//     req.Header.Set("Content-Type", "application/json")

//     resp, err := http.DefaultClient.Do(req)
//     if err != nil {
//         return "", fmt.Errorf("failed to upload image to GitHub: %w", err)
//     }
//     defer resp.Body.Close()

//     if resp.StatusCode != http.StatusCreated {
// 		body, _ := io.ReadAll(resp.Body)

//         return "", fmt.Errorf("GitHub API error: %s", body)
//     }

//     var result map[string]interface{}
//     if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
//         return "", fmt.Errorf("failed to parse GitHub API response: %w", err)
//     }

//     content, ok := result["content"].(map[string]interface{})
//     if !ok || content["download_url"] == nil {
//         return "", fmt.Errorf("GitHub API response missing download_url")
//     }

//     return content["download_url"].(string), nil
// }

func UploadImageToGitHub(c *fiber.Ctx, file *multipart.FileHeader, productName string) (string, error) {
	// 🔹 Ambil token dari Header Authorization
	githubToken := c.Get("Authorization")

	// ✅ Cek apakah token kosong
	if githubToken == "" {
		return "", fmt.Errorf("GitHub token is missing. Please provide a valid token")
	}

	// 🔹 Debugging untuk memastikan token terbaca
	fmt.Println("GitHub Token:", githubToken)

	// 🔹 Sesuaikan username dan nama repository GitHub
	repoOwner := "serlip06"
	repoName := "Gambar"

	// 🔹 Gunakan nama folder yang sesuai dengan GitHub
	folderName := "bukti_pembayaran" // ✅ Pastikan folder ini benar di GitHub
	filePath := fmt.Sprintf("%s/%d_%s.jpg", folderName, time.Now().Unix(), productName)

	// ✅ Debugging path
	fmt.Println("Uploading file to:", filePath)

	// 🔹 Buka file untuk dibaca
	fileContent, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %w", err)
	}
	defer fileContent.Close()

	// 🔹 Baca isi file
	imageData, err := io.ReadAll(fileContent)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	// 🔹 Encode file ke base64
	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	// 🔹 Siapkan payload untuk GitHub API
	payload := map[string]string{
		"message": fmt.Sprintf("Add image for product %s", productName),
		"content": encodedImage,
	}

	payloadBytes, _ := json.Marshal(payload)

	// 🔹 Buat request ke GitHub API
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, filePath)
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(payloadBytes))

	// ✅ Perbaikan: Gunakan format Authorization yang benar
	req.Header.Set("Authorization", "token "+githubToken) // ✅ Gunakan "token " bukan "Bearer "
	req.Header.Set("Content-Type", "application/json")

	// 🔹 Kirim request ke GitHub API
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to GitHub: %w", err)
	}
	defer resp.Body.Close()

	// 🔹 Jika terjadi error, tampilkan respon API
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("GitHub API Response:", string(body)) // ✅ Debugging untuk melihat respon API GitHub
		return "", fmt.Errorf("GitHub API error: %s", body)
	}

	// 🔹 Parse response JSON dari GitHub API
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse GitHub API response: %w", err)
	}

	// 🔹 Ambil download URL gambar yang sudah diupload
	content, ok := result["content"].(map[string]interface{})
	if !ok || content["download_url"] == nil {
		return "", fmt.Errorf("GitHub API response missing download_url")
	}

	return content["download_url"].(string), nil
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

	// Proses upload gambar bukti pembayaran jika ada
	if file, err := c.FormFile("bukti_pembayaran"); err == nil { // ✅ Pastikan menangani dua nilai yang dikembalikan
		// ✅ Tambahkan `c` sebagai parameter pertama saat memanggil UploadImageToGitHub
		imageURL, err := UploadImageToGitHub(c, file, "bukti_pembayaran")
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "Failed to upload payment proof: " + err.Error(),
			})
		}
		input.BuktiPembayaran = imageURL // ✅ Update dengan URL gambar baru
	}

	// Update transaksi dengan data terbaru
	update := bson.M{"$set": bson.M{
		"id_user":           input.IDUser,
		"id_cartitem":       input.IDCartItem,
		"total_harga":       calculatedTotal,
		"metode_pembayaran": input.MetodePembayaran,
		"bukti_pembayaran":  input.BuktiPembayaran, // ✅ Menyimpan URL gambar bukti pembayaran yang diupdate
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
