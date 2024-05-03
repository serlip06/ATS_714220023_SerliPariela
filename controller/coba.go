package controller

import (
	"github.com/serlip06/pointsalesofkantin"
	"github.com/gofiber/fiber/v2"
	
)

func Homepage(c *fiber.Ctx) error {
	ipaddr := _714220023.GetAllPelanggan()
	return c.JSON(ipaddr)
}
