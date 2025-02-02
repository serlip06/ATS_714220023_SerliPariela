package controller

import (
	"errors"
	"time"
	"fmt"
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
	"github.com/serlip06/ATS_714220023_SerliPariela/config"
	inimodel "github.com/serlip06/pointsalesofkantin/model"
	cek "github.com/serlip06/pointsalesofkantin/module"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

// func GetPelanggan(c *fiber.Ctx) error {
// 	ps := cek.GetAllPelanggan()
// 	return c.JSON(ps)
// }

//produk 
//get produk untuk all 
func GetProduks(c *fiber.Ctx) error {
	// Ambil query "kategori" dari URL, default kosong jika tidak diberikan
	kategori := c.Query("kategori", "")

	// Panggil fungsi GetAllProduks dari module
	produks, err := cek.GetAllProduks(kategori)
	if err != nil {
		// Jika ada error, kembalikan response dengan status 500
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error retrieving products: %v", err),
		})
	}

	// Tambahkan pengecekan dan pengisian CreatedAt jika kosong
	for i := range produks {
		if produks[i].CreatedAt.IsZero() {
			produks[i].CreatedAt = time.Now()
		}
	}

	// Jika sukses, kembalikan response dengan status 200
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Data produk berhasil diambil",
		"data":    produks,
	})
}


//customer 
//get data customer 

// GetPresensi godoc
// @Summary Get All Data Customer.
// @Description Mengambil semua data customer.
// @Tags Customer
// @Accept json
// @Produce json
// @Success 200 {object} Customer
// @Router /customer [get]
func GetCustomer(c *fiber.Ctx) error{
	ps := cek.GetAllCustomer()
	return c.JSON(ps)
}
// memanggil id 
// func GetPelangganByID(c *fiber.Ctx) {
// 	pelangganID, err := primitive.ObjectIDFromHex("615af14ae62f4c488e1d6d14")
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return
// 	}
// 	pelanggan := cek.GetPelangganByID(pelangganID)
// 	fmt.Println(pelanggan)
// }

//customer

//memanggil id customer 
// GetCustomerID godoc
// @Summary Get By ID Data Customer.
// @Description Ambil per ID data customer.
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path string true "Masukan ID"
// @Success 200 {object} Customer
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /customer/{id} [get]
func GetCustomerID(c *fiber.Ctx) error {
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
	ps, err := cek.GetCustomerFromID(objID, config.Ulbimongoconn, "customer")
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

//insert data customer 
// InsertDataCustomer godoc
// @Summary Insert data customer.
// @Description Input data customer.
// @Tags Customer
// @Accept json
// @Produce json
// @Param request body ReqCustomer true "Payload Body [RAW]"
// @Success 200 {object} Customer
// @Failure 400
// @Failure 500
// @Router /insert [post]
func InsertDataCustomer(c *fiber.Ctx) error {
	//db := config.Ulbimongoconn
	var customer inimodel.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	insertedID := cek.InsertCustomer(
		customer.Nama,
		customer.Phone_number,
		customer.Alamat,
		customer.Email,
	)

	if insertedID == "" { // Assuming an empty string means an error occurred
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Error inserting customer data",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Data berhasil disimpan.",
		"inserted_id": insertedID,
	})
}

// update data 
// UpdateData godoc
// @Summary Update data customer.
// @Description Ubah data customer.
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path string true "Masukan ID"
// @Param request body ReqCustomer true "Payload Body [RAW]"
// @Success 200 {object} Customer
// @Failure 400
// @Failure 500
// @Router /update/{id} [put]
func UpdateData(c *fiber.Ctx) error {
	db := config.Ulbimongoconn

	// Get the ID from the URL parameter
	id := c.Params("id")

	// Parse the ID into an ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Parse the request body into a Presensi object
	var customer inimodel.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Call the Updatecustomer function with the parsed ID and the Presensi object
	err = cek.UpdateCustomer(db, "customer",
		objectID,
		customer.Nama,
		customer.Phone_number,
		customer.Alamat,
		customer.Email,
		)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Data successfully updated",
	})
}

