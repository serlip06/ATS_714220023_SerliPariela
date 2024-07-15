package config

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

var origins = []string{
	"https://auth.ulbi.ac.id",
	"https://sip.ulbi.ac.id",
	"https://euis.ulbi.ac.id",
	"https://home.ulbi.ac.id",
	"https://alpha.ulbi.ac.id",
	"https://dias.ulbi.ac.id",
	"https://iteung.ulbi.ac.id",
	"https://whatsauth.github.io",
	"http://127.0.0.1:5500",
	"http://127.0.0.1:5502",
	"http://127.0.0.1:5503",
	"http://127.0.0.1:40013",
	"https://serlip06.github.io",
	"http://127.0.0.1:8080",
}

var Internalhost string = os.Getenv("INTERNALHOST") + ":" + os.Getenv("PORT")

var Cors = cors.Config{
	AllowOrigins:     strings.Join(origins[:], ","),
	AllowMethods: 	  "GET, HEAD, OPTIONS, POST, PUT, DETELE",
	AllowHeaders:     "Origin,Login, Content-Type",
	ExposeHeaders:    "Content-Length",
	AllowCredentials: true,
}
