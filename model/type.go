package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username        string            `bson:"username,omitempty" json:"username,omitempty"`
	Password     string				 `bson:"password,omitempty" json:"password,omitempty"`
}

type UserBiodata struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama         string             `bson:"nama,omitempty" json:"nama,omitempty"`
	Phone_number string             `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Email        string             `bson:"email,omitempty" json:"email,omitempty"`
	Umur          int               `bson:"umur,omitempty" json:"umur,omitempty"`
	Jenis_Kelamin  string     	 `bson:"jenis_kelamin,omitempty" json:"jenis_kelamin,omitempty"`
}

type Artikel struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Judul       string             `bson:"judul,omitempty" json:"judul,omitempty"`
	Konten      string             `bson:"konten,omitempty" json:"konten,omitempty"`
	Diterbitkan string             `bson:"diterbitkan,omitempty" json:"diterbitkan,omitempty"`
	Biodata      User              `bson:"user,omitempty" json:"user,omitempty"`
}

type Admin struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username    string             `bson:"username,omitempty" json:"username,omitempty"`
	Password    string             `bson:"password,omitempty" json:"password,omitempty"`
}

type Obat struct {
	ID       primitive.ObjectID     `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama_obat    string             `bson:"nama_obat,omitempty" json:"nama_obat,omitempty"`
	Jenis_obat    string            `bson:"jenis_obat,omitempty" json:"jenis_obat,omitempty"`
	Kategori_obat string            `bson:"kategori_obat,omitempty" json:"kategori_obat,omitempty"`
	Dosis		  string			`bson:"dosis,omitempty" json:"dosis,omitempty"`
	Tanggal_exp   time.Time	    `bson:"tanggal_exp,omitempty" json:"tanggal_exp,omitempty"`
}
