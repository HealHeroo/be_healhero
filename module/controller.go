package module

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/HealHeroo/be_healhero/model"
	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/argon2"
)

// var MongoString string = os.Getenv("MONGOSTRING")

func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

// crud
func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error GetAllDocs %s: %s", col, err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return err
	}
	return docs
}

func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func UpdateOneDoc(id primitive.ObjectID, db *mongo.Database, col string, doc interface{}) (err error) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("error update: %v", err)
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("tidak ada data yang diubah")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

// signup
func SignUpPengguna(db *mongo.Database, insertedDoc model.Pengguna) error {
	objectId := primitive.NewObjectID() 
	if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" || insertedDoc.Akun.Email == "" || insertedDoc.Akun.Password == "" {
		return fmt.Errorf("Dimohon untuk melengkapi data")
	} 
	if err := checkmail.ValidateFormat(insertedDoc.Akun.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	} 
	userExists, _ := GetUserFromEmail(insertedDoc.Akun.Email, db)
	if insertedDoc.Akun.Email == userExists.Email {
		return fmt.Errorf("Email sudah terdaftar")
	} 
	isValid, _ := ValidatePhoneNumber(insertedDoc.NomorHP)
	if !isValid {
		return fmt.Errorf("Nomor telepon tidak valid")
	}
	if strings.Contains(insertedDoc.Akun.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Akun.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	} 
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Akun.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"_id": objectId,
		"email": insertedDoc.Akun.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": "pengguna",
	}
	pengguna := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"tanggallahir": insertedDoc.TanggalLahir,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp": insertedDoc.NomorHP,
		"alamat": insertedDoc.Alamat,
		"akun": model.User {
			ID : objectId,
		},
	}
	_, err = InsertOneDoc(db, "user", user)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	_, err = InsertOneDoc(db, "pengguna", pengguna)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	// Send whatsapp confirmation
	err = SendWhatsAppConfirmation(insertedDoc.NamaLengkap, insertedDoc.NomorHP)
	if err != nil {
		return fmt.Errorf("SendWhatsAppConfirmation: %v", err)
	}
	return nil
}

func ValidatePhoneNumber(nomorhp string) (bool, error) {
	// Define the regular expression pattern for numeric characters
	numericPattern := `^[0-9]+$`

	// Compile the numeric pattern
	numericRegexp, err := regexp.Compile(numericPattern)
	if err != nil {
		return false, err
	}
	// Check if the phone number consists only of numeric characters
	if !numericRegexp.MatchString(nomorhp) {
		return false, nil
	}

	// Define the regular expression pattern for "62" followed by 6 to 12 digits
	pattern := `^62\d{6,13}$`

	// Compile the regular expression
	regexpPattern, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	// Test if the phone number matches the pattern
	isValid := regexpPattern.MatchString(nomorhp)

	return isValid, nil
}


func SignUpDriver(db *mongo.Database, insertedDoc model.Driver) error {
	objectId := primitive.NewObjectID()
	if insertedDoc.NamaLengkap == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" || insertedDoc.PlatMotor == "" ||  insertedDoc.Akun.Email == "" || insertedDoc.Akun.Password == "" {
		return fmt.Errorf("Dimohon untuk melengkapi data")
	} 
	if err := checkmail.ValidateFormat(insertedDoc.Akun.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	} 
	userExists, _ := GetUserFromEmail(insertedDoc.Akun.Email, db)
	if insertedDoc.Akun.Email == userExists.Email {
		return fmt.Errorf("Email sudah terdaftar")
	} 
	if strings.Contains(insertedDoc.Akun.Password, " ") {
		return fmt.Errorf("Password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Akun.Password) < 8 {
		return fmt.Errorf("Password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Akun.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"_id": objectId,
		"email": insertedDoc.Akun.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": "driver",
	}
	driver := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp": insertedDoc.NomorHP,
		"alamat": insertedDoc.Alamat,
		"platmotor": insertedDoc.PlatMotor,
		"akun": model.User {
			ID : objectId,
		},
	}
	_, err = InsertOneDoc(db, "user", user)
	if err != nil {
		return err
	}
	_, err = InsertOneDoc(db, "driver", driver)
	if err != nil {
		return err
	}
	return nil
}

