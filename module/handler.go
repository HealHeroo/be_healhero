package module

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/HealHeroo/be_healhero/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Response model.Response
	user model.User
	pengguna model.Pengguna
	driver model.Driver
	obat model.Obat
	order model.Order
	password model.Password

)

// signup
func GCFHandlerSignUpPengguna(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	var datapengguna model.Pengguna
	err := json.NewDecoder(r.Body).Decode(&datapengguna)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = SignUpPengguna(conn, datapengguna)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Halo " + datapengguna.NamaLengkap
	return GCFReturnStruct(Response)
}

func GCFHandlerSignUpDriver(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	var datadriver model.Driver
	err := json.NewDecoder(r.Body).Decode(&datadriver)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = SignUpDriver(conn, datadriver)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Halo " + datadriver.NamaLengkap
	return GCFReturnStruct(Response)
}

// login
func GCFHandlerLogin(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Credential
	Response.Status = false
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	user, err := LogIn(conn, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	tokenstring, err := Encode(user.ID, user.Role, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Response.Message = "Gagal Encode Token : " + err.Error()
	} else {
		Response.Message = "Selamat Datang " + user.Email
		Response.Token = tokenstring
		Response.Role = user.Role
	}
	return GCFReturnStruct(Response)
}

// get all
func GCFHandlerGetAll(MONGOCONNSTRINGENV, dbname, col string, docs interface{}) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	data := GetAllDocs(conn, col, docs)
	return GCFReturnStruct(data)
}

// user
func GCFHandlerUpdateEmailUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateEmailUser(user_login.Id, conn, user)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Email"
	return GCFReturnStruct(Response)
}

// func GCFHandlerUpdatePasswordUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)	
// 	Response.Status = false
// 	//
// 	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
// 	if err != nil {
// 		Response.Message = "Gagal Decode Token : " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	var password model.Password
// 	err = json.NewDecoder(r.Body).Decode(&password)
// 	if err != nil {
// 		Response.Message = "error parsing application/json: " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	err = UpdatePasswordUser(user_login.Id, conn, password)
// 	if err != nil {
// 		Response.Message = err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	//
// 	Response.Status = true
// 	Response.Message = "Berhasil Update Password Kamu"
// 	return GCFReturnStruct(Response)
// }

func GCFHandlerUpdatePasswordUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdatePasswordUser(user_login.Id, conn, password)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Password"
	return GCFReturnStruct(Response)
}


func GCFHandlerUpdateUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	var datauser model.User
	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateUser(payload.Id, conn, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Update User"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "admin" {
		return GCFHandlerGetUserFromID(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname, r)
	}
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllUserByAdmin(conn)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	data, err := GetUserFromID(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if data.Role == "pengguna" {
		datapengguna, err := GetPenggunaFromAkun(data.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		datapengguna.Akun = data
		return GCFReturnStruct(datapengguna) 
	}
	if data.Role == "driver" {
		datadriver, err := GetDriverFromAkun(data.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		datadriver.Akun = data
		return GCFReturnStruct(datadriver) 
	}
	Response.Message = "Tidak ada data"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetUserFromID(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	data, err := GetUserFromID(payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}
// get
func Get(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "admin" {
		Response.Message = "Kamu BUkan Admin"
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllUserByAdmin(conn)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	user, err := GetUserFromID(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if user.Role == "pengguna" {
		pengguna, err := GetPenggunaFromAkun(user.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		return GCFReturnStruct(pengguna)
	}
	if user.Role == "driver" {
		driver, err := GetDriverFromAkun(user.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		return GCFReturnStruct(driver)
	}
	
	if user.Role == "admin" {
		admin, err := GetUserFromID(user_login.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		return GCFReturnStruct(admin)
	}
	//
	Response.Message = "Tidak ada data"
	return GCFReturnStruct(Response)
}

//email
func Put(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateEmailUser(user_login.Id, conn, user)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Email"
	return GCFReturnStruct(Response)
}


// func GCFHandlerGetAllUserByAdmin(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	var Response model.Response
// 	Response.Status = false
// 	tokenstring := r.Header.Get("Authorization")
// 	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
// 	if err != nil {
// 		Response.Message = "Gagal Decode Token : " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	if payload.Role != "admin" {
// 		Response.Message = "Anda bukan admin"
// 		return GCFReturnStruct(Response)
// 	}
// 	data, err := GetAllUser(conn)
// 	if err != nil {
// 		Response.Message = err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	return GCFReturnStruct(data)
// }

func GCFHandlerGetAllUserByAdmin(conn *mongo.Database) string {
	Response.Status = false
	//
	data, err := GetAllUser(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(data)
}

// pengguna
func GCFHandlerUpdatePengguna(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "pengguna" {
		Response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	var datapengguna model.Pengguna
	err = json.NewDecoder(r.Body).Decode(&datapengguna)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdatePengguna(idparam, payload.Id, conn, datapengguna)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Update Pengguna"
	return GCFReturnStruct(Response)
}

// func GCFHandlerUpdatePengguna(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	Response.Status = false
// 	//
// 	id := GetID(r)
// 	if id == "" {
// 		Response.Message = "Wrong parameter"
// 		return GCFReturnStruct(Response)
// 	}
// 	idparam, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		Response.Message = "Invalid id parameter"
// 		return GCFReturnStruct(Response)
// 	}
// 	err = json.NewDecoder(r.Body).Decode(&pengguna)
// 	if err != nil {
// 		Response.Message = "error parsing application/json: " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
// 	if err != nil {
// 		Response.Message = "Gagal Decode Token : " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	if user_login.Role == "pengguna" {
// 		return GCFHandlerUpdateByPengguna(idparam, user_login.Id, pengguna, conn, r)
// 	}
// 	if user_login.Role == "admin" {
// 		return GCFHandlerUpdateByAdmin(idparam, pengguna, conn, r)
// 	}
	
// 	Response.Message = "Kamu tidak memiliki akses"
// 	return GCFReturnStruct(Response)
// }

func GCFHandlerUpdateByPengguna(idparam, iduser primitive.ObjectID,  pengguna model.Pengguna, conn *mongo.Database, r *http.Request) string {
	Response.Status = false
	//
	err := UpdatePengguna(idparam, iduser, conn, pengguna)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Pengguna"
	return GCFReturnStruct(Response)
}

// func GCFHandlerUpdateByAdmin(idparam primitive.ObjectID, pengguna model.Pengguna, conn *mongo.Database, r *http.Request) string {
// 	Response.Status = false
// 	//
// 	err := UpdatePenggunaByAdmin(idparam, conn, pengguna)
// 	if err != nil {
// 		Response.Message = err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	//
// 	Response.Status = true
// 	Response.Message = "Berhasil Update Data"
// 	return GCFReturnStruct(Response)
// }


func GCFHandlerGetAllPengguna(MONGOCONNSTRINGENV, dbname string) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	data, err := GetAllPengguna(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetPengguna(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllPengguna(MONGOCONNSTRINGENV, dbname)
	}

	idParam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid ID parameter"
		return GCFReturnStruct(Response)
	}

	obat, err := GetPenggunaFromID(idParam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	return GCFReturnStruct(obat)
}


// func GCFHandlerGetPengguna(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	var Response model.Response
// 	Response.Status = false
// 	tokenstring := r.Header.Get("Authorization")
// 	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
// 	if err != nil {
// 		Response.Message = "Gagal Decode Token : " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	if payload.Role != "admin" {
// 		return GCFHandlerGetPenggunaFromID(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname, r)
// 	}
// 	id := GetID(r)
// 	if id == "" {
// 		return GCFHandlerGetAllPenggunaByAdmin(conn)
// 	}
// 	idparam, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		Response.Message = "Invalid id parameter"
// 		return GCFReturnStruct(Response)
// 	}
// 	data, err := GetUserFromID(idparam, conn)
// 	if err != nil {
// 		Response.Message = err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	if data.Role == "pengguna" {
// 		datapengguna, err := GetPenggunaFromAkun(data.ID, conn)
// 		if err != nil {
// 			Response.Message = err.Error()
// 			return GCFReturnStruct(Response)
// 		}
// 		datapengguna.Akun = data
// 		return GCFReturnStruct(datapengguna) 
// 	}
// 	Response.Message = "Tidak ada data"
// 	return GCFReturnStruct(Response)
// }

func GCFHandlerGetAllPenggunaByAdmin(conn *mongo.Database) string {
	Response.Status = false
	//
	data, err := GetAllUser(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(data)
}


func GCFHandlerGetPenggunaFromID(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "pengguna" {
		return GCFHandlerGetPenggunaByPengguna(user_login.Id, conn)
	}
	if user_login.Role == "admin" {
		return GCFHandlerGetPenggunaByAdmin(conn, r)
	}
	Response.Message = "Kamu tidak memiliki akses"
	return GCFReturnStruct(Response)
	
}



func GCFHandlerGetPenggunaByAdmin(conn *mongo.Database, r *http.Request) string {
	Response.Status = false
	//
	id := GetID(r)
	if id == "" {
		pengguna, err := GetAllPenggunaByAdmin(conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		return GCFReturnStruct(pengguna)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	pengguna, err := GetPenggunaFromIDByAdmin(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(pengguna)
}

func GCFHandlerGetPenggunaByPengguna(iduser primitive.ObjectID, conn *mongo.Database) string {
	Response.Status = false
	//
	pengguna, err := GetPenggunaFromAkun(iduser, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(pengguna)
}

// driver
// func GCFHandlerUpdateDriver(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	var Response model.Response
// 	Response.Status = false
// 	tokenstring := r.Header.Get("Authorization")
// 	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
// 	if err != nil {
// 		Response.Message = "Gagal Decode Token : " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	if payload.Role != "driver" {
// 		Response.Message = "Anda tidak memiliki akses"
// 		return GCFReturnStruct(Response)
// 	}
// 	id := GetID(r)
// 	if id == "" {
// 		Response.Message = "Wrong parameter"
// 		return GCFReturnStruct(Response)
// 	}
// 	idparam, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		Response.Message = "Invalid id parameter"
// 		return GCFReturnStruct(Response)
// 	}
// 	var datadriver model.Driver
// 	err = json.NewDecoder(r.Body).Decode(&datadriver)
// 	if err != nil {
// 		Response.Message = "error parsing application/json: " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	err = UpdateDriver(idparam, payload.Id, conn, datadriver)
// 	if err != nil {
// 		Response.Message = err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	Response.Status = true
// 	Response.Message = "Berhasil Update Driver"
// 	return GCFReturnStruct(Response)
// }

func GCFHandlerUpdateDriver(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Driver) error {
	driver, err := GetDriverFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if driver.ID != idparam {
		return fmt.Errorf("kamu bukan pemilik data ini")
	}
	if insertedDoc.NamaLengkap == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" || insertedDoc.PlatMotor == ""  {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	mtr := bson.M{
		"namalengkap": 		insertedDoc.NamaLengkap,
		"jeniskelamin": 	insertedDoc.JenisKelamin,
		"nomorhp":          insertedDoc.NomorHP,
		"alamat":     	 	insertedDoc.Alamat,
		"platmotor": 		insertedDoc.PlatMotor,
		"akun": model.User{
			ID: driver.Akun.ID,
		},
	}
	err = UpdateOneDoc(idparam, db, "driver", mtr)
	if err != nil {
		return err
	}
	return nil
}

func GCFHandlerGetDriver(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllDriver(MONGOCONNSTRINGENV, dbname)
	}

	idParam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid ID parameter"
		return GCFReturnStruct(Response)
	}

	obat, err := GetDriverFromID(idParam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	return GCFReturnStruct(obat)
}

func GCFHandlerGetAllDriver(MONGOCONNSTRINGENV, dbname string) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	data, err := GetAllDriver(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}


func GCFHandlerGetDriverFromID(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "driver" {
		return GCFHandlerGetDriverByDriver(user_login.Id, conn)
	}
	if user_login.Role == "admin" {
		return GCFHandlerGetDriverByAdmin(conn, r)
	}
	Response.Message = "Kamu tidak memiliki akses"
	return GCFReturnStruct(Response)
	
}

func GCFHandlerGetDriverByAdmin(conn *mongo.Database, r *http.Request) string {
	Response.Status = false
	//
	id := GetID(r)
	if id == "" {
		pengguna, err := GetAllPenggunaByAdmin(conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		return GCFReturnStruct(pengguna)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	pengguna, err := GetPenggunaFromIDByAdmin(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(pengguna)
}

func GCFHandlerGetDriverByDriver(iduser primitive.ObjectID, conn *mongo.Database) string {
	Response.Status = false
	//
	pengguna, err := GetPenggunaFromAkun(iduser, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(pengguna)
}

// func GCFHandlerGetDriverFromID(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	var Response model.Response
// 	Response.Status = false
// 	tokenstring := r.Header.Get("Authorization")
// 	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
// 	if err != nil {
// 		Response.Message = err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	if payload.Role != "driver" {
// 		Response.Message = "Maaf Kamu bukan driver"
// 		return GCFReturnStruct(Response)
// 	}
// 	data, err := GetDriverFromAkun(payload.Id, conn)
// 	if err != nil {
// 		Response.Message = err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	return GCFReturnStruct(data)
// }

// obat
func GCFHandlerInsertObat(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	var dataobat model.Obat
	err = json.NewDecoder(r.Body).Decode(&dataobat)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = InsertObat(payload.Id, conn, dataobat)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Insert Obat"
	return GCFReturnStruct(Response)
}

func GCFHandlerUpdateObat(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	var dataobat model.Obat
	err = json.NewDecoder(r.Body).Decode(&dataobat)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateObat(idparam, payload.Id, conn, dataobat)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Update Obat"
	return GCFReturnStruct(Response)
}

// func GCFHandlerUpdateObat(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	Response.Status = false
// 	//
// 	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
// 	if err != nil {
// 		Response.Message = "Gagal Decode Token : " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	if user_login.Role != "admin" {
// 		Response.Message = "Kamu tidak memiliki akses"
// 		return GCFReturnStruct(Response)
// 	}
// 	id := GetID(r)
// 	if id == "" {
// 		Response.Message = "Wrong parameter"
// 		return GCFReturnStruct(Response)
// 	}
// 	idobat, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		Response.Message = "Invalid id parameter"
// 		return GCFReturnStruct(Response)
// 	}
// 	err = json.NewDecoder(r.Body).Decode(&obat)
// 	if err != nil {
// 		Response.Message = "error parsing application/json: " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	err = UpdateObat(idobat, user_login.Id, conn, obat)
// 	if err != nil {
// 		Response.Message = err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	//
// 	Response.Status = true
// 	Response.Message = "Berhasil Update Obat"
// 	return GCFReturnStruct(Response)
// }

func GCFHandlerDeleteObat(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	err = DeleteObat(idparam, payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Delete Obat"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllObat(MONGOCONNSTRINGENV, dbname string) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	data, err := GetAllObat(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetObatFromID(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllObat(MONGOCONNSTRINGENV, dbname)
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	data, err := GetObatFromID(objID, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}


func GCFHandlerGetObat(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllObat(MONGOCONNSTRINGENV, dbname)
	}

	idParam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid ID parameter"
		return GCFReturnStruct(Response)
	}

	obat, err := GetObatFromID(idParam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	return GCFReturnStruct(obat)
}


//order
func GCFHandlerInsertOrder(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	var dataorder model.Order
	err = json.NewDecoder(r.Body).Decode(&dataorder)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}

	idParam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid ID parameter"
		return GCFReturnStruct(Response)
	}
	err = InsertOrder(idParam, payload.Id, conn, dataorder)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Insert Order"
	return GCFReturnStruct(Response)
}



func GCFHandlerDeleteOrder(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	err = DeleteOrder(idparam, payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Delete Order"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllOrder(MONGOCONNSTRINGENV, dbname string) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	data, err := GetAllOrder(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetOrderFromID(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllOrder(MONGOCONNSTRINGENV, dbname)
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	data, err := GetOrderFromID(objID, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetOrder(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllOrder(MONGOCONNSTRINGENV, dbname)
	}

	idParam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid ID parameter"
		return GCFReturnStruct(Response)
	}

	obat, err := GetOrderFromID(idParam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	return GCFReturnStruct(obat)
}


//pesanan


// func GCFHandlerInsertPesanan(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	var Response model.Response
// 	Response.Status = false
// 	tokenstring := r.Header.Get("Authorization")
// 	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
// 	if err != nil {
// 		Response.Message = "Gagal Decode Token : " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	var datapesanan model.Pesanan
// 	err = json.NewDecoder(r.Body).Decode(&datapesanan)
// 	if err != nil {
// 		Response.Message = "error parsing application/json: " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	err = InsertPesanan(payload.Id, conn, datapesanan)
// 	if err != nil {
// 		Response.Message = err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	Response.Status = true
// 	Response.Message = "Berhasil Insert Pesanan"
// 	return GCFReturnStruct(Response)
// }


func GCFHandlerInsertPesanan(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	var datapesanan model.Pesanan
	err = json.NewDecoder(r.Body).Decode(&datapesanan)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}

	idParam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid ID parameter"
		return GCFReturnStruct(Response)
	}
	err = InsertPesanan(idParam, payload.Id, conn, datapesanan)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Insert Pesanan"
	return GCFReturnStruct(Response)
}


func GCFHandlerDeletePesanan(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	err = DeletePesanan(idparam, payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Delete Pesanan"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllPesanan(MONGOCONNSTRINGENV, dbname string) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	data, err := GetAllPesanan(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetPesananFromID(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllPesanan(MONGOCONNSTRINGENV, dbname)
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	data, err := GetOrderFromID(objID, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetPesanan(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllPesanan(MONGOCONNSTRINGENV, dbname)
	}

	idParam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid ID parameter"
		return GCFReturnStruct(Response)
	}

	pesanan, err := GetPesananFromID(idParam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	return GCFReturnStruct(pesanan)
}



// return struct
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// get user login
func GetUserLogin(PASETOPUBLICKEYENV string, r *http.Request) (model.Payload, error) {
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		return payload, err
	}
	return payload, nil
}

// get id
func GetID(r *http.Request) string {
    return r.URL.Query().Get("id")
}