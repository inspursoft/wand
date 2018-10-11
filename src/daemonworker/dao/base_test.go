package dao_test

import (
	"log"
	"os"
	"testing"

	"github.com/inspursoft/wand/src/daemonworker/dao"
)

func TestInitDb(t *testing.T) {
	dao.InitDB()
}

func TestAddOrUpdateUserAccess(t *testing.T) {
	err := dao.AddOrUpdateUserAccess("tester1", "123456")
	if err != nil {
		log.Printf("Failed to insert data:%+v\n", err)
	}
	err = dao.AddOrUpdateUserAccess("tester1", "456789")
	if err != nil {
		log.Printf("Failed to insert data:%+v\n", err)
	}
}

func TestGetUserAccess(t *testing.T) {
	userAccess := dao.GetUserAccess("tester1")
	log.Printf("User access:%+v\n", userAccess)
}

func TestCleanUp(t *testing.T) {
	os.Remove("storage.db")
}