// login
func LogIn(db *mongo.Database, insertedDoc model.User) (user model.User, err error) {
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return user, fmt.Errorf("Dimohon untuk melengkapi data")
	} 
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return user, fmt.Errorf("Email tidak valid")
	} 
	existsDoc, err := GetUserFromEmail(insertedDoc.Email, db)
	if err != nil {
		return 
	}
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		return user, fmt.Errorf("kesalahan server : salt")
	}
	hash := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		return user, fmt.Errorf("password salah")
	}
	return existsDoc, nil
}

//user
func UpdateEmailUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.User) error {
	dataUser, err := GetUserFromID(iduser, db)
	if err != nil {
		return err
	}
	if insertedDoc.Email == "" {
		return fmt.Errorf("Dimohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	}
	existsDoc, _ := GetUserFromEmail(insertedDoc.Email, db)
	if existsDoc.Email == insertedDoc.Email {
		return fmt.Errorf("Email sudah terdaftar")
	}
	user := bson.M{
		"email": insertedDoc.Email,
		"password": dataUser.Password,
		"salt": dataUser.Salt,
		"role": dataUser.Role,
	}
	err = UpdateOneDoc(iduser, db, "user", user)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePasswordUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Password) error {
	dataUser, err := GetUserFromID(iduser, db)
	if err != nil {
		return err
	}
	salt, err := hex.DecodeString(dataUser.Salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	if insertedDoc.Newpassword == ""  {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if strings.Contains(insertedDoc.Newpassword, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Newpassword) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt = make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Newpassword), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"email": dataUser.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": dataUser.Role,
	}
	err = UpdateOneDoc(iduser, db, "user", user)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.User) error {
	dataUser, err := GetUserFromID(iduser, db)
	if err != nil {
		return err
	}
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	existsDoc, _ := GetUserFromEmail(insertedDoc.Email, db)
	if existsDoc.Email == insertedDoc.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	if strings.Contains(insertedDoc.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"email": insertedDoc.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": dataUser.Role,
	}
	err = UpdateOneDoc(iduser, db, "user", user)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUser(db *mongo.Database) (user []model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return user, fmt.Errorf("error GetAllUser mongo: %s", err)
	}
	err = cursor.All(context.Background(), &user)
	if err != nil {
		return user, fmt.Errorf("error GetAllUser context: %s", err)
	}
	return user, nil
}

func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetUserFromEmail(email string, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// pengguna
func InsertPengguna(db *mongo.Database, insertedDoc model.Pengguna) error {
	objectId := primitive.NewObjectID() 
	if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" || insertedDoc.Akun.Email == "" || insertedDoc.Akun.Password == "" {
		return fmt.Errorf("Dimohon untuk melengkapi data")
	} 
	if err := checkmail.ValidateFormat(insertedDoc.Akun.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	} 
	userExists, _ := GetUserFromEmail(insertedDoc.Akun.Email, db)
	if insertedDoc.Akun.Email == userExists.Email {
		return fmt.Errorf("Email sudah terdaftar")
	} 
	if strings.Contains(insertedDoc.Akun.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Akun.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	} 
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Akun.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"_id": objectId,
		"email": insertedDoc.Akun.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": "pengguna",
	}
	pengguna := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"tanggallahir": insertedDoc.TanggalLahir,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp": insertedDoc.NomorHP,
		"alamat": insertedDoc.Alamat,
		"akun": model.User {
			ID : objectId,
		},
	}
	_, err = InsertOneDoc(db, "user", user)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	_, err = InsertOneDoc(db, "pengguna", pengguna)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	return nil
}

