package module

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/HealHeroo/be_healhero/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Response model.Response
	user model.User
	pengguna model.Pengguna
	driver model.Driver
	obat model.Obat
	order model.Order
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

func GCFHandlerUpdatePasswordUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)	
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	var password model.Password
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
		return GCFHandlerGetAllUserByAdmin(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname, r)
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

func GCFHandlerGetAllUserByAdmin(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
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
		Response.Message = "Anda bukan admin"
		return GCFReturnStruct(Response)
	}
	data, err := GetAllUser(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
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

func GCFHandlerGetPenggunaFromID(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "pengguna" {
		Response.Message = "Anda bukan pengguna"
		return GCFReturnStruct(Response)
	}
	data, err := GetPenggunaFromAkun(payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

// driver
func GCFHandlerUpdateDriver(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "driver" {
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
	var datadriver model.Driver
	err = json.NewDecoder(r.Body).Decode(&datadriver)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateDriver(idparam, payload.Id, conn, datadriver)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Update Driver"
	return GCFReturnStruct(Response)
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
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "driver" {
		Response.Message = "Kamu bukan driver"
		return GCFReturnStruct(Response)
	}
	data, err := GetDriverFromAkun(payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

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
	err = InsertOrder(payload.Id, conn, dataorder)
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
		Response.Message = "Missing ID parameter"
		return GCFReturnStruct(Response)
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