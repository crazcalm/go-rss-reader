package database

import (
	"fmt"
	"strings"
	"testing"
)

func TestFilterFeedTags(t *testing.T) {
	file := "./testing/filter_feed_tags.db"
	db := createTestDB(file)

	var feedID int64
	feedURL := "filter_feed_tags.com"
	var allTags = make(map[int64]string)
	var passedInTags = make(map[int64]string)
	var diffTags = make(map[int64]string)
	tagCount := 5
	passedInTagCount := 2 // must be less than tagCount

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error happened when adding a feed to the database")
	}

	for i := 1; i < tagCount; i++ {
		tag := fmt.Sprintf("tag%d", i)
		tagID, err := AddTag(db, tag)
		if err != nil {
			t.Errorf("Error happened while adding a tag the the database: %s", err.Error())
		}

		//Add tags to feed
		_, err = AddTagToFeed(db, feedID, tagID)
		if err != nil {
			t.Errorf("Error happened while adding a tag to a feed: %s", err.Error())
		}

		//Track all tags
		allTags[tagID] = tag

		//Passed in tags and diff tags
		if i < passedInTagCount {
			passedInTags[tagID] = tag
		} else {
			diffTags[tagID] = tag
		}
	}

	//Actual test
	results := FilterFeedTags(db, feedID, passedInTags)

	if len(results) != len(diffTags) {
		t.Errorf("Expected a total of %d tags, but got %d", len(diffTags), len(results))
	}

	for key, value := range diffTags {
		exist := false
		for resultKey, resultValue := range results {
			if key == resultKey && strings.EqualFold(value, resultValue) {
				exist = true
			}
		}

		if !exist {
			t.Errorf("Expected tag (%s) to be in the results, but it was not found", value)
		}
	}
}

func TestAddTagToFeed(t *testing.T) {
	file := "./testing/add_tag_to_feed.db"
	db := createTestDB(file)

	feedID, err := AddFeedURL(db, "url1")
	if err != nil {
		t.Errorf("Error happened while adding a feed to the database: %s", err.Error())
	}

	tagID, err := AddTag(db, "tag1")
	if err != nil {
		t.Errorf("Error happened while adding a tag to the database: %s", err.Error())
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

func TestDeleteTagFromFeed(t *testing.T) {
	file := "./testing/delete_tag_from_feed.db"
	db := createTestDB(file)

	var feedID int64
	feedURL := "delete_tag_from_feed.com"
	var tags = make(map[int64]string)
	numOfTags := 5

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error when adding feed to thedatabase: %s", err.Error())
	}

	for i := 1; i < numOfTags; i++ {
		tag := fmt.Sprintf("tag%d", i)
		tagID, err := AddTag(db, tag) // Adding tag to the database
		if err != nil {
			t.Errorf("Error happened while adding a tag: %s", err.Error())
		}

		_, err = AddTagToFeed(db, feedID, tagID)
		if err != nil {
			t.Errorf("Error occured while adding tag (%s) to feed: %s", tag, err.Error())
		}

		//Addind tag to the list
		tags[tagID] = tag
	}

	//Actual Test
	var count int

	for key := range tags {
		err := DeleteTagFromFeed(db, feedID, key)
		if err != nil {
			t.Errorf("Error occured while trying to delete a tag from the database: %s", err.Error())
		}
		delete(tags, key)
		count++

		if count >= 2 {
			break
		}
	}

	dbTags := AllActiveFeedTags(db, feedID)

	for keyID, value := range tags {
		result := false

		for dbID, dbValue := range dbTags {
			if dbID == keyID && strings.EqualFold(value, dbValue) {
				result = true
			}
		}

		if !result {
			t.Errorf("Expected feed to have tag (%s), but it was not found", value)
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
		tag := fmt.Sprintf("tag%d", i)
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
	dbTags := AllActiveFeedTags(db, feedID)

	for keyID, value := range tagsAddedToFeed {
		result := false

		for dbID, dbValue := range dbTags {
			if dbID == keyID && strings.EqualFold(value, dbValue) {
				result = true
			}
		}

		if !result {
			t.Errorf("Expected feed to have tag (%s), but it was not found", value)
		}
	}
}