//UpdatePengguna
func UpdatePengguna(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Pengguna) error {
    _, err := GetPenggunaFromID(idparam, db)
    if err != nil {
        return err
    }
    if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == ""{
        return fmt.Errorf("Dimohon untuk melengkapi data")
    } 
    pgn := bson.M{
        "namalengkap": insertedDoc.NamaLengkap,
        "tanggallahir": insertedDoc.TanggalLahir,
        "jeniskelamin": insertedDoc.JenisKelamin,
        "nomorhp": insertedDoc.NomorHP,
        "alamat": insertedDoc.Alamat,
        "akun": model.User {
            ID : pengguna.Akun.ID,
        },
    }
    err = UpdateOneDoc(idparam, db, "pengguna", pgn)
    if err != nil {
        return err
    }
    return nil
}


func DeletePengguna(idparam, iduser primitive.ObjectID, db *mongo.Database) error {
	_, err := GetPenggunaFromID(idparam, db)
	if err != nil {
		return err
	}
	err = DeleteOneDoc(idparam, db, "pengguna")
	if err != nil {
		return err
	}
	return nil
}

func GetAllPengguna(db *mongo.Database) (pengguna []model.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna mongo: %s", err)
	}
	err = cursor.All(context.Background(), &pengguna)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna context: %s", err)
	}
	return pengguna, nil
}

func GetPenggunaFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}


func GetPenggunaFromAkun(akun primitive.ObjectID, db *mongo.Database) (doc model.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{"akun._id": akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("pengguna tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}



//by admin
func GetPenggunaFromIDByAdmin(idparam primitive.ObjectID, db *mongo.Database) (pengguna model.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{
		"_id": idparam,
	}
	err = collection.FindOne(context.Background(), filter).Decode(&pengguna)
	if err != nil {
		return pengguna, fmt.Errorf("error GetPenggunaFromID mongo: %s", err)
	}
	user, err := GetUserFromID(pengguna.Akun.ID, db)
	if err != nil {
		return pengguna, fmt.Errorf("error GetPenggunaFromID mongo: %s", err)
	}
	akun := model.User{
		ID: user.ID,
		Email: user.Email,
		Role: user.Role,
	}
	pengguna.Akun = akun
	return pengguna, nil
}

func GetAllPenggunaByAdmin(db *mongo.Database) (pengguna []model.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna mongo: %s", err)
	}
	err = cursor.All(context.Background(), &pengguna)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna context: %s", err)
	}
	return pengguna, nil
}

func UpdatePenggunaByAdmin(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Pengguna) error {
	pengguna, err := GetPenggunaFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if pengguna.ID != idparam {
		return fmt.Errorf("Anda bukan pemilik data ini")
	}
	if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == ""{
		return fmt.Errorf("Dimohon untuk melengkapi data")
	} 
	pgn := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"tanggallahir": insertedDoc.TanggalLahir,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp": insertedDoc.NomorHP,
		"alamat": insertedDoc.Alamat,
		"akun": model.User {
			ID : pengguna.Akun.ID,
		},
	}
	err = UpdateOneDoc(idparam, db, "pengguna", pgn)
	if err != nil {
		return err
	}
	return nil
}

func GetDriverFromIDByAdmin(idparam primitive.ObjectID, db *mongo.Database) (driver model.Driver, err error) {
	collection := db.Collection("driver")
	filter := bson.M{
		"_id": idparam,
	}
	err = collection.FindOne(context.Background(), filter).Decode(&driver)
	if err != nil {
		return driver, err
	}
	user, err := GetUserFromID(driver.Akun.ID, db)
	if err != nil {
		return driver, err
	}
	akun := model.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	driver.Akun = akun
	return driver, nil
}

// driver

