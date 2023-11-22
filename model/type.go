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
	Email 	 string            `bson:"email" json:"email"`
	Password string            `bson:"password" json:"password"`
	Salt 	 string			   `bson:"salt,omitempty" json:"salt,omitempty"`
	Role     string            `bson:"role" json:"role"`
}

type Password struct {
	Password        string         	   `bson:"password,omitempty" json:"password,omitempty"`
	Newpassword 	string         	   `bson:"newpass,omitempty" json:"newpass,omitempty"`
}

type Pengguna struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap  	string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	TanggalLahir	string             `bson:"tanggallahir,omitempty" json:"tanggallahir,omitempty"`
	JenisKelamin  	string             `bson:"jeniskelamin,omitempty" json:"jeniskelamin,omitempty"`
	NomorHP  		string             `bson:"nomorhp,omitempty" json:"nomorhp,omitempty"`
	Alamat			string             `bson:"alamat,omitempty" json:"alamat,omitempty"`
	Akun     		User               `bson:"akun" json:"akun"`
}

type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Akun     User            	`bson:"akun" json:"akun"`
}

type Driver struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap  	string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	JenisKelamin  	string             `bson:"jeniskelamin,omitempty" json:"jeniskelamin,omitempty"`
	NomorHP  		string             `bson:"nomorhp,omitempty" json:"nomorhp,omitempty"`
	Alamat			string             `bson:"alamat,omitempty" json:"alamat,omitempty"`
	PlatMotor  		string             `bson:"platmotor,omitempty" json:"platmotor,omitempty"`
	JenisMotor  		string         `bson:"jenismotor,omitempty" json:"jenismotor,omitempty"`
	Akun     		User           	   `bson:"akun" json:"akun"`
}

type Obat struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaObat    string            `json:"nama_obat" bson:"nama_obat"`
	JenisObat   string            `json:"jenis_obat" bson:"jenis_obat"`
	Keterangan  string            `json:"keterangan" bson:"keterangan"`
	Harga       string           `json:"harga" bson:"harga"`
}

type Order struct {
	ID        	  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PenggunaID    primitive.ObjectID `bson:"pengguna_id" json:"pengguna_id"`
	DriverID  	  primitive.ObjectID `bson:"driver_id" json:"driver_id"`
	ObatID    	  primitive.ObjectID `bson:"obat_id" json:"obat_id"`
	NamaObat  	  string             `bson:"namaobat" json:"namaobat"`
	Quantity  	  string             `bson:"quantity" json:"quantity"`
	TotalCost 	  string           	 `bson:"total_cost" json:"total_cost"`
	Status    	  string           	 `bson:"status" json:"status"`
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