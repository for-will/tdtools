package nosql

import "testing"

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