func InsertDriver(db *mongo.Database, insertedDoc model.Driver) error {
	objectId := primitive.NewObjectID()
	if insertedDoc.NamaLengkap == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" || insertedDoc.PlatMotor == "" ||  insertedDoc.Akun.Email == "" || insertedDoc.Akun.Password == "" {
		return fmt.Errorf("Dimohon untuk melengkapi data")
	} 
	if err := checkmail.ValidateFormat(insertedDoc.Akun.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	} 
	userExists, _ := GetUserFromEmail(insertedDoc.Akun.Email, db)
	if insertedDoc.Akun.Email == userExists.Email {
		return fmt.Errorf("Email sudah terdaftar")
	} 
	if strings.Contains(insertedDoc.Akun.Password, " ") {
		return fmt.Errorf("Password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Akun.Password) < 8 {
		return fmt.Errorf("Password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Akun.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"_id": objectId,
		"email": insertedDoc.Akun.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": "driver",
	}
	driver := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp": insertedDoc.NomorHP,
		"alamat": insertedDoc.Alamat,
		"platmotor": insertedDoc.PlatMotor,
		"akun": model.User {
			ID : objectId,
		},
	}
	_, err = InsertOneDoc(db, "user", user)
	if err != nil {
		return err
	}
	_, err = InsertOneDoc(db, "driver", driver)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDriver(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Driver) error {
	_, err := GetDriverFromID(idparam, db)
	if err != nil {
		return err
	}
	if insertedDoc.NamaLengkap == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" || insertedDoc.PlatMotor == ""{
		return fmt.Errorf("dimohon untuk melengkapi data")
	} 
	drv := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp": insertedDoc.NomorHP,
		"alamat": insertedDoc.Alamat,
		"platmotor": insertedDoc.PlatMotor,
		"akun": model.User {
			ID : driver.Akun.ID,
		},
	}

	err = UpdateOneDoc(idparam, db, "driver", drv)
	if err != nil {
		return err
	}
	return nil
}


func DeleteDriver(idparam, iduser primitive.ObjectID, db *mongo.Database) error {
	_, err := GetDriverFromID(idparam, db)
	if err != nil {
		return err
	}
	err = DeleteOneDoc(idparam, db, "driver")
	if err != nil {
		return err
	}
	return nil
}

func GetAllDriver(db *mongo.Database) (driver []model.Driver, err error) {
	collection := db.Collection("driver")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return driver, fmt.Errorf("error GetAllDriver mongo: %s", err)
	}
	err = cursor.All(context.Background(), &driver)
	if err != nil {
		return driver, fmt.Errorf("error GetAllDriver context: %s", err)
	}
	return driver, nil
}


func GetDriverFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Driver, err error) {
	collection := db.Collection("driver")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
            return doc, fmt.Errorf("no data found for ID %s", _id)
        }
        return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetDriverFromAkun(akun primitive.ObjectID, db *mongo.Database) (doc model.Driver, err error) {
	collection := db.Collection("driver")
	filter := bson.M{"akun._id": akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("Driver tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

//obat

func InsertObat(iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Obat) error {
	if insertedDoc.NamaObat == "" || insertedDoc.JenisObat == "" || insertedDoc.Keterangan == "" || insertedDoc.Harga == "" {
		return fmt.Errorf("Mohon untuk melengkapi data")
	}

	obt := bson.M{
		"nama_obat":    insertedDoc.NamaObat,
		"jenis_obat":   insertedDoc.JenisObat,
		"keterangan":   insertedDoc.Keterangan,
		"harga":        insertedDoc.Harga,
	}

	_, err := InsertOneDoc(db, "obat", obt)
	if err != nil {
		return fmt.Errorf("error saat menyimpan data obat: %s", err)
	}
	return nil
}

func UpdateObat(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Obat) error {
	_, err := GetObatFromID(idparam, db)
	if err != nil {
		return err
	}
	if insertedDoc.NamaObat == "" || insertedDoc.JenisObat == "" || insertedDoc.Keterangan == "" || insertedDoc.Harga == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	obt := bson.M{
		"nama_obat": insertedDoc.NamaObat,
		"jenis_obat": insertedDoc.JenisObat,
		"keterangan": insertedDoc.Keterangan,
		"harga": insertedDoc.Harga,
		
	}

	err = UpdateOneDoc(idparam, db, "obat", obt)
	if err != nil {
		return err
	}
	return nil
}


func DeleteObat(idparam, iduser primitive.ObjectID, db *mongo.Database) error {
	_, err := GetObatFromID(idparam, db)
	if err != nil {
		return err
	}
	err = DeleteOneDoc(idparam, db, "obat")
	if err != nil {
		return err
	}
	return nil
}

func GetAllObat(db *mongo.Database) (obat []model.Obat, err error) {
	collection := db.Collection("obat")
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return obat, fmt.Errorf("error GetAllObat mongo: %s", err)
	}
	err = cursor.All(context.TODO(), &obat)
	if err != nil {
		return obat, fmt.Errorf("error GetAllObat context: %s", err)
	}
	return obat, nil
}

func GetObatFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Obat, err error) {
	collection := db.Collection("obat")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

//order

// func InsertOrder(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Order) error {

// 	if insertedDoc.NamaObat == "" || insertedDoc.Quantity == "" || insertedDoc.TotalCost == "" || insertedDoc.Status == "" {
// 		return fmt.Errorf("harap lengkapi semua data order")
// 	}

// 	ord := bson.M{
// 		"pengguna": bson.M{
// 			"_id" : iduser,
// 			"namalengkap" : iduser,
// 		},
// 		"driver": bson.M{
// 			"_id" : insertedDoc.Driver.ID,
// 		},
// 		"obat": bson.M{
// 			"_id" : idparam,
// 			"nama_obat" : idparam,
// 		},
// 		"namaobat":    insertedDoc.NamaObat,
// 		"quantity":    insertedDoc.Quantity,
// 		"total_cost":   insertedDoc.TotalCost,
// 		"status":   insertedDoc.Status,
// 	}

// 	_, err := InsertOneDoc(db, "order", ord)
// 	if err != nil {
// 		return fmt.Errorf("error saat menyimpan data order: %s", err)
// 	}
// 	return nil
// }

// func InsertOrder(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Order) error {

// 	if insertedDoc.NamaObat == "" || insertedDoc.Quantity == "" || insertedDoc.TotalCost == "" || insertedDoc.Status == "" {
// 		return fmt.Errorf("harap lengkapi semua data order")
// 	}

// 	ord := bson.M{
// 		"pengguna": bson.M{
// 			"_id" : pengguna.ID,
// 		},
// 		"driver": bson.M{
// 			"_id" : insertedDoc.Driver.ID,
// 		},
// 		"obat": bson.M{
// 			"_id" : obat.ID,
	
// 		},
// 		"namaobat":    insertedDoc.NamaObat,
// 		"quantity":    insertedDoc.Quantity,
// 		"total_cost":   insertedDoc.TotalCost,
// 		"status":   insertedDoc.Status,
// 	}

// 	_, err := InsertOneDoc(db, "order", ord)
// 	if err != nil {
// 		return fmt.Errorf("error saat menyimpan data order: %s", err)
// 	}
// 	return nil
// }
func InsertOrder(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Order) error {

	if insertedDoc.NamaObat == "" || insertedDoc.Quantity == "" || insertedDoc.TotalCost == "" || insertedDoc.Status == "" {
		return fmt.Errorf("harap lengkapi semua data order")
	}

	ord := bson.M{
		"pengguna": bson.M{
			"namapengguna": insertedDoc.Pengguna.NamaLengkap,
			"alamat":       insertedDoc.Pengguna.Alamat,
			"nohp":         insertedDoc.Pengguna.NomorHP,
		},
		"driver": bson.M{
			"namadriver": insertedDoc.Driver.NamaLengkap,
		},
		"obat": bson.M{
			"namaobat": insertedDoc.Obat.NamaObat,
		},
		"namaobat":   insertedDoc.NamaObat,
		"quantity":   insertedDoc.Quantity,
		"total_cost": insertedDoc.TotalCost,
		"status":     insertedDoc.Status,
	}

	_, err := InsertOneDoc(db, "order", ord)
	if err != nil {
		return fmt.Errorf("error saat menyimpan data order: %s", err)
	}
	return nil
}


//update status pengiriman
func UpdateStatusOrder(idorder primitive.ObjectID, db *mongo.Database, insertedDoc model.Order) error {
	order, err := GetOrderFromID(idorder, db)
	if err != nil {
		return err
	}

	data := bson.M{
		"namaobat":    order.NamaObat,
		"quantity":    order.Quantity,
		"total_cost":   order.TotalCost,
		"status": insertedDoc.Status,
	}

	err = UpdateOneDoc(idorder, db, "order", data)
	if err != nil {
		return err
	}
	return nil
}


func DeleteOrder(idparam, iduser primitive.ObjectID, db *mongo.Database) error {
	_, err := GetOrderFromID(idparam, db)
	if err != nil {
		return err
	}
	err = DeleteOneDoc(idparam, db, "order")
	if err != nil {
		return err
	}
	return nil
}

func GetOrderFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Order, err error) {
	collection := db.Collection("order")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("_id tidak ditemukan")
		}
		return doc, err
	}
	return doc, nil
}


