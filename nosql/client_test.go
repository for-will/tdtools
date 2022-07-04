package nosql

import (
	"fmt"
	"testing"
)

func TestFIndOne(t *testing.T) {
	FindOne()
}

func TestUpsert(t *testing.T) {
	UpdateOne()
}

func TestUpsertInterface(t *testing.T) {
	InsertOne()
}

func TestInsertMany(t *testing.T) {
	InsertMany()
}

func TestFindMany(t *testing.T) {
	FindMany()
}

func TestUpsert1(t *testing.T) {
	for i := 10; i < 20; i++ {
		userId := fmt.Sprintf("%v", i+10000)
		UniqueInsert(userId, "120000")
		UniqueInsert(userId, "110000")
	}
}
