package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
)

func createTestDB(file string) *sql.DB {
	testDB := fmt.Sprintf("file:%s%s", file, foreignKeySupport)

	db, err := Init(testDB, true)
	if err != nil {
		log.Fatal("Error when trying to create the database")
	}
	return db
}

func TestFeedExist(t *testing.T) {

}

func TestTagExist(t *testing.T) {

}

func TestFeedHasTag(t *testing.T) {

}

func TestAddFeedURL(t *testing.T) {
	file := "./testing/add_feed_url.db"
	db := createTestDB(file)

	tests := []struct {
		URL       string
		Count     int64
		ExpectErr bool
	}{
		{"url1", 1, false},
		{"url1", 1, true}, // Tests FeedURLExist
	}

	for _, test := range tests {
		_, err := AddFeedURL(db, test.URL)

		if err != nil && !test.ExpectErr {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		if err == nil && test.ExpectErr {
			t.Errorf("Expected an error, but none was received")
		}

		//Case: expected error, received an error
		if err != nil && test.ExpectErr {
			continue
		}

		if err == nil && !test.ExpectErr {
			var count int64
			row := db.QueryRow("SELECT COUNT(*) FROM feeds")
			err := row.Scan(&count)
			if err != nil {
				t.Errorf("Error happened when trying to obtain count of feeds")
			}

			if count != test.Count {
				t.Errorf("Expected %d feeds, but got %d", test.Count, count)
			}
		}
	}
}

func TestAddTag(t *testing.T) {
	file := "./testing/add_tag.db"
	db := createTestDB(file)

	tests := []struct {
		Tag       string
		Count     int64
		ExpectErr bool
	}{
		{"tag1", 1, false},
		{"tag1", 1, true}, //Tests TagExist
	}

	for _, test := range tests {
		_, err := AddTag(db, test.Tag)

		if err != nil && !test.ExpectErr {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		if err == nil && test.ExpectErr {
			t.Errorf("Expected an error, but none was received")
		}

		//Case: expected error, received an error
		if err != nil && test.ExpectErr {
			continue
		}

		if err == nil && !test.ExpectErr {
			var count int64
			row := db.QueryRow("SELECT COUNT(*) FROM tags")
			err := row.Scan(&count)
			if err != nil {
				t.Errorf("Error happened when trying to obtain count of tags")
			}

			if count != test.Count {
				t.Errorf("Expected %d tags, but got %d", test.Count, count)
			}
		}
	}
}

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

func TestCreate(t *testing.T) {
	file := "./testing/create_test_file.db"

	//Need to create the test db file
	_, err := os.Create(file)
	if err != nil {
		t.Errorf("Unexpected error when create the database: %s", err.Error())
	}

	if !Exist(file) {
		t.Errorf("File: %s does not exist", file)
	}

	err = os.Remove(file)
	if err != nil {
		t.Errorf("Error while removing file (%s): %s", file, err.Error())
	}
}

func TestInit(t *testing.T) {
	file := "./testing/init_test_file.db"
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
}
