package module

import (
	"context"
	"fmt"
	"os"

	"github.com/HealHeroo/be_healhero/model"
	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoString string = os.Getenv("MONGOSTRING")

var MongoInfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "healhero_db",
}

var MongoConn = atdb.MongoConnect(MongoInfo)

func InsertOneDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}

func InsertUser(db *mongo.Database, col string,username string,password string) (InsertedID interface{}) {
	var user model.User
	user.Username = username
	user.Password = password
	return InsertOneDoc(db, col, user)
}

func GetUserFromUsername(username string, db *mongo.Database, col string) (user model.UserBiodata) {
	data_profile := db.Collection(col)
	filter := bson.M{"username": username}
	err := data_profile.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("getUserFromUsername: %v\n", err)
	}
	return user
}


func InsertUser_Biodata(db *mongo.Database, col string,nama string,phone_number string, email string, umur int,jenis_kelamin string) (InsertedID interface{}) {
	var userbiodata model.UserBiodata
	userbiodata.Nama = nama
	userbiodata.Phone_number = phone_number
	userbiodata.Email = email
	userbiodata.Umur = umur
	userbiodata.Jenis_Kelamin = jenis_kelamin
	return InsertOneDoc(db, col, userbiodata)
}

func GetUserBiodataFromEmail(email string, db *mongo.Database, col string) (userbiodata model.UserBiodata) {
	data_profile := db.Collection(col)
	filter := bson.M{"email": email}
	err := data_profile.FindOne(context.TODO(), filter).Decode(&userbiodata)
	if err != nil {
		fmt.Printf("getUserFromEmail: %v\n", err)
	}
	return userbiodata
}


func InsertArtikel(db *mongo.Database, col string, judul string,konten string, diterbitkan string, biodata model.User) (InsertedID interface{}) {
	var artikel model.Artikel
	artikel.Judul = judul
	artikel.Konten = konten
	artikel.Diterbitkan = diterbitkan
	artikel.Biodata = biodata
	return InsertOneDoc(db,col, artikel)
}


func GetArtikelFromJudul(judul string, db *mongo.Database, col string) (artikel model.Artikel) {
	data_profile := db.Collection(col)
	filter := bson.M{"judul": judul}
	err := data_profile.FindOne(context.TODO(), filter).Decode(&artikel)
	if err != nil {
		fmt.Printf("getArtikelFromJudul: %v\n", err)
	}
	return artikel
}

func InsertAdmin(db *mongo.Database, col string, username string, password string) (InsertedID interface{}) {
	var admin model.Admin
	admin.Username = username
	admin.Password = password
	return InsertOneDoc(db, col, admin)
}


func GetAdminFromUsername(username string, db *mongo.Database, col string) (admin model.Admin) {
	data_profile := db.Collection(col)
	filter := bson.M{"username": username}
	err := data_profile.FindOne(context.TODO(), filter).Decode(&admin)
	if err != nil {
		fmt.Printf("getAdminFromUsername: %v\n", err)
	}
	return admin
}