//delete data 
// DeleteCustomerByID godoc
// @Summary Delete data customer.
// @Description Hapus data customer.
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path string true "Masukan ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /delete/{id} [delete]
func DeleteCustomerByID(c *fiber.Ctx) error {
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

	err = cek.DeleteCustomerByID(objID, config.Ulbimongoconn, "customer")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error deleting data for id %s", id),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Data with id %s deleted successfully", id),
	})
}

//Produk
//function untuk mengambil data produk by ID
func GetProduksID(c *fiber.Ctx) error {
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
	ps, err := cek.GetProduksFromID(objID, config.Ulbimongoconn, "produk")
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
// function untuk get all pruduk 
// catatan ini masi ga pake filter ya ntar kalo berhail dia manual nambahin jalurnya sendiri di linknya 
func GetAllProduks(c *fiber.Ctx) error {
	// Ambil parameter query "kategori"
	kategori := c.Query("kategori", "") // Default-nya kosong jika tidak diisi

	// Panggil fungsi GetAllProduks dari backend
	produks, err := cek.GetAllProduks(kategori)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error retrieving products: %v", err),
		})
	}

	// Periksa dan tambahkan properti CreatedAt jika kosong
	for i := range produks {
		if produks[i].CreatedAt.IsZero() {
			produks[i].CreatedAt = time.Now()
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Data produk berhasil diambil",
		"data":    produks,
	})
}


// fitur tambahan kalo mau endpointnya nambahin produk/makanan atau produk/minuman 
func GetAllProduksByKategori(c *fiber.Ctx, kategori string) error {
    // Pastikan kategori valid
    if kategori == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Kategori parameter is required",
        })
    }

    produks, err := cek.GetAllProduks(kategori)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": fmt.Sprintf("Error retrieving products for kategori '%s': %v", kategori, err),
        })
    }

    if len(produks) == 0 {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{
            "status":  http.StatusNotFound,
            "message": fmt.Sprintf("No products found for kategori '%s'", kategori),
        })
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "status":  http.StatusOK,
        "message": fmt.Sprintf("Data produk kategori '%s' berhasil diambil", kategori),
        "data":    produks,
    })
}

//insert data produk
//ini diperbaharui karena ada kategori 
func InsertDataProduk(c *fiber.Ctx) error {
	var produk inimodel.Produk

	// Parse body
	if err := c.BodyParser(&produk); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Validasi kategori
	if produk.Kategori != "Makanan" && produk.Kategori != "Minuman" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Kategori harus berupa 'Makanan' atau 'Minuman'",
		})
	}

	// Insert data produk ke backend
	insertedID, err := cek.InsertDataProduk(
		produk.Nama_Produk,
		produk.Deskripsi,
		produk.Harga,
		produk.Gambar,
		produk.Stok,
		produk.Kategori,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Set properti CreatedAt dengan waktu saat ini
	produk.CreatedAt = time.Now()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Data produk berhasil disimpan",
		"inserted_id": insertedID,
		"created_at":  produk.CreatedAt,
	})
}


//update data produk 
// ini juga di perbaharui 
func UpdateDataProduk(c *fiber.Ctx) error {
	db := config.Ulbimongoconn

	// Get the ID from the URL parameter
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid id parameter",
		})
	}

	// Parse the request body into Produk object
	var produk inimodel.Produk
	if err := c.BodyParser(&produk); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Validasi kategori
	if produk.Kategori != "Makanan" && produk.Kategori != "Minuman" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Kategori harus berupa 'Makanan' atau 'Minuman'",
		})
	}

	// Update data produk
	err = cek.UpdateProduks(
		db, "produk", objectID,
		produk.Nama_Produk,
		produk.Deskripsi,
		produk.Harga,
		produk.Gambar,
		produk.Stok,
		produk.Kategori,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Data produk berhasil diperbarui",
	})
}

//delete data produk
func DeleteProduksByID(c *fiber.Ctx) error {
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

	err = cek.DeleteProduksByID(objID, config.Ulbimongoconn, "produk")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error deleting data for id %s", id),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Data with id %s deleted successfully", id),
	})
}

// cartitem bagian keranjang 

// get all cartitem 
func GetCartItem(c *fiber.Ctx) error {
	ps := cek.GetAllCartItems()
	return c.JSON(ps)
}

