package controller

import (

	//"fmt"
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
	//"github.com/serlip06/ATS_714220023_SerliPariela/config"
	cek "github.com/serlip06/pointsalesofkantin/module"
	//"github.com/serlip06/ATS_714220023_SerliPariela/config"
	//"go.mongodb.org/mongo-driver/bson/primitive"
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
// func GetAllPelanggan(c *fiber.Ctx) error {
// 	ps := GetAllPelanggan(config.Ulbimongoconn, "kantin_pelanggan")
// 	return c.JSON(ps)
// }

