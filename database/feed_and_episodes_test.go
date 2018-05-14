package database

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestGetEpisodeHeaderData(t *testing.T) {
	testName := "get_episode_header_data"
	file := fmt.Sprintf("./testing/%s.db", testName)
	db := createTestDB(file)

	feedName := "feed name"
	feedName2 := "feed name 2"
	feedURL := fmt.Sprintf("%s.com", testName)
	feedURL2 := fmt.Sprintf("%s.zh", testName)
	epTitle1 := "episode title 1"
	epURL1 := fmt.Sprintf("%s/1", feedURL)
	epTitle2 := "episode title 2"
	epURL2 := fmt.Sprintf("%s/2", feedURL)
	epTitle3 := "episode title 3"
	epURL3 := fmt.Sprintf("%s/3", feedURL)
	feedAuthorName := "feed author name"
	feedAuthorEmail := "feed author email"
	epAuthorName := "episode author name"
	epAuthorEmail := "episode author email"
	mediaContent := "media content"
	rawData := "raw Data"
	date := time.Now()

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Error(err)
	}

	feedID2, err := AddFeedURL(db, feedURL2)
	if err != nil {
		t.Error(err)
	}

	err = UpdateFeedTitle(db, feedID, feedName)
	if err != nil {
		t.Error(err)
	}

	err = UpdateFeedTitle(db, feedID2, feedName2)
	if err != nil {
		t.Error(err)
	}

	feedAuthorID, err := AddAuthor(db, feedAuthorName, feedAuthorEmail)
	if err != nil {
		t.Error(err)
	}

	epAuthorID, err := AddAuthor(db, epAuthorName, epAuthorEmail)
	if err != nil {
		t.Error(err)
	}

	err = UpdateFeedAuthor(db, feedID, feedAuthorID)
	if err != nil {
		t.Error(err)
	}

	epID1, err := AddEpisode(db, feedID, epURL1, epTitle1, &date, rawData)
	if err != nil {
		t.Errorf("Failed to add an episode to the database: %s", err.Error())
	}

	err = UpdateEpisodeAuthor(db, epID1, epAuthorID)
	if err != nil {
		t.Error(err)
	}

	epID2, err := AddEpisode(db, feedID, epURL2, epTitle2, &date, rawData)
	if err != nil {
		t.Errorf("Failed to add an episode to the database: %s", err.Error())
	}

	err = UpdateEpisodeMediaContent(db, epID2, mediaContent)
	if err != nil {
		t.Error(err)
	}

	epID3, err := AddEpisode(db, feedID2, epURL3, epTitle3, &date, rawData)
	if err != nil {
		t.Errorf("Failed to add an episode to the database: %s", err.Error())
	}

	tests := []struct {
		FeedID       int64
		EpisodeID    int64
		FeedName     string
		EpTitle      string
		Author       string
		EpLink       string
		MediaContent string
	}{
		{feedID, epID1, feedName, epTitle1, "episode author name (episode author email)", epURL1, ""},
		{feedID, epID2, feedName, epTitle2, "feed author name (feed author email)", epURL2, mediaContent},
		{feedID2, epID3, feedName2, epTitle3, "", epURL3, ""},
	}

	for _, test := range tests {
		dbFeedName, dbEpTitle, dbAuthor, dbEpLink, _, dbMediaContent, err := GetEpisodeHeaderData(db, test.FeedID, test.EpisodeID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		if !strings.EqualFold(dbFeedName, test.FeedName) {
			t.Errorf("(epID = %d) Expected %s, but got %s", test.EpisodeID, test.FeedName, dbFeedName)

		}
		if !strings.EqualFold(dbEpTitle, test.EpTitle) {
			t.Errorf("(epID = %d) Expected %s, but got %s", test.EpisodeID, test.EpTitle, dbEpTitle)
		}

		if !strings.EqualFold(dbAuthor, test.Author) {
			t.Errorf("(epID = %d) Expected %s, but got %s", test.EpisodeID, test.Author, dbAuthor)
		}

		if !strings.EqualFold(dbEpLink, test.EpLink) {
			t.Errorf("(epID = %d) Expected %s, but got %s", test.EpisodeID, test.EpLink, dbEpLink)
		}

		if !strings.EqualFold(dbMediaContent, test.MediaContent) {
			t.Errorf("(epID = %d) Expected %s, but got %s", test.EpisodeID, dbMediaContent, test.MediaContent)
		}
	}
}

