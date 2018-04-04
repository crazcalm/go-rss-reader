package database

import (
	"fmt"
	"os"
	"testing"
)

func TestExist(t *testing.T) {
	tests := []struct {
		File     string
		Expected bool
	}{
		{"db.go", true},
		{"DoesNotExist", false},
	}

	for _, test := range tests {
		result := Exist(test.File)

		if result != test.Expected {
			t.Errorf("For file %s, expected existence was %t, but got %t", test.File, test.Expected, result)
		}
	}
}

func TestInit(t *testing.T) {
	file := "./testing/init_test_file.db"

	//Need to create the test db file
	_, err := os.Create(file)
	if err != nil {
		t.Errorf("Unexpected error when create the database: %s", err.Error())
	}

	dbPath := fmt.Sprintf("file:%s?_foreign_keys=1", file)

	tests := []struct {
		File  string
		Reset bool
	}{
		{dbPath, false},
		{dbPath, true},
	}

	for _, test := range tests {

		_, err := Init(test.File, test.Reset)

		if err != nil {
			t.Errorf("Unexpected err: %s", err.Error())
		}
	}
	//Remove db
	err = os.Remove(file)
	if err != nil {
		t.Errorf("Error while removing file (%s): %s", dbPath, err.Error())
	}
}
