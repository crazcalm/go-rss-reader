package database

import (
	"strings"
	"testing"
)

func TestAddAuthor(t *testing.T) {
	file := "./testing/add_author.db"
	db := createTestDB(file)

	name := "John Doe"
	email := "John.Doe@gmail.com"

	_, err := AddAuthor(db, name, email)
	if err != nil {
		t.Errorf("Failed to add author: %s", err.Error())
	}
}

func TestGetAuthor(t *testing.T) {
	file := "./testing/get_author.db"
	db := createTestDB(file)

	name := "John Doe"
	email := "John.Doe@gmail.com"
	id, err := AddAuthor(db, name, email)
	if err != nil {
		t.Errorf("Error happened when trying to add a author: %s", err.Error())
	}

	dbName, dbEmail, err := GetAuthor(db, id)
	if err != nil {
		t.Errorf("Failed to get author: %s", err.Error())
	}

	if !strings.EqualFold(name, dbName) || !strings.EqualFold(email, dbEmail) {
		t.Errorf("Expected name and email to be %s and %s, but got %s and %s", name, email, dbName, dbEmail)
	}

}
