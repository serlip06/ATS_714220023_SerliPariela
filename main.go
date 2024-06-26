package main

import (
	"log"

	"github.com/serlip06/ATS_714220023_SerliPariela/config"
	//"github.com/serlip06/pointsalesofkantin/module"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/serlip06/ATS_714220023_SerliPariela/url"
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
)

func main() {
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.Web(site)
	log.Fatal(site.Listen(musik.Dangdut()))
}