// get cartitemfromID
func GetCartItemID(c *fiber.Ctx) error {
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
	ps, err := cek.GetCartItemFromID(objID, config.Ulbimongoconn, "cart_items")
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

// InsertDataCartItem (Keranjang)
func InsertDataCartItem(c *fiber.Ctx) error {
	var input struct {
		IDProduk  primitive.ObjectID `json:"id_produk"`
		Quantity  int                `json:"quantity"`
		IDUser    primitive.ObjectID `json:"id_user"`
		IsSelected bool              `json:"is_selected"` // Tambahkan field ini
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	produk, err := cek.GetProduksFromID(input.IDProduk, config.Ulbimongoconn, "produk")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error retrieving product: %v", err),
		})
	}

	if produk.Stok < input.Quantity {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Stok tidak cukup untuk produk ini",
		})
	}

	subTotal := input.Quantity * produk.Harga
	cartItem := inimodel.CartItem{
		IDCartItem: primitive.NewObjectID(),
		IDProduk:   input.IDProduk,
		IDUser:     input.IDUser,
		Nama_Produk: produk.Nama_Produk,
		Harga:      produk.Harga,
		Quantity:   input.Quantity,
		SubTotal:   subTotal,
		Gambar:     produk.Gambar,
		IsSelected:  input.IsSelected, // Simpan nilai is_selected
	}

	collection := config.Ulbimongoconn.Collection("cart_items")
	_, err = collection.InsertOne(c.Context(), cartItem)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Data cart item berhasil disimpan",
		"cart_item":   cartItem,
	})
}