func GetAllOrder(db *mongo.Database) (order []model.Order, err error) {
	collection := db.Collection("order")
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return order, fmt.Errorf("error GetAllOrder mongo: %s", err)
	}
	err = cursor.All(context.TODO(), &order)
	if err != nil {
		return order, fmt.Errorf("error GetAllOrder context: %s", err)
	}
	
	return order, nil
}


//pesanan 

func InsertPesanan(iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Pesanan) error {

	if insertedDoc.Nama == "" || insertedDoc.Alamat == "" || insertedDoc.NomorHP == "" || insertedDoc.NamaObat == "" || insertedDoc.Quantity == "" || insertedDoc.Harga == "" || insertedDoc.TotalHarga == ""|| insertedDoc.Status == ""{
		return fmt.Errorf("harap lengkapi semua data pesanan")
	}

	psn := bson.M{

		"nama":    insertedDoc.Nama,
		"alamat":    insertedDoc.Alamat,
		"nomorhp":    insertedDoc.NomorHP,
		"namaobat":    insertedDoc.NamaObat,
		"quantity":    insertedDoc.Quantity,
		"harga":    insertedDoc.Harga,
		"totalharga":    insertedDoc.TotalHarga,
		"status":   insertedDoc.Status,
	}

	_, err := InsertOneDoc(db, "pesanan", psn)
	if err != nil {
		return fmt.Errorf("error saat menyimpan data pesanan: %s", err)
	}
	return nil
}

func DeletePesanan(idparam, iduser primitive.ObjectID, db *mongo.Database) error {
	_, err := GetPesananFromID(idparam, db)
	if err != nil {
		return err
	}
	err = DeleteOneDoc(idparam, db, "pesanan")
	if err != nil {
		return err
	}
	return nil
}



func GetAllPesanan(db *mongo.Database) (pesanan []model.Pesanan, err error) {
	collection := db.Collection("pesanan")
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return pesanan, fmt.Errorf("error GetAllPesanan mongo: %s", err)
	}
	err = cursor.All(context.TODO(), &pesanan)
	if err != nil {
		return pesanan, fmt.Errorf("error GetAllPesanan context: %s", err)
	}
	
	return pesanan, nil
}


func GetPesananFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Pesanan, err error) {
	collection := db.Collection("pesanan")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("_id tidak ditemukan")
		}
		return doc, err
	}
	return doc, nil
}

//sendwhatsapp
func SendWhatsAppConfirmation(namalengkap, nomorhp string) error {
	url := "https://api.wa.my.id/api/send/message/text"

	// Data yang akan dikirimkan dalam format JSON
	jsonStr := []byte(`{
        "to": "` + nomorhp + `",
        "isgroup": false,
        "messages": "Hello ` + namalengkap + `!!! ˗ˏˋ ♡ ˎˊ˗\nTerima kasih telah melakukan Registrasi akun di HealHeroo, silakan login atau tekan link dibawah ini untuk melanjutkan.\n⬇ ⬇ ⬇ ⬇ ⬇ \nhttps://healhero.my.id/signin.html"
    }`)

	// Membuat permintaan HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	// Menambahkan header ke permintaan
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Token", "v4.public.eyJleHAiOiIyMDI0LTAyLTIzVDE1OjA1OjI3WiIsImlhdCI6IjIwMjQtMDEtMjRUMTU6MDU6MjdaIiwiaWQiOiI2Mjg5NTgwNjg0NDU1NCIsIm5iZiI6IjIwMjQtMDEtMjRUMTU6MDU6MjdaIn2rLHTLg6rDPzKgKR4wr_smlabvVARrT-iXzbUDlp-fPfapNnPRf5_8mxqz9DnwMp_fQ5KJ5q8sPfLPB_VZSpAD")
	// req.Header.Set("Token", "v4.public.eyJleHAiOiIyMDI0LTAyLTE5VDIxOjA3OjM2WiIsImlhdCI6IjIwMjQtMDEtMjBUMjE6MDc6MzZaIiwiaWQiOiI2MjgyMzE3MTUwNjgxIiwibmJmIjoiMjAyNC0wMS0yMFQyMTowNzozNloiff1YQuHHPwSzGpisAMb9rTLP58-jCqtByzePJACBLghprkq2HXtTSbVTShc49m3GIVkU42VSl8uSGme8c4vXnQc")
	req.Header.Set("Content-Type", "application/json")

	// Melakukan permintaan HTTP POST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Menampilkan respons dari server
	fmt.Println("Response Status:", resp.Status)

	return nil
}

func GetUserFromToken(db *mongo.Database, col string, _id primitive.ObjectID) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}

	err = cols.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("no data found for ID", _id)
		} else {
			fmt.Println("error retrieving data for ID", _id, ":", err.Error())
		}
	}

	return user, nil
}

// func GetOrderByAdmin(db *mongo.Database) (order []model.Order, err error) {
// 	collection := db.Collection("order")
// 	filter := bson.M{}
// 	cursor, err := collection.Find(context.Background(), filter)
// 	if err != nil {
// 		return order, fmt.Errorf("error GetOrderByAdmin mongo: %s", err)
// 	}
// 	err = cursor.All(context.Background(), &order)
// 	if err != nil {
// 		return order, fmt.Errorf("error GetOrderByAdmin context: %s", err)
// 	}
// 	for _, m := range order {
// 		order, err := GetOrderFromID(m.Order.ID, db)
// 		if err != nil {
// 			return order, fmt.Errorf("error GetOrderByAdmin get order: %s", err)
// 		}
// 		m.order = order
// 		order, err := GetOrderFromID(m.Order.ID, db)
// 		if err != nil {
// 			return order, fmt.Errorf("error GetOrderByAdmin get order: %s", err)
// 		}
// 		m.Order = order
// 		Pengguna, _ := GetPenggunaFromID(m.Pengguna.ID, db)
// 		m.Pengguna = pengguna
// 		mentor, _ := GetMentorFromID(m.Mentor.ID, db)
// 		m.Mentor = mentor
// 		order = append(order, m)
// 		order = order[1:]
// 	}
// 	return order, nil
// }