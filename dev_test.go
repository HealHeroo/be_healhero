package HealHero

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/HealHeroo/be_healhero/model"
	"github.com/HealHeroo/be_healhero/module"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/argon2"
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

func TestInsertUser(t *testing.T) {
	var doc model.User
	doc.Email = "admin@gmail.com"
	password := "admin123"
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		t.Errorf("kesalahan server : salt")
	} else {
		hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
		user := bson.M{
			"email": doc.Email,
			"password": hex.EncodeToString(hashedPassword),
			"salt": hex.EncodeToString(salt),
			"role": "admin",
		}
		_, err = module.InsertOneDoc(db, "user", user)
		if err != nil {
			t.Errorf("gagal insert")
		} else {
			fmt.Println("berhasil insert")
		}
	}
}

func TestGetUserByAdmin(t *testing.T) {
	id := "655c3b9a1d6524f2f1200fc5"
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("Error converting id to objectID: %v", err)
	}
	data, err := module.GetUserFromID(idparam, db)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		if data.Role == "pengguna" {
			datapengguna, err := module.GetPenggunaFromAkun(data.ID, db)
			if err != nil {
				t.Errorf("Error getting document: %v", err)
			} else {
				datapengguna.Akun = data
				fmt.Println(datapengguna) 
			}
		}
		if data.Role == "driver" {
			datadriver, err := module.GetDriverFromAkun(data.ID, db)
			if err != nil {
				t.Errorf("Error getting document: %v", err)
			} else {
				datadriver.Akun = data
				fmt.Println(datadriver)
			}
		}
	}
}

func TestSignUpPengguna(t *testing.T) {
	var doc model.Pengguna
	doc.NamaLengkap = "Marlina"
	doc.TanggalLahir = "30/08/2003"
	doc.JenisKelamin = "Perempuan"
	doc.NomorHP = "081284739485"
	doc.Alamat = "Jalan Sarijadi No 56"
	doc.Akun.Email = "marlina@gmail.com"
	doc.Akun.Password = "marlinacantik"
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

func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := module.GenerateKey()
	fmt.Println("ini private key :", privateKey)
	fmt.Println("ini public key :", publicKey)
	id := "655c3b9a1d6524f2f1200fc6"
	objectId, err := primitive.ObjectIDFromHex(id)
	role := "pengguna"
	if err != nil{
		t.Fatalf("error converting id to objectID: %v", err)
	}
	hasil, err := module.Encode(objectId, role, privateKey)
	fmt.Println("ini hasil :", hasil, err)
}

func TestUpdatePengguna(t *testing.T) {
	var doc model.Pengguna
	id := "655c3b9a1d6524f2f1200fc8"
	objectId, _ := primitive.ObjectIDFromHex(id)
	id2 := "655c3b9a1d6524f2f1200fc6"
	userid, _ := primitive.ObjectIDFromHex(id2)
	doc.NamaLengkap = "Marlina Lubis"
	doc.TanggalLahir = "30/08/2003"
	doc.JenisKelamin = "Perempuan"
	doc.NomorHP = "081237629371"
	doc.Alamat = "Jalan Sarijadi No 53"
	if doc.NamaLengkap == "" || doc.TanggalLahir == "" || doc.JenisKelamin == "" || doc.NomorHP == "" || doc.Alamat == "" {
		t.Errorf("mohon untuk melengkapi data")
	} else {
		err := module.UpdatePengguna(objectId, userid, db, doc)
		if err != nil {
			t.Errorf("Error inserting document: %v", err)
			fmt.Println("Data tidak berhasil diupdate")
		} else {
			fmt.Println("Data berhasil diupdate")
		}
	}
}

