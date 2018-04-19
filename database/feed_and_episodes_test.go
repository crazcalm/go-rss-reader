package database

import (
	"testing"
	"time"
)

func TestGetFeedEpisodeIDs(t *testing.T) {
	file := "./testing/get_feed_episode_ids.db"
	db := createTestDB(file)
	feedURL := "get_feed_episode_ids.com"
	totalEpisodes := 2

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Failed to add feed to database: %s", err.Error())
	}

	episodeID1, err := AddEpisode(db, feedID, "ep1.com", "ep1 title", time.Now(), "ep1 raw data")
	if err != nil {
		t.Errorf("failed to add episode to database: %s", err.Error())
	}

	episodeID2, err := AddEpisode(db, feedID, "ep2.com", "ep2 title", time.Now(), "ep2 raw data")
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

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Failed to add feed to database: %s", err.Error())
	}

	episodeID1, err := AddEpisode(db, feedID, "ep1.com", "ep1 title", time.Now(), "ep1 raw data")
	if err != nil {
		t.Errorf("failed to add episode to database: %s", err.Error())
	}

	_, err = AddEpisode(db, feedID, "ep2.com", "ep2 title", time.Now(), "ep2 raw data")
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
