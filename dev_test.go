// package HealHero

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/aiteung/atdb"
// 	"github.com/whatsauth/watoken"
// 	"go.mongodb.org/mongo-driver/bson"
// )

// // func TestUpdateGetData(t *testing.T) {
// // 	mconn := SetConnection("MONGOSTRING", "healhero_db")
// // 	datagedung := GetAllBangunanLineString(mconn, "healhero_db")
// // 	fmt.Println(datagedung)
// // }

// func TestGeneratePasswordHash(t *testing.T) {
// 	password := "cantik"
// 	hash, _ := HashPassword(password) // ignore error for the sake of simplicity
// 	fmt.Println("Password:", password)
// 	fmt.Println("Hash:    ", hash)

// 	match := CheckPasswordHash(password, hash)
// 	fmt.Println("Match:   ", match)
// }
// func TestGeneratePrivateKeyPaseto(t *testing.T) {
// 	privateKey, publicKey := watoken.GenerateKey()
// 	fmt.Println("private Key")
// 	fmt.Println(privateKey)
// 	fmt.Println("public Key")
// 	fmt.Println(publicKey)
// 	hasil, err := watoken.Encode("rizkyria", privateKey)
// 	fmt.Println(hasil, err)
// }

// func TestHashFunction(t *testing.T) {
// 	mconn := SetConnection("MONGOSTRING", "healhero_db")
// 	var userdata User
// 	userdata.Username = "rizkyria"
// 	userdata.Password = "cantik"

// 	filter := bson.M{"username": userdata.Username}
// 	res := atdb.GetOneDoc[User](mconn, "user", filter)
// 	fmt.Println("Mongo User Result: ", res)
// 	hash, _ := HashPassword(userdata.Password)
// 	fmt.Println("Hash Password : ", hash)
// 	match := CheckPasswordHash(userdata.Password, res.Password)
// 	fmt.Println("Match:   ", match)

// }

// func TestIsPasswordValid(t *testing.T) {
// 	mconn := SetConnection("MONGOSTRING", "healhero_db")
// 	var userdata User
// 	userdata.Username = "rizkyria"
// 	userdata.Password = "cantik"

// 	anu := IsPasswordValid(mconn, "user", userdata)
// 	fmt.Println(anu)
// }

// func TestInsertUser(t *testing.T) {
// 	mconn := SetConnection("MONGOSTRING", "healhero_db")
// 	var userdata User
// 	userdata.Username = "rizkyria"
// 	userdata.Password = "cantik"

// 	nama := InsertUser(mconn, "user", userdata)
// 	fmt.Println(nama)
// }

package HealHero

import (
	"fmt"
	"testing"

	"github.com/HealHeroo/be_healhero/model"
	"github.com/HealHeroo/be_healhero/module"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

var db = module.MongoConnect("MONGOSTRING", "healhero_db")

func TestGetUserFromEmail(t *testing.T) {
	email := "admin@gmail.com"
	hasil, err := module.GetUserFromEmail(email, db)
	if err != nil {
		t.Errorf("Error TestGetUserFromEmail: %v", err)
	} else {
		fmt.Println(hasil)
	}
}

func TestInsertOneObat(t *testing.T) {
	var doc model.Obat
   doc.NamaObat= "Paracetamol"
   doc.JenisObat = "Analgesik dan Antipiretik"
   doc.Keterangan = "500 mg"
   doc.Harga = "RP 8.000"
   if  doc.NamaObat == "" || doc.JenisObat == "" || doc.Keterangan == "" || doc.Harga == ""   {
	   t.Errorf("mohon untuk melengkapi data")
   } else {
	   insertedID, err := module.InsertOneDoc(db, "obat", doc)
	   if err != nil {
		   t.Errorf("Error inserting document: %v", err)
		   fmt.Println("Data tidak berhasil disimpan")
	   } else {
	   fmt.Println("Data berhasil disimpan dengan id :", insertedID.Hex())
	   }
   }
}

type Userr struct {
	ID           	primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email  			string             `bson:"email,omitempty" json:"email,omitempty"`
	Role     		string			   `bson:"role,omitempty" json:"role,omitempty"`
}

func TestGetAllDoc(t *testing.T) {
	hasil := module.GetAllDocs(db, "user", []Userr{})
	fmt.Println(hasil)
}


func TestSignUpPengguna(t *testing.T) {
	var doc model.Pengguna
	doc.NamaLengkap = "Rizkyria Hutabarat"
	doc.TanggalLahir = "18/11/2003"
	doc.JenisKelamin = "Perempuan"
	doc.NomorHP = "081219882869"
	doc.Alamat = "Jalan Sarijadi No 56"
	doc.Akun.Email = "rizkyria@gmail.com"
	doc.Akun.Password = "rizkyriacantik"
	err := module.SignUpPengguna(db, doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
	fmt.Println("Data berhasil disimpan dengan nama :", doc.NamaLengkap)
	}
}

func TestSignUpDriver(t *testing.T) {
	var doc model.Driver
	doc.NamaLengkap = "Wawan Setiawan"
	doc.JenisKelamin = "Laki-laki"
	doc.NomorHP = "088475638475"
	doc.Alamat = "Jalan Jingga No 54"
	doc.PlatMotor = "D 8392 SDE"
	doc.Akun.Email = "wawan@gmail.com"
	doc.Akun.Password = "driverwawan"
	err := module.SignUpDriver(db, doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
	fmt.Println("Data berhasil disimpan dengan nama :", doc.NamaLengkap)
	}
}


func TestLogIn(t *testing.T) {
	var doc model.User
	doc.Email = "wawan@gmail.com"
	doc.Password = "driverwawan"
	user, err := module.LogIn(db, doc)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		fmt.Println("Selamat datang Driver:", user)
	}
}