func TestWatoken(t *testing.T) {
	body, err := module.Decode("fe58577f04c139838907cc8c298b6d0c6844aa7d14ef2e99d8b4d26f1b02ce01", "v4.public.eyJleHAiOiIyMDIzLTExLTI3VDE3OjExOjA5KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yN1QxNToxMTowOSswNzowMCIsImlkIjoiNjU1YzRjMWMzNTg0M2VkNTYzMWM5MDNkIiwibmJmIjoiMjAyMy0xMS0yN1QxNToxMTowOSswNzowMCIsInJvbGUiOiJkcml2ZXIifb4fz7rJ3dZV2qtQ1BGL19-pOEaU7evhdXP8910lkpGKM3dnWoKG0qxeZVObnk58hzaZAbQDKyBegF6-_R1rvwU")
	fmt.Println("isi : ", body, err)
}


func TestInsertOneOrder(t *testing.T) {
	var doc model.Order
	doc.NamaObat= "Vometa"
   doc.Quantity= "1"
   doc.TotalCost = "Rp 60.000"
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

// test obat
func TestInsertObat(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
	payload, err := module.Decode("2d2bdc3a1fca7cc064174e2a6e63e2b78f1db16de9e9ed42e63646709de4a1a", "v4.public.eyJleHAiOiIyMDIzLTExLTIxVDE2OjA2OjM4KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yMVQxNDowNjozOCswNzowMCIsImlkIjoiNjU1YzNiOWI0ZjZhNzVkZGFlZWNhMTkxIiwibmJmIjoiMjAyMy0xMS0yMVQxNDowNjozOCswNzowMCIsInJvbGUiOiJhZG1pbiJ9mGLShR1CooldqYp11ygx8dJt0UNUrj4XfIegnwhriKeZSfuv-9SOcr2XG5KKO1r0hL3_V8QFCev__cJgEaTzBA")
	if err != nil {
		t.Errorf("Error decode token: %v", err)
	}
	// if payload.Role != "mitra" {
	// 	t.Errorf("Error role: %v", err)
	// }
	var dataobat model.Obat
	dataobat.NamaObat = "Paracetamol"
	dataobat.JenisObat = "Analgesik dan Antipiretik"
	dataobat.Keterangan = "500 Mg"
	dataobat.Harga = "Rp 8.000"
	err = module.InsertObat(payload.Id, conn, dataobat)
	if err != nil {
		t.Errorf("Error insert : %v", err)
	} else {
		fmt.Println("Success!!!")
	}
}

// func TestUpdateObat(t *testing.T) {
// 	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
// 	payload, err := module.Decode("2d2bdc3a1fca7cc064174e2a6e63e2b78f1db16de9e9ed42e63646709de4a1a", "v4.public.eyJleHAiOiIyMDIzLTExLTIxVDE2OjA2OjM4KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yMVQxNDowNjozOCswNzowMCIsImlkIjoiNjU1YzNiOWI0ZjZhNzVkZGFlZWNhMTkxIiwibmJmIjoiMjAyMy0xMS0yMVQxNDowNjozOCswNzowMCIsInJvbGUiOiJhZG1pbiJ9mGLShR1CooldqYp11ygx8dJt0UNUrj4XfIegnwhriKeZSfuv-9SOcr2XG5KKO1r0hL3_V8QFCev__cJgEaTzBA")
// 	if err != nil {
// 		t.Errorf("Error decode token: %v", err)
// 	}
// 	if payload.Role != "admin" {
// 		t.Errorf("Error role: %v", err)
// 	}
// 	var dataobat model.Obat
// 	dataobat.NamaObat = "Paracetamol 500Mg"
// 	dataobat.JenisObat = "Analgesik dan Antipiretik"
// 	dataobat.Keterangan = "500 Mg"
// 	dataobat.Harga = "Rp 10.000"
// 	id := "65687b87d83eb5fac7049b29"
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

func TestUpdateObat(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
	payload, err := module.Decode("2d2bdc3a1fca7cc064174e2a6e63e2b78f1db16de9e9ed42e63646709de4a1a", "v4.public.eyJleHAiOiIyMDIzLTExLTIxVDE2OjA2OjM4KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yMVQxNDowNjozOCswNzowMCIsImlkIjoiNjU1YzNiOWI0ZjZhNzVkZGFlZWNhMTkxIiwibmJmIjoiMjAyMy0xMS0yMVQxNDowNjozOCswNzowMCIsInJvbGUiOiJhZG1pbiJ9mGLShR1CooldqYp11ygx8dJt0UNUrj4XfIegnwhriKeZSfuv-9SOcr2XG5KKO1r0hL3_V8QFCev__cJgEaTzBA")
	if err != nil {
		t.Errorf("Error decode token: %v", err)
	}
	if payload.Role != "admin" {
		t.Errorf("Error role: %v", err)
	}
	var dataobat model.Obat
	id := "656879cdd83eb5fac7049b27"
	objectId, _ := primitive.ObjectIDFromHex(id)
	dataobat.NamaObat = "Marlina Lubis"
	dataobat.JenisObat = "30/08/2003"
	dataobat.Keterangan = "Perempuan"
	dataobat.Harga = "081237629371"
	if dataobat.NamaObat == "" || dataobat.JenisObat == "" || dataobat.Keterangan == "" || dataobat.Harga == "" {
		t.Errorf("mohon untuk melengkapi data")
	} else {
		err := module.UpdateObat(objectId, payload.Id,conn,  dataobat)
		if err != nil {
			t.Errorf("Error inserting document: %v", err)
			fmt.Println("Data tidak berhasil diupdate")
		} else {
			fmt.Println("Data berhasil diupdate")
		}
	}
}

func TestDeleteObat(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
	payload, err := module.Decode("d2bdc3a1fca7cc064174e2a6e63e2b78f1db16de9e9ed42e63646709de4a1a", "v4.public.eyJleHAiOiIyMDIzLTExLTIxVDE2OjA2OjM4KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yMVQxNDowNjozOCswNzowMCIsImlkIjoiNjU1YzNiOWI0ZjZhNzVkZGFlZWNhMTkxIiwibmJmIjoiMjAyMy0xMS0yMVQxNDowNjozOCswNzowMCIsInJvbGUiOiJhZG1pbiJ9mGLShR1CooldqYp11ygx8dJt0UNUrj4XfIegnwhriKeZSfuv-9SOcr2XG5KKO1r0hL3_V8QFCev__cJgEaTzBA")
	if err != nil {
		t.Errorf("Error decode token: %v", err)
	}
	// if payload.Role != "mitra" {
	// 	t.Errorf("Error role: %v", err)
	// }
	id := "65687b87d83eb5fac7049b29"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		t.Fatalf("error converting id to objectID: %v", err)
	}
	err = module.DeleteObat(objectId, payload.Id, conn)
	if err != nil {
		t.Errorf("Error delete : %v", err)
	} else {
		fmt.Println("Success!!!")
	}
}



