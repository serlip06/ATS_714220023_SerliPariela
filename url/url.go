package url

import (
	"github.com/serlip06/ATS_714220023_SerliPariela/controller"
	
	"github.com/gofiber/fiber/v2"
)

func Web(page *fiber.App) {
	// page.Post("/api/whatsauth/request", controller.PostWhatsAuthRequest)  //API from user whatsapp message from iteung gowa
	// page.Get("/ws/whatsauth/qr", websocket.New(controller.WsWhatsAuthQR)) //websocket whatsauth
	page.Get("/", controller.Sink)
	page.Post("/", controller.Sink)
	page.Put("/", controller.Sink)
	page.Patch("/", controller.Sink)
	page.Delete("/", controller.Sink)
	page.Options("/", controller.Sink)

	page.Get("/checkip", controller.Homepage)
	page.Get("/pelanggan", controller.GetPelanggan)	
	page.Get("/customer", controller.GetCustomer)//memanggil customer
	page.Get("/customer/:id", controller.GetCustomerID)//memanggil data customer berdasarkan id
	page.Post("/insert", controller.InsertDataCustomer)//post data customer (insert data)	
	page.Put("/update/:id", controller.UpdateData)//update data 
	page.Delete("/delete/:id", controller.DeleteCustomerByID)
}
