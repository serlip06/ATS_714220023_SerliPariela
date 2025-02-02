package url

import (
	"github.com/serlip06/ATS_714220023_SerliPariela/controller"
	//"go.mongodb.org/mongo-driver/mongo"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger" // swagger handler
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
	// page.Get("/pelanggan", controller.GetPelanggan)	
	page.Get("/customer", controller.GetCustomer)//memanggil customer
	page.Get("/customer/:id", controller.GetCustomerID)//memanggil data customer berdasarkan id
	page.Post("/insert", controller.InsertDataCustomer)//post data customer (insert data)	
	page.Put("/update/:id", controller.UpdateData)//update data 
	page.Delete("/delete/:id", controller.DeleteCustomerByID)

	//endpoint bagian produk 
	page.Get("/produk", controller.GetProduks)//memanggil data produk
	page.Get("/produk/makanan", func(c *fiber.Ctx) error { return controller.GetAllProduksByKategori(c, "Makanan") })//produk kategori berdasarkan makanan 
	page.Get("/produk/minuman", func(c *fiber.Ctx) error { return controller.GetAllProduksByKategori(c, "Minuman")  }) // Produk kategori Minuman
	page.Get("/produk/:id", controller.GetProduksID)//memangil data berdasarkan id 
	page.Post("/insertproduk", controller.InsertDataProduk)//insert data produk
	page.Put("/updateproduk/:id", controller.UpdateDataProduk)//update data produk
	page.Delete("/deleteproduk/:id", controller.DeleteProduksByID)

	// endpoint bagian chartitem
	page.Get("/chartitem", controller.GetCartItem)//memanggil data chart item
	page.Get("/chartitem/:id", controller.GetCartItemID)//memanggil data chart item by id
	page.Post("/insertchartitem", controller.InsertDataCartItem)//insert chart item
	page.Put("/updatechartitem/:id", controller.UpdateDataCartItem)//update chart item
	page.Delete("/deletechartitem/:id", controller.DeleteCartItemByID)//delete chart item
	page.Post("/checkout", controller.CheckoutFromCart) // Checkout cart item
	page.Put("/cart-items/:id/select", controller.UpdateCartItemSelection)



	//login register untuk user 
	page.Post("/register", controller.RegisterHandler)
	page.Post("/login", controller.LoginHandler)
	page.Post("/approve-regis/:id", controller.ApproveRegistrationHandler)
	// get data user 
	page.Get("/user", controller.GetAllUsers) 
	page.Get("/user/:id", controller.GetUserByID)
	page.Get("/pendingregis", controller.GetAllPendingRegistrations)
	//notifikasi 
	page.Post("/update-product", controller.UpdateProductHandler)  // Misalnya menggunakan metode PUT
	page.Post("/add-product", controller.AddProductHandler)

	//transaksi
	page.Get("/transaksi", controller.GetAllTransaksi)
	page.Get("/transaksi/:id", controller.GetTransaksiByID)
	page.Get("/transaksi-user/:id", controller.GetTransaksiByUserID)
	page.Post("/inserttransaksi", controller.InsertTransaksi)
	page.Put("/updatetransaksi/:id", controller.UpdateTransaksi)
	page.Delete("/deletetransaksi/:id", controller.DeleteTransaksiByID)
	//swager
	page.Get("/docs/*", swagger.HandlerDefault)
}