func TestInsertOneOrder(t *testing.T) {
	var doc model.Order
	doc.NamaObat= "Amoxsilin"
   doc.Quantity= "1"
   doc.TotalCost = "Rp 7.000"
   doc.Status = "Pending"
   if  doc.Quantity == "" || doc.TotalCost == "" || doc.Status == ""    {
	   t.Errorf("mohon untuk melengkapi data")
   } else {
	   insertedID, err := module.InsertOneDoc(db, "order", doc)
	   if err != nil {
		   t.Errorf("Error inserting document: %v", err)
		   fmt.Println("Data tidak berhasil disimpan")
	   } else {
	   fmt.Println("Data berhasil disimpan dengan id :", insertedID.Hex())
	   }
   }
}


// func TestGeneratePrivateKeyPaseto(t *testing.T) {
// 	privateKey, publicKey := module.GenerateKey()
// 	fmt.Println("ini private key :", privateKey)
// 	fmt.Println("ini public key :", publicKey)
// 	id := "64d0b1104255ba95ba588512"
// 	objectId, err := primitive.ObjectIDFromHex(id)
// 	role := "admin"
// 	if err != nil{
// 		t.Fatalf("error converting id to objectID: %v", err)
// 	}
// 	hasil, err := module.Encode(objectId, role, privateKey)
// 	fmt.Println("ini hasil :", hasil, err)
// }

// func TestWatoken(t *testing.T) {
// 	body, err := module.Decode("f3248b509d9731ebd4e0ccddadb5a08742e083db01678e8a1d734ce81298868f", "v4.public.eyJlbWFpbCI6ImZheEBnbWFpbC5jb20iLCJleHAiOiIyMDIzLTEwLTIyVDAwOjQxOjQ1KzA3OjAwIiwiZmlyc3RuYW1lIjoiRmF0d2EiLCJpYXQiOiIyMDIzLTEwLTIxVDIyOjQxOjQ1KzA3OjAwIiwiaWQiOiI2NDkwNjNkM2FkNzJlMDc0Mjg2YzYxZTgiLCJsYXN0bmFtZSI6IkZhdGFoaWxsYWgiLCJuYmYiOiIyMDIzLTEwLTIxVDIyOjQxOjQ1KzA3OjAwIiwicm9sZSI6InBlbGFtYXIifR_Q4b9X7WC7up7dUUxz_Yki39M-ReovTIoTFfdJmFYRF5Oer0zQZx_ZQamkOsogJ6RuGJhxT3OxrXFS5p6dMg0")
// 	fmt.Println("isi : ", body, err)
// }

// <<<<<<< HEAD
// =======
// 	match := CheckPasswordHash(password, hash)
// 	fmt.Println("Match:   ", match)
// }
// func TestGeneratePrivateKeyPaseto(t *testing.T) {
// 	privateKey, publicKey := watoken.GenerateKey()
// 	fmt.Println("private Key")
// 	fmt.Println(privateKey)
// 	fmt.Println("publicKey")
// 	fmt.Println(publicKey)
// 	hasil, err := watoken.Encode("rizkyria", privateKey)
// 	fmt.Println(hasil, err)
// }
// >>>>>>> d4d8e3c7ef4ba1e17afe1ad98d8cc294de290774


