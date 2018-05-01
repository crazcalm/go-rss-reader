package database

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetFeedInfo(t *testing.T) {
	file := "./testing/get_feed_info.db"
	db := createTestDB(file)

	tests := []struct {
		FeedURL string
	}{
		{"http://www.leoville.tv/podcasts/sn.xml"},
	}

	for _, test := range tests {
		feedID, err := AddFeedURL(db, test.FeedURL)
		if err != nil {
			t.Errorf("Error while inserting a feed into the database: %s", err.Error())
		}

		err = GetFeedInfo(db, feedID)
		if err != nil {
			t.Errorf("Error while trying to get the feed info for url (%s): %s", test.FeedURL, err.Error())
		}
	}
}

func TestLoadFeed(t *testing.T) {
	file := "./testing/load_feed.db"
	db := createTestDB(file)
	feedURL := "load_feed.com"
	rawData := testRawData
	defer db.Close()

	feedID1, err := AddFeedURL(db, fmt.Sprintf("%s1", feedURL))
	if err != nil {
		t.Errorf("Error while inserting a feed into the database: %s", err.Error())
	}

	tagID, err := AddTag(db, "tag")
	if err != nil {
		t.Errorf("Error happened while adding a tag to the database: %s", err.Error())
	}

	_, err = AddTagToFeed(db, feedID1, tagID)
	if err != nil {
		t.Errorf("Error happened while adding a feed tag to the database: %s", err.Error())
	}

	feedID2, err := AddFeedURL(db, fmt.Sprintf("%s2", feedURL))
	if err != nil {
		t.Errorf("Error while inserting a feed into the database: %s", err.Error())
	}

	err = UpdateFeedRawData(db, feedID1, rawData)
	if err != nil {
		t.Errorf("Error happened while trying to update the raw_data for a feed: %s", err.Error())
	}

	tests := []struct {
		FeedID int64
	}{
		{feedID1},
		{feedID2},
	}

	for _, test := range tests {
		_, err := LoadFeed(db, test.FeedID)
		if err != nil {
			t.Errorf("Error happened when trying to load a feed (id = %d) from the database: %s", test.FeedID, err.Error())
		}
	}

}

func TestGetFeedRawData(t *testing.T) {
	file := "./testing/get_feed_raw_data.db"
	db := createTestDB(file)
	feedURL := "get_feed_raw_data.com"
	rawData := testRawData

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error while inserting a feed into the database: %s", err.Error())
	}

	err = UpdateFeedRawData(db, feedID, rawData)
	if err != nil {
		t.Errorf("Error happened while trying to update the raw_data for a feed: %s", err.Error())
	}

	//Actual test
	dbRawData, err := GetFeedRawData(db, feedID)
	if err != nil {
		t.Errorf("Failed to get the raw_data of a feed: %s", err.Error())
	}

	if !strings.EqualFold(rawData, dbRawData) {
		t.Errorf("Expected the raw data to be %s, but get %s", rawData, dbRawData)
	}

}

func TestUpdateFeedRawData(t *testing.T) {
	file := "./testing/update_feed_raw_data.db"
	db := createTestDB(file)
	feedURL := "update_feed_raw_data.com"
	rawData := testRawData

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error while inserting a feed into the database: %s", err.Error())
	}

	//Actual test
	err = UpdateFeedRawData(db, feedID, rawData)
	if err != nil {
		t.Errorf("Error happened while trying to update the raw_data for a feed: %s", err.Error())
	}
}

func TestGetFeedURL(t *testing.T) {
	file := "./testing/get_feed_url.db"
	db := createTestDB(file)
	feedURL := "get_feed_url.com"

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error while inserting a feed into the database: %s", err.Error())
	}

	dbURL, err := GetFeedURL(db, feedID)
	if err != nil {
		t.Errorf("Failed while trying to get the url for a feed: %s", err.Error())
	}

	if !strings.EqualFold(feedURL, dbURL) {
		t.Errorf("Expected the url to be %s, but got %s", feedURL, dbURL)
	}
}