func TestGetEpisodeIDByFeedIDAndTitle(t *testing.T) {
	file := "./testing/get_episode_id_by_feed_id_and_title.db"
	db := createTestDB(file)
	feedURL := "get_episode_id_by_feed_id_and_title.com"
	episodeTitle := "Episode Title"
	date := time.Now()

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Failed to add feed to database: %s", err.Error())
	}

	episodeID, err := AddEpisode(db, feedID, "ep1.com", episodeTitle, &date, "ep raw data")
	if err != nil {
		t.Errorf("failed to add episode to database: %s", err.Error())
	}

	//Actual test
	dbID, err := GetEpisodeIDByFeedIDAndTitle(db, feedID, episodeTitle)
	if err != nil {
		t.Errorf("Error happened while trying to get the episode id by feed id and episode title (%s): %s", episodeTitle, err.Error())
	}

	if episodeID != dbID {
		t.Errorf("Expected episode id to be %d, but got %d", episodeID, dbID)
	}

	//Part 2: partial title
	dbID2, err := GetEpisodeIDByFeedIDAndTitle(db, feedID, "Ep")
	if err != nil {
		t.Errorf("Error happened while trying to get the episode id by feed id and episode title (%s): %s", episodeTitle, err.Error())
	}

	if episodeID != dbID2 {
		t.Errorf("Expected episode id to be %d, but got %d", episodeID, dbID2)
	}
}

func TestGetFeedEpisodeIDs(t *testing.T) {
	file := "./testing/get_feed_episode_ids.db"
	db := createTestDB(file)
	feedURL := "get_feed_episode_ids.com"
	totalEpisodes := 2
	date := time.Now()

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Failed to add feed to database: %s", err.Error())
	}

	episodeID1, err := AddEpisode(db, feedID, "ep1.com", "ep1 title", &date, "ep1 raw data")
	if err != nil {
		t.Errorf("failed to add episode to database: %s", err.Error())
	}

	episodeID2, err := AddEpisode(db, feedID, "ep2.com", "ep2 title", &date, "ep2 raw data")
	if err != nil {
		t.Errorf("failed to add episode to database: %s", err.Error())
	}

	//Actual test
	ids, err := GetFeedEpisodeIDs(db, feedID)
	if err != nil {
		t.Errorf("Failed to get episode ids for a feed: %s", err.Error())
	}

	if len(ids) != totalEpisodes {
		t.Errorf("Expected %d episodes, but got %d", totalEpisodes, len(ids))
	}

	if ids[0] != episodeID1 || ids[1] != episodeID2 {
		t.Errorf("Expected the episode ids to be %d and %d, but got %d and %d", episodeID1, episodeID2, ids[0], ids[1])
	}
}

func TestGetFeedEpisodeSeenRatio(t *testing.T) {
	file := "./testing/get_feed_episode_seen_ratio.db"
	db := createTestDB(file)
	feedURL := "get_feed_episode_seen_ratio.com"
	var seen int64 = 1
	var total int64 = 2
	date := time.Now()

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Failed to add feed to database: %s", err.Error())
	}

	episodeID1, err := AddEpisode(db, feedID, "ep1.com", "ep1 title", &date, "ep1 raw data")
	if err != nil {
		t.Errorf("failed to add episode to database: %s", err.Error())
	}

	_, err = AddEpisode(db, feedID, "ep2.com", "ep2 title", &date, "ep2 raw data")
	if err != nil {
		t.Errorf("failed to add episode to database: %s", err.Error())
	}

	err = MarkEpisodeAsSeen(db, episodeID1)
	if err != nil {
		t.Errorf("Failed to mark episode as seen: %s", err.Error())
	}

	//Actual test
	dbSeen, dbTotal, err := GetFeedEpisodeSeenRatio(db, feedID)
	if err != nil {
		t.Errorf("Failed to get the episodes seen ration from the database: %s", err.Error())
	}

	if seen != dbSeen || total != dbTotal {
		t.Errorf("Expected %d seen and %d total, but got %d and %d", seen, total, dbSeen, dbTotal)
	}

}
