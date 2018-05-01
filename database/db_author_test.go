package database

import (
	"strings"
	"testing"
)

func TestAuthorExist(t *testing.T) {
	file := "./testing/author_exist.db"
	db := createTestDB(file)

	name := "Jane Doe"
	email := "Jane.Doe@gmail.com"
	notName := "Not a name"
	notEmail := "Not an email"

	_, err := AddAuthor(db, name, email)
	if err != nil {
		t.Errorf("Failed to add author: %s", err.Error())
	}

	tests := []struct {
		Name   string
		Email  string
		Expect bool
	}{
		{name, email, true},
		{name, notEmail, false},
		{notName, email, false},
		{notName, notEmail, false},
	}

	for _, test := range tests {
		answer := AuthorExist(db, test.Name, test.Email)

		if answer != test.Expect {
			t.Errorf("When checking if the author (%s, %s) existed: Expected %t, but got %t", test.Name, test.Email, test.Expect, answer)
		}
	}
}

func TestGetAuthorByNameAndEmail(t *testing.T) {
	file := "./testing/get_author_by_name_and_email.db"
	db := createTestDB(file)

	name := "Jane Doe"
	email := "Jane.Doe@gmail.com"
	notName := "Not a name"
	notEmail := "Not an email"

	authorID, err := AddAuthor(db, name, email)
	if err != nil {
		t.Errorf("Failed to add author: %s", err.Error())
	}

	tests := []struct {
		Name  string
		Email string
		ID    int64
		Error bool
	}{
		{name, email, authorID, false},
		{name, notEmail, authorID, true},
		{notName, email, authorID, true},
		{notName, notEmail, authorID, true},
	}

	for _, test := range tests {
		dbID, err := GetAuthorByNameAndEmail(db, test.Name, test.Email)

		if test.Error && err == nil {
			t.Error("Expected and error, but none was received")
		}

		if !test.Error && err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		//Wanted error case
		if test.Error && err != nil {
			continue
		}

		if dbID != test.ID {
			t.Errorf("Expected the id to be %d, but got %d", test.ID, dbID)
		}
	}
}

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
