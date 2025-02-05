package controller

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Pelanggan struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama         string             `bson:"nama,omitempty" json:"nama,omitempty"`
	Phone_number string             `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Alamat      string             `bson:"alamat,omitempty" json:"alamat,omitempty"`
	Email   []string           `bson:"email,omitempty" json:"email,omitempty"`
}

type Produk struct {
	IDProduk    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" example:"1234567"`
	Nama_Produk string             `bson:"nama_produk,omitempty" json:"nama_produk,omitempty"  example:"nama makanan/minuman"`
	Deskripsi   string             `bson:"deskripsi,omitempty" json:"deskripsi,omitempty" example:"minuman teh manis yang menyegarkan"`
	Harga       int                `bson:"harga,omitempty" json:"harga,omitempty" example:"10000"`
	Gambar      string             `bson:"gambar,omitempty" json:"gambar,omitempty" example:"https://i.pinimg.com/564x/94/82/ab/9482ab2e248d249e7daa7fd6924c8d3b.jpg" `
	Stok        int                `bson:"stok,omitempty" json:"stok,omitempty" example:"5" `
	Kategori    string             `bson:"kategori,omitempty" json:"kategori,omitempty"  example:"makanan"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type ReqProduk struct {
	Nama_Produk string             `bson:"nama_produk,omitempty" json:"nama_produk,omitempty"  example:"test swagger"`
	Deskripsi   string             `bson:"deskripsi,omitempty" json:"deskripsi,omitempty" example:"minuman teh manis yang menyegarkan"`
	Harga       int                `bson:"harga,omitempty" json:"harga,omitempty" example:"10000"`
	Gambar      string             `bson:"gambar,omitempty" json:"gambar,omitempty" example:"https://i.pinimg.com/564x/94/82/ab/9482ab2e248d249e7daa7fd6924c8d3b.jpg" `
	Stok        int                `bson:"stok,omitempty" json:"stok,omitempty" example:"5" `
	Kategori    string             `bson:"kategori,omitempty" json:"kategori,omitempty"  example:"makanan"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type Transaksi struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Metode_Pembayaran    string          `bson:"metode_pembayaran,omitempty" json:"metode_pembayaran,omitempty"`
	Tanggal_Waktu     string            `bson:"tanggal_waktu,omitempty" json:"tanggal_waktu,omitempty"`
}

type Customer struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" example:"1234567"`
	Nama         string             `bson:"nama,omitempty" json:"nama,omitempty" example:"xavieraa putri"`
	Phone_number string             `bson:"phone_number,omitempty" json:"phone_number,omitempty" example:"085798654096"`
	Alamat       string             `bson:"alamat,omitempty" json:"alamat,omitempty" example:"jl.sarijadi"`
	Email        []string           `bson:"email,omitempty" json:"email" example:"Xaviera_89@gmail.com,Putri_90@gmail.com"`
	Nama_Produk  string             `bson:"nama_produk,omitempty" json:"nama_produk,omitempty" example:"Nasi Goreng"`
	Deskripsi    string             `bson:"deskripsi,omitempty" json:"deskripsi,omitempty" example:"nasi goreng dengan telor dan daging"`
	Harga        int                `bson:"harga,omitempty" json:"harga,omitempty" example:"15000"`
	Gambar       string             `bson:"gambar,omitempty" json:"gambar,omitempty" example:"https://i.pinimg.com/564x/94/82/ab/9482ab2e248d249e7daa7fd6924c8d3b.jpg"`
	Stok         string             `bson:"stok,omitempty" json:"stok,omitempty" example:"10"`
}

type ReqCustomer struct{
	Nama         string             `bson:"nama,omitempty" json:"nama,omitempty" example:"Tes swager"`
	Phone_number string             `bson:"phone_number,omitempty" json:"phone_number,omitempty" example:"085798654096"`
	Alamat       string             `bson:"alamat,omitempty" json:"alamat,omitempty" example:"jl.sarijadi"`
	Email        []string           `bson:"email,omitempty" json:"email" example:"Xaviera_89@gmail.com,Putri_90@gmail.com"`
	Nama_Produk  string             `bson:"nama_produk,omitempty" json:"nama_produk,omitempty" example:"Nasi Goreng"`
	Deskripsi    string             `bson:"deskripsi,omitempty" json:"deskripsi,omitempty" example:"nasi goreng dengan telor dan daging"`
	Harga        int                `bson:"harga,omitempty" json:"harga,omitempty" example:"15000"`
	Gambar       string             `bson:"gambar,omitempty" json:"gambar,omitempty" example:"https://i.pinimg.com/564x/94/82/ab/9482ab2e248d249e7daa7fd6924c8d3b.jpg"`
	Stok         string             `bson:"stok,omitempty" json:"stok,omitempty" example:"10"`
}

type CartItem struct {
	IDCartItem  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" example:"1234567"`             // ID unik untuk item keranjang
	IDProduk    primitive.ObjectID `bson:"id_produk,omitempty" json:"id_produk,omitempty" example:"1234567"` // Referensi ke ID Produk
	IDUser      primitive.ObjectID `bson:"id_user" json:"id_user" example:"1234567"`
	Nama_Produk string             `bson:"nama_produk,omitempty" json:"nama_produk,omitempty" example:"ikan bakar"` //nama untuk produknya
	Harga       int                `bson:"harga,omitempty" json:"harga,omitempty" example:"5000"`             // Harga produk pada saat dimasukkan ke keranjang
	Quantity    int                `bson:"quantity,omitempty" json:"quantity,omitempty" example:"1"`       // Jumlah produk dalam keranjang
	SubTotal    int                `bson:"sub_total,omitempty" json:"sub_total,omitempty" example:"2000"`     // Total harga (Harga * Quantity)
	Gambar      string             `bson:"gambar,omitempty" json:"gambar,omitempty"  example:"https://i.pinimg.com/564x/94/82/ab/9482ab2e248d249e7daa7fd6924c8d3b.jpg" `           // Gambar produk
	IsSelected  bool               `bson:"is_selected,omitempty" json:"is_selected,omitempty"  example:"true"` // Tambahkan flag ini
	
}
type ReqCartItem struct {
	IDProduk    primitive.ObjectID `bson:"id_produk,omitempty" json:"id_produk,omitempty" example:"1234567"` // Referensi ke ID Produk
	IDUser      primitive.ObjectID `bson:"id_user" json:"id_user" example:"1234567"`
	Nama_Produk string             `bson:"nama_produk,omitempty" json:"nama_produk,omitempty" example:"ikan bakar"` //nama untuk produknya
	Harga       int                `bson:"harga,omitempty" json:"harga,omitempty" example:"5000"`             // Harga produk pada saat dimasukkan ke keranjang
	Quantity    int                `bson:"quantity,omitempty" json:"quantity,omitempty" example:"1"`       // Jumlah produk dalam keranjang
	SubTotal    int                `bson:"sub_total,omitempty" json:"sub_total,omitempty" example:"2000"`     // Total harga (Harga * Quantity)
	Gambar      string             `bson:"gambar,omitempty" json:"gambar,omitempty"  example:"https://i.pinimg.com/564x/94/82/ab/9482ab2e248d249e7daa7fd6924c8d3b.jpg" `           // Gambar produk
	IsSelected  bool               `bson:"is_selected,omitempty" json:"is_selected,omitempty"  example:"true"` // Tambahkan flag ini
	
}
