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
		log.Fatalf("Error when trying to create the database (%s): %s", file, err.Error())
	}
	return db
}

func TestAllActiveFeedTags(t *testing.T) {
	file := "./testing/all_active_feed_tags.db"
	db := createTestDB(file)

	//Variables
	numOfTags := 5
	var feedID int64
	feedURL := "all_active_feed_tags.com"
	var tagsAddedToFeed = make(map[int64]string)

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error when adding feed to the database: %s", err.Error())
	}

	for i := 1; i < numOfTags; i++ {
		tag := fmt.Sprintf("tag%s", i)
		tagID, err := AddTag(db, tag) // Adding tag to the database
		if err != nil {
			t.Errorf("Error happaned while adding a tag: %s", err.Error())
		}

		if i > 2 {
			_, err := AddTagToFeed(db, feedID, tagID)
			if err != nil {
				t.Errorf("Error happened while adding tag (%s) to a feed: %s", tag, err.Error())
			}

			//Adding tag to list
			tagsAddedToFeed[tagID] = tag
		}
	}

	//Actual test part
	for keyID, value := range tagsAddedToFeed {
		if !FeedHasTag(db, feedID, keyID) {
			t.Errorf("Feed(%d) is missing tag (%s)", feedID, value)
		}
	}
}

func TestDeleteAllTagsFromFeed(t *testing.T) {
	file := "./testing/delete_all_tags_from_feed.db"
	db := createTestDB(file)
	var tags []string
	var tagIDs []int64
	numOfTags := 5
	feedURL := "delete_all_tags_from_feed.com"

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error when adding feed to the database: %s", err.Error())
	}

	for i := 1; i < numOfTags; i++ {
		tag := fmt.Sprintf("tag%d", i)
		tagID, err := AddTag(db, tag) // Adding tag to the database
		if err != nil {
			t.Errorf("Error happened while adding a tag: %s", err.Error())
		}
		tagIDs = append(tagIDs, tagID)

		//Adding tag to feed
		_, err = AddTagToFeed(db, feedID, tagID)
		if err != nil {
			t.Errorf("Error when adding tag to feed: %s", err.Error())
		}
		tags = append(tags, tag)
	}

	err = DeleteAllTagsFromFeed(db, feedID)
	if err != nil {
		t.Errorf("Error while deleting all the feeds from a tag: %s", err.Error())
	}

	for index, tagID := range tagIDs {
		if !FeedHasTag(db, feedID, tagID) {
			t.Errorf("Feed not expected to have tag: %s", tags[index])
		}
	}

}

func TestFilterFeeds(t *testing.T) {
	file := "./testing/filter_feeds.db"
	db := createTestDB(file)
	var feeds = make(map[int64]string)
	expected := 3

	for i := 1; i <= 5; i++ {
		url := fmt.Sprintf("url%d", i)
		id, err := AddFeedURL(db, url)
		if err != nil {
			t.Errorf("Error while inserting feeds into the database: %s", err.Error())
		}
		feeds[id] = url
	}

	var count int
	for key := range feeds {
		delete(feeds, key)
		count++
		if count == 3 {
			break
		}
	}

	result := FilterFeeds(db, feeds)

	if len(result) != expected {
		t.Errorf("Expected %d feeds, but got %d feeds", expected, len(result))
	}

}

func TestDeleteFeed(t *testing.T) {
	file := "./testing/delete_feed.db"
	db := createTestDB(file)
	var feedToDelete int64 = 2
	var expected int64 = 4

	for i := 0; i < 5; i++ {
		_, err := AddFeedURL(db, fmt.Sprintf("url%d", i))
		if err != nil {
			t.Errorf("Error while inserting feeds into the database: %s", err.Error())
		}
	}

	err := DeleteFeed(db, feedToDelete)
	if err != nil {
		t.Errorf("Error while deleting the feed: %s", err.Error())
	}

	var count int64
	row := db.QueryRow("SELECT COUNT(*) FROM feeds WHERE deleted = 0")
	err = row.Scan(&count)
	if err != nil {
		t.Errorf("Error happened when trying to obtain count of feeds: %s", err.Error())
	}

	if count != expected {
		t.Errorf("Expected %d feeds, but got %d", expected, count)
	}

}

func TestAllActiveFeeds(t *testing.T) {
	file := "./testing/all_active_feeds.db"
	db := createTestDB(file)
	var expected int64 = 5

	for i := 0; i < 5; i++ {
		_, err := AddFeedURL(db, fmt.Sprintf("url%d", i))
		if err != nil {
			t.Errorf("Error while inserting feeds into the database: %s", err.Error())
		}
	}

	_, err := db.Exec("INSERT INTO feeds (uri, deleted) VALUES ($1, $2)", "deletedURL", 1)
	if err != nil {
		t.Errorf("Error during test setup: %s", err.Error())
	}

	var count int64
	row := db.QueryRow("SELECT COUNT(*) FROM feeds WHERE deleted = 0")
	err = row.Scan(&count)
	if err != nil {
		t.Errorf("Error happened when trying to obtain count of feeds: %s", err.Error())
	}

	if count != expected {
		t.Errorf("Expected %d feeds, but got %d", expected, count)
	}
}

func TestAddTagToFeed(t *testing.T) {
	file := "./testing/add_tag_to_feed.db"
	db := createTestDB(file)

	feedID, err := AddFeedURL(db, "url1")
	if err != nil {
		t.Errorf("Error happened when adding a feed to the database")
	}

	tagID, err := AddTag(db, "tag1")
	if err != nil {
		t.Errorf("Error happened when adding a tag to the database")
	}

	tests := []struct {
		FeedID    int64
		TagID     int64
		ExpectErr bool
		Count     int64
	}{
		{feedID, tagID, false, 1},
		{feedID, tagID, true, 1},
	}

	for _, test := range tests {
		_, err := AddTagToFeed(db, test.FeedID, test.TagID)

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
			row := db.QueryRow("SELECT COUNT(*) FROM feeds_and_tags")
			err := row.Scan(&count)
			if err != nil {
				t.Errorf("Error happened when trying to obtain count of feeds_and_tags")
			}

			if count != test.Count {
				t.Errorf("Expected the count to be %d, but got %d", test.Count, count)
			}
		}
	}

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
