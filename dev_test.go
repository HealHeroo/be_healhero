package HealHero

import (
	"fmt"
	"testing"

	// "github.com/HealHeroo/be_healhero/model"
	"github.com/HealHeroo/be_healhero/module"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsertUser(t *testing.T) {
	username := "rizkyria"
	password := "admin123"
	hasil := module.InsertUser(module.MongoConn, "user", username, password)
	fmt.Println(hasil)
}



func TestGetUserFromUsername(t *testing.T) {
	username := "rizkyria"
	data := module.GetUserFromUsername(username, module.MongoConn, "user")
	fmt.Println(data)
}
