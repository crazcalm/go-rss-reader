package database

import (
	"fmt"
	"testing"
)

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

func TestIsFeedDeleted(t *testing.T) {
	file := "./testing/is_feed_deleted.db"
	db := createTestDB(file)

	feedID1, err := AddFeedURL(db, "url1")
	if err != nil {
		t.Errorf("Error while adding feed to test database: %s", err.Error())
	}

	feedID2, err := AddFeedURL(db, "url2")
	if err != nil {
		t.Errorf("Error while adding feed to test database: %s", err.Error())
	}

	err = DeleteFeed(db, feedID2)
	if err != nil {
		t.Errorf("Error while trying to delete a feed: %s", err.Error())
	}

	tests := []struct {
		FeedID   int64
		Expected bool
	}{
		{feedID1, false},
		{feedID2, true},
	}

	for _, test := range tests {
		result := IsFeedDeleted(db, test.FeedID)
		if test.Expected != result {
			t.Errorf("IsFeedDeleted test expected %t, but got %t", test.Expected, result)
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