// // test obat
// func TestInsertObat(t *testing.T) {
// 	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
// 	payload, err := module.Decode("b95509d9634ed137b5ccdd07a7534ab0dcede0f310c09634afbf0262c7a4ce1c", "v4.public.eyJleHAiOiIyMDIzLTEwLTMxVDA4OjQ4OjIyWiIsImlhdCI6IjIwMjMtMTAtMzFUMDY6NDg6MjJaIiwiaWQiOiI2NTQwNjMyODI4NzY0ZDk2YzY0OWYyOWQiLCJuYmYiOiIyMDIzLTEwLTMxVDA2OjQ4OjIyWiJ9lXy1b5nOEYuCn7_o-TcFuR-3OOm__T7SHlAdx3PQl4Du9EAr8pu85lvU6SddRar7YB3DEbf-zwfY_zytj7jrAQ")
// 	if err != nil {
// 		t.Errorf("Error decode token: %v", err)
// 	}
// 	// if payload.Role != "mitra" {
// 	// 	t.Errorf("Error role: %v", err)
// 	// }
// 	var dataobat model.Obat
// 	dataobat.NamaObat = "Paracetamol"
// 	dataobat.JenisObat = "Analgesik dan Antipiretik"
// 	dataobat.Keterangan = "500 Mg"
// 	dataobat.Harga = "Rp 8.000"
// 	err = module.InsertObat(payload.Id, conn, dataobat)
// 	if err != nil {
// 		t.Errorf("Error insert : %v", err)
// 	} else {
// 		fmt.Println("Success!!!")
// 	}
// }

// func TestUpdateObat(t *testing.T) {
// 	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
// 	payload, err := module.Decode("b95509d9634ed137b5ccdd07a7534ab0dcede0f310c09634afbf0262c7a4ce1c", "v4.public.eyJleHAiOiIyMDIzLTExLTAxVDA2OjQ5OjQ0WiIsImlhdCI6IjIwMjMtMTEtMDFUMDQ6NDk6NDRaIiwiaWQiOiI2NTQwNjMyODI4NzY0ZDk2YzY0OWYyOWQiLCJuYmYiOiIyMDIzLTExLTAxVDA0OjQ5OjQ0WiJ92RxBGslXaHBoLQhvMJLQO7uEBG5c5FmkpZgakPjmk1aUFDdRkw3m3r-7BpkhDmCtByoARDr36X3DhjCL8HT8AQ")
// 	if err != nil {
// 		t.Errorf("Error decode token: %v", err)
// 	}
// 	// if payload.Role != "mitra" {
// 	// 	t.Errorf("Error role: %v", err)
// 	// }
// 	var dataobat model.Obat
// 	dataobat.NamaObat = "Paracetamol"
// 	dataobat.JenisObat = "Analgesik dan Antipiretik"
// 	dataobat.Keterangan = "500 Mg"
// 	dataobat.Harga = "Rp 8.000"
// 	id := "65406377996edfaee3ed9a19"
// 	objectId, err := primitive.ObjectIDFromHex(id)
// 	if err != nil{
// 		t.Fatalf("error converting id to objectID: %v", err)
// 	}
// 	err = module.UpdateObat(objectId, payload.Id, conn, dataobat)
// 	if err != nil {
// 		t.Errorf("Error update : %v", err)
// 	} else {
// 		fmt.Println("Success!!!")
// 	}
// }

// func TestDeleteObat(t *testing.T) {
// 	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
// 	payload, err := module.Decode("b95509d9634ed137b5ccdd07a7534ab0dcede0f310c09634afbf0262c7a4ce1c", "v4.public.eyJleHAiOiIyMDIzLTExLTAxVDA2OjQ5OjQ0WiIsImlhdCI6IjIwMjMtMTEtMDFUMDQ6NDk6NDRaIiwiaWQiOiI2NTQwNjMyODI4NzY0ZDk2YzY0OWYyOWQiLCJuYmYiOiIyMDIzLTExLTAxVDA0OjQ5OjQ0WiJ92RxBGslXaHBoLQhvMJLQO7uEBG5c5FmkpZgakPjmk1aUFDdRkw3m3r-7BpkhDmCtByoARDr36X3DhjCL8HT8AQ")
// 	if err != nil {
// 		t.Errorf("Error decode token: %v", err)
// 	}
// 	// if payload.Role != "mitra" {
// 	// 	t.Errorf("Error role: %v", err)
// 	// }
// 	id := "65406377996edfaee3ed9a19"
// 	objectId, err := primitive.ObjectIDFromHex(id)
// 	if err != nil{
// 		t.Fatalf("error converting id to objectID: %v", err)
// 	}
// 	err = module.DeleteObat(objectId, payload.Id, conn)
// 	if err != nil {
// 		t.Errorf("Error delete : %v", err)
// 	} else {
// 		fmt.Println("Success!!!")
// 	}
// }



// func TestGetAllObat(t *testing.T) {
// 	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
// 	data, err := module.GetAllObat(conn)
// 	if err != nil {
// 		t.Errorf("Error get all : %v", err)
// 	} else {
// 		fmt.Println(data)
// 	}
// }

// func TestGetObatFromID(t *testing.T) {
// 	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
// 	id := "65406377996edfaee3ed9a19"
// 	objectId, err := primitive.ObjectIDFromHex(id)
// 	if err != nil{
// 		t.Fatalf("error converting id to objectID: %v", err)
// 	}
// 	obat, err := module.GetObatFromID(objectId, conn)
// 	if err != nil {
// 		t.Errorf("Error get obat : %v", err)
// 	} else {
// 		fmt.Println(obat)
// 	}
// }

