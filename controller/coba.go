package controller

import (
	//"errors"
	"fmt"
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
	//"github.com/serlip06/ATS_714220023_SerliPariela/config"
	cek "github.com/serlip06/pointsalesofkantin/module"
	//"github.com/serlip06/ATS_714220023_SerliPariela/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
	//"net/http"
)

func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

func GetPelanggan(c *fiber.Ctx) error {
	ps := cek.GetAllPelanggan()
	return c.JSON(ps)
}

// memanggil id 
func GetPelangganByID(c *fiber.Ctx) {
	pelangganID, err := primitive.ObjectIDFromHex("615af14ae62f4c488e1d6d14")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	pelanggan := cek.GetPelangganByID(pelangganID)
	fmt.Println(pelanggan)
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