func TestGetAllObat(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
	data, err := module.GetAllObat(conn)
	if err != nil {
		t.Errorf("Error get all : %v", err)
	} else {
		fmt.Println(data)
	}
}

func TestGetObatFromID(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
	id := "655c3b9b4f6a75ddaeeca191"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		t.Fatalf("error converting id to objectID: %v", err)
	}
	obat, err := module.GetObatFromID(objectId, conn)
	if err != nil {
		t.Errorf("Error get obat : %v", err)
	} else {
		fmt.Println(obat)
	}
}

//order
func TestGetOrderFromID(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "healhero_db")
	id := "6565cd8f8b1e02b3244cd8a8"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		t.Fatalf("error converting id to objectID: %v", err)
	}
	order, err := module.GetOrderFromID(objectId, conn)
	if err != nil {
		t.Errorf("Error get order : %v", err)
	} else {
		fmt.Println(order)
	}
}


func TestReturnStruct(t *testing.T){
	id := "655c4c1b35843ed5631c903b"
	objectId, _ := primitive.ObjectIDFromHex(id)
	user, _ := module.GetUserFromID(objectId, db)
	data := model.User{ 
		ID : user.ID,
		Email: user.Email,
		Role : user.Role,
	}
	hasil := module.GCFReturnStruct(data)
	fmt.Println(hasil)
}

