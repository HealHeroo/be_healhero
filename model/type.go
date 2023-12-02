// package HealHero

// import (
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	// "time"
// )

// type User struct {
// 	ID           primitive.ObjectID 	`bson:"_id,omitempty" json:"_id,omitempty"`
// 	Username string `json:"username" bson:"username"`
// 	Password string `json:"password" bson:"password"`
// 	// Email		 string             	`bson:"email,omitempty" json:"email,omitempty"`
// }

// type Credential struct {
// 	Status  bool   `json:"status" bson:"status"`
// 	Token   string `json:"token,omitempty" bson:"token,omitempty"`
// 	Message string `json:"message,omitempty" bson:"message,omitempty"`
// }

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email 	 string            `bson:"email,omitempty" json:"email,omitempty"`
	Password string            `bson:"password,omitempty" json:"password,omitempty"`
	Salt 	 string			   `bson:"salt,omitempty,omitempty" json:"salt,omitempty"`
	Role     string            `bson:"role,omitempty" json:"role,omitempty"`
}

type Password struct {
	Password        string         	   `bson:"password,omitempty" json:"password,omitempty"`
	Newpassword 	string         	   `bson:"newpass,omitempty" json:"newpass,omitempty"`
}

type Pengguna struct {
	ID       primitive.ObjectID 		`bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap  	string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	TanggalLahir	string             `bson:"tanggallahir,omitempty" json:"tanggallahir,omitempty"`
	JenisKelamin  	string             `bson:"jeniskelamin,omitempty" json:"jeniskelamin,omitempty"`
	NomorHP  		string             `bson:"nomorhp,omitempty" json:"nomorhp,omitempty"`
	Alamat			string             `bson:"alamat,omitempty" json:"alamat,omitempty"`
	Akun     		User               `bson:"akun,omitempty" json:"akun,omitempty"`
}

type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty,omitempty" json:"_id,omitempty"`
	Akun     User            	`bson:"akun,omitempty" json:"akun,omitempty"`
}

type Driver struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap  	string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	JenisKelamin  	string             `bson:"jeniskelamin,omitempty" json:"jeniskelamin,omitempty"`
	NomorHP  		string             `bson:"nomorhp,omitempty" json:"nomorhp,omitempty"`
	Alamat			string             `bson:"alamat,omitempty" json:"alamat,omitempty"`
	PlatMotor  		string             `bson:"platmotor,omitempty" json:"platmotor,omitempty"`
	Akun     		User           	   `bson:"akun,omitempty" json:"akun,omitempty"`
}

type Obat struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaObat    string            `json:"nama_obat,omitempty" bson:"nama_obat,omitempty"`
	JenisObat   string            `json:"jenis_obat,omitempty" bson:"jenis_obat,omitempty"`
	Keterangan  string            `json:"keterangan,omitempty" bson:"keterangan,omitempty"`
	Harga       string           `json:"harga,omitempty" bson:"harga,omitempty"`
}

type Pesanan struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama    string            `json:"nama,omitempty" bson:"nama,omitempty"`
	Alamat   string            `json:"alamat,omitempty" bson:"alamat,omitempty"`
	NomorHP   string            `json:"nomorhp,omitempty" bson:"nomorhp,omitempty"`
	NamaObat  string            `json:"namaobat,omitempty" bson:"namaobat,omitempty"`
	Quantity  string            `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Harga       string           `json:"harga,omitempty" bson:"harga,omitempty"`
	TotalHarga       string           `json:"totalharga,omitempty" bson:"totalharga,omitempty"`
	Status    	  string           	 `json:"status,omitempty" bson:"status,omitempty"`
}

type Order struct {
	ID        	  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Pengguna    	Pengguna `bson:"pengguna,omitempty" json:"pengguna,omitempty"`
	Driver  	  Driver `bson:"driver,omitempty" json:"driver,omitempty"`
	Obat    	  Obat `bson:"obat,omitempty" json:"obat,omitempty"`
	NamaObat  	  string             `bson:"namaobat,omitempty" json:"namaobat,omitempty"`
	Quantity  	  string             `bson:"quantity,omitempty" json:"quantity,omitempty"`
	TotalCost 	  string           	 `bson:"total_cost,omitempty" json:"total_cost,v"`
	Status    	  string           	 `bson:"status,omitempty" json:"status,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Role	string `json:"role,omitempty" bson:"role,omitempty"`
}

type Response struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Payload struct {
	Id           	primitive.ObjectID `json:"id"`
	Role           	string			   `json:"role"`
	Exp 			time.Time 	 	   `json:"exp"`
	Iat 			time.Time 		   `json:"iat"`
	Nbf 			time.Time 		   `json:"nbf"`
}