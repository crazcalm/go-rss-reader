package database

import (
	"testing"
)

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
