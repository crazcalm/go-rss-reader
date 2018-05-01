package database

import (
	"strings"
	"testing"
	"time"
)

const (
	testRawData = `<!--                    
	Description: rss channel link                   
	-->                     
	<rss version="2.0">     
	  <channel>             
	    <link>http://example.org</link>             
	  </channel>            
	</rss> `
)

func TestUpdateEpisodeAuthor(t *testing.T) {
	file := "./testing/update_episode_author.db"
	db := createTestDB(file)
	feedURL := "update_episode_author.com"
	episodeURL := "update_episode_author.com/1"
	episodeTitle := "Episode Title"
	rawData := testRawData
	date := time.Now()
	authorName := "Jane Doe"
	authorEmail := "jane.doe@gmail.com"

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error happened when adding a feed to the database: %s", err.Error())
	}

	episodeID, err := AddEpisode(db, feedID, episodeURL, episodeTitle, &date, rawData)
	if err != nil {
		t.Errorf("Failed to add an episode to the database: %s", err.Error())
	}

	authorID, err := AddAuthor(db, authorName, authorEmail)
	if err != nil {
		t.Errorf("Failed to add an auhtor to the database: %s", err.Error())
	}

	//Actual test
	err = UpdateEpisodeAuthor(db, episodeID, authorID)
	if err != nil {
		t.Errorf("Error while trying to update the author of an episode: %s", err.Error())
	}
}

func TestUpdateEpisodeMediaContent(t *testing.T) {
	file := "./testing/update_episode_media_content.db"
	db := createTestDB(file)
	feedURL := "update_episode_media_content.com"
	episodeURL := "update_episode_media_content.com/1"
	episodeTitle := "Episode Title"
	rawData := testRawData
	date := time.Now()
	mediaContent := "media content"

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error happened when adding a feed to the database: %s", err.Error())
	}

	episodeID, err := AddEpisode(db, feedID, episodeURL, episodeTitle, &date, rawData)
	if err != nil {
		t.Errorf("Failed to add an episode to the database: %s", err.Error())
	}

	//Actual test
	err = UpdateEpisodeMediaContent(db, episodeID, mediaContent)
	if err != nil {
		t.Errorf("Error while trying to update the media content of an episode: %s", err.Error())
	}
}

func TestEpisodeExist(t *testing.T) {
	file := "./testing/episode_exist.db"
	db := createTestDB(file)
	feedURL := "episode_exist.com"
	episodeURL := "episode_exist.com/1"
	episodeTitle := "Episode Title"
	notEpisodeTitle := "Not the TITLE!!!"
	rawData := testRawData
	date := time.Now()

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error happened when adding a feed to the database: %s", err.Error())
	}

	_, err = AddEpisode(db, feedID, episodeURL, episodeTitle, &date, rawData)
	if err != nil {
		t.Errorf("Failed to add an episode to the database: %s", err.Error())
	}

	//Actual test
	if !EpisodeExist(db, episodeTitle) {
		t.Errorf("Expected an episode titled %s to exist in the database", episodeTitle)
	}

	if EpisodeExist(db, notEpisodeTitle) {
		t.Errorf("Did not expect an episode titled %s to exist in the database", notEpisodeTitle)
	}
}

func TestMarkEpisodeAsSeen(t *testing.T) {
	file := "./testing/mark_episode_as_seen.db"
	db := createTestDB(file)
	feedURL := "mark_episode_as_seen.com"
	episodeURL := "mark_episode_as_seen.com/1"
	title := "Episode title"
	rawData := testRawData
	date := time.Now()

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error happened when adding a feed to the database: %s", err.Error())
	}

	episodeID, err := AddEpisode(db, feedID, episodeURL, title, &date, rawData)
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
	rawData := testRawData
	date := time.Now()

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error happened when adding a feed to the database: %s", err.Error())
	}

	episodeID, err := AddEpisode(db, feedID, episodeURL, title, &date, rawData)
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
	rawData := testRawData
	date := time.Now()

	feedID, err := AddFeedURL(db, feedURL)
	if err != nil {
		t.Errorf("Error happened when adding a feed to the database: %s", err.Error())
	}

	_, err = AddEpisode(db, feedID, episodeURL, title, &date, rawData)
	if err != nil {
		t.Errorf("Failed to add an episode to the database: %s", err.Error())
	}

}
