package controller

import (
	"errors"
	"fmt"
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
	"github.com/serlip06/ATS_714220023_SerliPariela/config"
	inimodel "github.com/serlip06/pointsalesofkantin/model"
	cek "github.com/serlip06/pointsalesofkantin/module"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

func GetPelanggan(c *fiber.Ctx) error {
	ps := cek.GetAllPelanggan()
	return c.JSON(ps)
}

func GetProduks(c *fiber.Ctx) error {
	ps := cek.GetAllProduks()
	return c.JSON(ps)
}

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

//function untuk produk 
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
//insert data produk
func InsertDataProduk(c *fiber.Ctx) error {
	// db := config.Ulbimongoconn
	var produk inimodel.Produk
	if err := c.BodyParser(&produk); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	insertedID := cek.InsertDataProduk(
		produk.Nama_Produk,
		produk.Deskripsi,
		produk.Harga,
		produk.Gambar,
		produk.Stok,
	)

	if insertedID == "" { // Assuming an empty string means an error occurred
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Error inserting product data",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Data berhasil disimpan.",
		"inserted_id": insertedID,
	})
}
//update data produk 
func UpdateDataProduk(c *fiber.Ctx) error {
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
	var produk inimodel.Produk
	if err := c.BodyParser(&produk); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Call the Updatecustomer function with the parsed ID and the Presensi object
	err = cek.UpdateProduks(db, "produk",
		objectID,
		produk.Nama_Produk,
		produk.Deskripsi,
		produk.Harga,
		produk.Gambar,
		produk.Stok)
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

