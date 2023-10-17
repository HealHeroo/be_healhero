package HealHero

import (
	"fmt"
	"testing"
)

func TestInsertUser(t *testing.T) {
	username := "rizkyria"
	password := "admin123"
	hasil:=TestInsertUser(username ,password)
	fmt.Println(hasil)
}

func TestGetUserFromUsername(t *testing.T) {
	username := "rizkyria"
	biodata:=GetUserFromUsername(username)
	fmt.Println(biodata)
}