// funsi untuk checkout tambahan 
func UpdateCartItemSelection(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID format",
		})
	}

	var input struct {
		IsSelected bool `json:"is_selected"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Failed to parse request body",
		})
	}

	collection := config.Ulbimongoconn.Collection("cart_items")
	update := bson.M{"$set": bson.M{"is_selected": input.IsSelected}}
	_, err = collection.UpdateOne(c.Context(), bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to update selection",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Cart item selection updated successfully",
	})
}
// function untuk checkout 
func CheckoutFromCart(c *fiber.Ctx) error {
	var input struct {
		IDUser           primitive.ObjectID `json:"id_user"`
		MetodePembayaran string            `json:"metode_pembayaran"`
		BuktiPembayaran  string            `json:"bukti_pembayaran"`
		Alamat           string            `json:"alamat"`
	}

	// Perbaikan: Status error menjadi 400 (Bad Request)
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
		})
	}

	// Periksa apakah `IDUser` valid (jangan kosong)
	if input.IDUser.IsZero() {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "IDUser is required",
		})
	}

	// Pastikan pemanggilan module sesuai dengan struktur kode kamu
	transaksiID, err := cek.CheckoutFromCart(config.Ulbimongoconn, input.IDUser, input.MetodePembayaran, input.BuktiPembayaran, input.Alamat)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Checkout failed",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":         http.StatusOK,
		"message":        "Checkout successful",
		"transaction_id": transaksiID,
	})
}


//update data keranjang
func UpdateDataCartItem(c *fiber.Ctx) error {
    // Parse ID dari parameter URL
    id := c.Params("id")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Invalid ID format",
        })
    }

    // Parse body hanya untuk Quantity
    var input struct {
        Quantity int `json:"quantity"`
    }
    if err := c.BodyParser(&input); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Failed to parse request body",
        })
    }

    // Validasi Quantity
    if input.Quantity <= 0 {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Quantity must be greater than 0",
        })
    }

    // Ambil data CartItem berdasarkan ID
    collection := config.Ulbimongoconn.Collection("cart_items")
    var existingCartItem inimodel.CartItem
    err = collection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&existingCartItem)
    if err != nil {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{
            "status":  http.StatusNotFound,
            "message": "Cart item not found",
        })
    }

    // Hitung ulang SubTotal berdasarkan Quantity baru dan Harga asli
    existingCartItem.Quantity = input.Quantity
    existingCartItem.SubTotal = existingCartItem.Quantity * existingCartItem.Harga

    // Update data di database
    update := bson.M{
        "$set": bson.M{
            "quantity":  existingCartItem.Quantity,
            "sub_total": existingCartItem.SubTotal,
        },
    }
    _, err = collection.UpdateOne(c.Context(), bson.M{"_id": objectID}, update)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": "Failed to update cart item",
        })
    }

    // Berikan respon sukses
    return c.Status(http.StatusOK).JSON(fiber.Map{
        "status":      http.StatusOK,
        "message":     "Cart item successfully updated",
        "updated_data": existingCartItem, // Kembalikan data terbaru
    })
}

// hapus dari keranjang
func DeleteCartItemByID(c *fiber.Ctx) error {
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

	// Hapus data dari database tanpa validasi IDUser
	collection := config.Ulbimongoconn.Collection("cart_items")
	result, err := collection.DeleteOne(c.Context(), bson.M{"_id": objID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error deleting data for id %s", id),
		})
	}

	// Periksa apakah ada dokumen yang dihapus
	if result.DeletedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": "Cart item not found",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Data with id %s deleted successfully", id),
	})
}


// cadangan untuk cartitem

// func InsertDataCartItem(c *fiber.Ctx) error {
// 	var input struct {
// 		IDProduk  primitive.ObjectID `json:"id_produk"`  // ID Produk sebagai input
// 		Quantity  int                `json:"quantity"`   // Jumlah produk
// 		IDUser    primitive.ObjectID `json:"id_user"`    // ID User sebagai input baru
// 	}

// 	// Parse body
// 	if err := c.BodyParser(&input); err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": err.Error(),
// 		})
// 	}

// 	// Validasi input
// 	if input.Quantity <= 0 {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"status":  http.StatusBadRequest,
// 			"message": "Quantity harus lebih besar dari 0",
// 		})
// 	}

// 	// Ambil data produk berdasarkan IDProduk
// 	produk, err := cek.GetProduksFromID(input.IDProduk, config.Ulbimongoconn, "produk")
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": fmt.Sprintf("Error retrieving product: %v", err),
// 		})
// 	}

// 	// Validasi stok
// 	if produk.Stok < input.Quantity {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"status":  http.StatusBadRequest,
// 			"message": "Stok tidak cukup untuk produk ini",
// 		})
// 	}

// 	// Hitung subtotal
// 	subTotal := input.Quantity * produk.Harga

// 	// Buat data cart item
// 	cartItem := inimodel.CartItem{
// 		IDCartItem: primitive.NewObjectID(), // Buat ID baru untuk CartItem
// 		IDProduk:   input.IDProduk,
// 		IDUser:     input.IDUser,            // Tambahkan IDUser ke dalam data
// 		Nama_Produk: produk.Nama_Produk,
// 		Harga:      produk.Harga,
// 		Quantity:   input.Quantity,
// 		SubTotal:   subTotal,
// 		Gambar:     produk.Gambar,
// 	}

// 	// Insert data ke MongoDB
// 	collection := config.Ulbimongoconn.Collection("cart_items")
// 	_, err = collection.InsertOne(c.Context(), cartItem)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": err.Error(),
// 		})
// 	}

// 	return c.Status(http.StatusOK).JSON(fiber.Map{
// 		"status":      http.StatusOK,
// 		"message":     "Data cart item berhasil disimpan",
// 		"cart_item":   cartItem,
// 	})
// }










// func GetPelangganByID(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	if id == "" {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": "Wrong parameter",
// 		})
// 	}
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"status":  http.StatusBadRequest,
// 			"message": "Invalid id parameter",
// 		})
// 	}
// 	ps, err := cek.GetPelangganFromID(objID, config.Ulbimongoconn, "kantin_pelangan")
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return c.Status(http.StatusNotFound).JSON(fiber.Map{
// 				"status":  http.StatusNotFound,
// 				"message": fmt.Sprintf("No data found for id %s", id),
// 			})
// 		}
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": fmt.Sprintf("Error retrieving data for id %s", id),
// 		})
// 	}
// 	return c.JSON(ps)
// }
// func GetAllPelanggan(c *fiber.Ctx) error {
// 	ps := GetAllPelanggan(config.Ulbimongoconn, "kantin_pelanggan")
// 	return c.JSON(ps)
// }

