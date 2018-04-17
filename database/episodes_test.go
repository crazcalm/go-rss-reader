package database

import (
	"strings"
	"testing"
	"time"
)

func TestMarkEpisodeAsSeen(t *testing.T) {
	file := "./testing/mark_episode_as_seen.db"
	db := createTestDB(file)
	feedURL := "mark_episode_as_seen.com"
	episodeURL := "mark_episode_as_seen.com/1"
	title := "Episode title"
	rawData := "Episode Raw Data"
	date := time.Now()

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error happened when adding a feed to the database: %s", err.Error())
	}

	episodeID, err := AddEpisode(db, feedID, episodeURL, title, date, rawData)
	if err != nil {
		t.Errorf("Failed to add an episode to the database: %s", err.Error())
	}

	//Actual test
	err = MarkEpisodeAsSeen(db, episodeID)
	if err != nil {
		t.Errorf("Error happened while trying to mark an episode as seen: %s", err.Error())
	}

	_, _, _, dbSeen, _, err := GetEpisode(db, episodeID)
	if err != nil {
		t.Errorf("Error happened when trying to get an episode from the database: %s", err.Error())
	}

	if dbSeen != 1 {
		t.Errorf("Expected seen to be %d, but got %d", 1, dbSeen)
	}
}

func TestGetEpisode(t *testing.T) {
	file := "./testing/get_episode.db"
	db := createTestDB(file)
	feedURL := "get_episode.com"
	episodeURL := "get_episode.com/1"
	title := "Episode Title"
	rawData := "Episode RawData"
	date := time.Now()

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error happened when adding a feed to the database: %s", err.Error())
	}

	episodeID, err := AddEpisode(db, feedID, episodeURL, title, date, rawData)
	if err != nil {
		t.Errorf("Failed to add an episode to the database: %s", err.Error())
	}

	//Actual test
	dbURL, dbTitle, dbDate, dbSeen, dbRawData, err := GetEpisode(db, episodeID)
	if err != nil {
		t.Errorf("Error happened when trying to get an episode from the database: %s", err.Error())
	}

	if !strings.EqualFold(dbURL, episodeURL) {
		t.Errorf("Expected %s, but got %s", episodeURL, dbURL)
	}

	if !strings.EqualFold(dbTitle, title) {
		t.Errorf("Expected %s, but got %s", title, dbTitle)
	}

	if !strings.EqualFold(dbRawData, rawData) {
		t.Errorf("Expected %s, but got %s", rawData, dbRawData)
	}

	if dbSeen != 0 {
		t.Errorf("Expected seen to be %d, but got %d", 0, dbSeen)
	}

	if !date.Equal(dbDate) {
		timeFormat := "Mon Jan 2 15:04:05 -0700 MST 2006"
		t.Errorf("Expected %s, but got %s", date.Format(timeFormat), dbDate.Format(timeFormat))
	}

}

func TestAddEpisode(t *testing.T) {
	file := "./testing/add_epiode.db"
	db := createTestDB(file)
	feedURL := "add_episode.com"
	episodeURL := "add_episode.com/1"
	title := "Episode Title"
	rawData := "Episode Raw Data"
	date := time.Now()

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error happened when adding a feed to the database: %s", err.Error())
	}

	_, err = AddEpisode(db, feedID, episodeURL, title, date, rawData)
	if err != nil {
		t.Errorf("Failed to add an episode to the database: %s", err.Error())
	}

}
