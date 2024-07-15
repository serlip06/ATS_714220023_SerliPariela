package main

import (
	"log"

	"github.com/serlip06/ATS_714220023_SerliPariela/config"
	//"github.com/serlip06/pointsalesofkantin/module"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/serlip06/ATS_714220023_SerliPariela/url"
	"github.com/aiteung/musik"
	_ "github.com/serlip06/ATS_714220023_SerliPariela/docs" 
	"github.com/gofiber/fiber/v2"
)
// @title TES SWAGGER ULBI
// @version 1.0
// @description This is a sample swagger for Fiber

// @contact.name API Support
// @contact.url https://github.com/serlip06
// @contact.email 714220023@std.ulbi.ac.id

// @host ats-714220023-serlipariela-38bba14820aa.herokuapp.com
// @BasePath /
// @schemes https http
func main() {
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.Web(site)
	log.Fatal(site.Listen(musik.Dangdut()))
}