func TestGetFeedTitle(t *testing.T) {
	file := "./testing/get_feed_title.db"
	db := createTestDB(file)
	feedURL := "get_feed_title.com"
	title := "Feed Title"

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error while inserting a feed into the database: %s", err.Error())
	}

	err = UpdateFeedTitle(db, feedID, title)
	if err != nil {
		t.Errorf("Error happened while trying to update the title for a feed: %s", err.Error())
	}

	//Actual test
	dbTitle, err := GetFeedTitle(db, feedID)
	if err != nil {
		t.Errorf("Failed to get feed title from database: %s", err.Error())
	}

	if !strings.EqualFold(title, dbTitle) {
		t.Errorf("Expected feed title to be %s, but got %s", title, dbTitle)
	}
}

func TestUpdateFeedTitle(t *testing.T) {
	file := "./testing/update_feed_title.db"
	db := createTestDB(file)
	feedURL := "update_feed_title.com"
	title := "Feed Title"

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error while inserting a feed into the database: %s", err.Error())
	}

	//Actual test
	err = UpdateFeedTitle(db, feedID, title)
	if err != nil {
		t.Errorf("Error happened while trying to update the title for a feed: %s", err.Error())
	}
}

func TestGetFeedAuthorID(t *testing.T) {
	file := "./testing/get_feed_author_id.db"
	db := createTestDB(file)

	authorName := "John Doe"
	authorEmail := "john.doe@gmail.com"
	authorID, err := AddAuthor(db, authorName, authorEmail)
	if err != nil {
		t.Errorf("Failed to add author: %s", err.Error())
	}

	feedURL := "update_feed_author.com"
	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error while inserting a feed into the database: %s", err.Error())
	}

	err = UpdateFeedAuthor(db, feedID, authorID)
	if err != nil {
		t.Errorf("Failed to update the author of a feed: %s", err.Error())
	}

	//Actual test
	dbID, err := GetFeedAuthorID(db, feedID)
	if err != nil {
		t.Errorf("Failed to get the author id associate with a feed: %s", err.Error())
	}

	if dbID != authorID {
		t.Errorf("Expected authorID to be %d, but got %d", authorID, dbID)
	}
}

func TestUpdateFeedAuthor(t *testing.T) {
	file := "./testing/update_feed_author.db"
	db := createTestDB(file)

	authorName := "John Doe"
	authorEmail := "john.doe@gmail.com"
	authorID, err := AddAuthor(db, authorName, authorEmail)
	if err != nil {
		t.Errorf("Failed to add author: %s", err.Error())
	}

	feedURL := "update_feed_author.com"
	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error while inserting a feed into the database: %s", err.Error())
	}

	//Acutal test
	err = UpdateFeedAuthor(db, feedID, authorID)
	if err != nil {
		t.Errorf("Failed to update the author of a feed: %s", err.Error())
	}
}

func TestGetFeedID(t *testing.T) {
	file := "./testing/get_feed_id.db"
	db := createTestDB(file)
	feedURL := "get_feed_id.com"
	feedTitle := "Feed Title"
	tests := []string{feedURL, feedTitle}

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error while inserting a feed into the database: %s", err.Error())
	}

	err = UpdateFeedTitle(db, feedID, feedTitle)
	if err != nil {
		t.Errorf("Error while trying to update the feed title")
	}

	for _, test := range tests {
		result, err := GetFeedID(db, test)
		if err != nil {
			t.Errorf("Error while trying to get the feed id for a url/title: %s", err.Error())
		}

		if feedID != result {
			t.Errorf("Expected the feed id to be %d, but got %d", feedID, result)
		}
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

func TestUndeleteFeed(t *testing.T) {
	file := "./testing/undelete_feed.db"
	db := createTestDB(file)
	var feed int64 = 2
	var expected int64 = 4 // 5 feeds were created, but one was deleted

	for i := 0; i < 5; i++ {
		_, err := AddFeedURL(db, fmt.Sprintf("url%d", i))
		if err != nil {
			t.Errorf("Error while inserting feeds into the database: %s", err.Error())
		}
	}

	err := DeleteFeed(db, feed)
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

	expected++ // Undeleted the previously deleted feed
	err = UndeleteFeed(db, feed)
	if err != nil {
		t.Errorf("Error happened while undeleting a feed: %s", err.Error())
	}

	row = db.QueryRow("SELECT COUNT(*) FROM feeds WHERE deleted = 0")
	err = row.Scan(&count)
	if err != nil {
		t.Errorf("Error happened when trying to obtain count of feeds: %s", err.Error())
	}

	if count != expected {
		t.Errorf("Expected %d feeds, but got %d", expected, count)
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
