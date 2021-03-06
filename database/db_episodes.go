package database

import (
	"database/sql"
	"log"
	"time"
)

//GetEpisodeAuthor -- returns the episode author
func GetEpisodeAuthor(db *sql.DB, episodeID int64) (name, email string, err error) {
	stmt := "SELECT authors.name, authors.email FROM episodes INNER JOIN authors ON authors.id = episodes.author_id WHERE episodes.id = $1"
	row := db.QueryRow(stmt, episodeID)
	err = row.Scan(&name, &email)
	return
}

//EpisodeHasAuthor -- returns true is an author id exists and false otherwise
func EpisodeHasAuthor(db *sql.DB, episodeID int64) (result bool) {
	var count int64
	row := db.QueryRow("SELECT COUNT(author_id) FROM episodes WHERE id = $1", episodeID)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		result = true
	}
	return
}

//UpdateEpisodeAuthor -- Updates the author associated with the episode
func UpdateEpisodeAuthor(db *sql.DB, episodeID, authorID int64) (err error) {
	_, err = db.Exec("UPDATE episodes SET author_id = $1 WHERE id = $2", authorID, episodeID)
	return
}

//GetEpisodeMediaContent -- Gets the episode's media content from the database
func GetEpisodeMediaContent(db *sql.DB, episodeID int64) (mediaContent string, err error) {
	stmt := "SELECT media_content FROM episodes WHERE id = $1"
	row := db.QueryRow(stmt, episodeID)
	err = row.Scan(&mediaContent)
	return
}

//EpisodeHasMediaContent -- returns true is an author id exists and false otherwise
func EpisodeHasMediaContent(db *sql.DB, episodeID int64) (result bool) {
	var count int64
	row := db.QueryRow("SELECT COUNT(media_content) FROM episodes WHERE id = $1", episodeID)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		result = true
	}
	return
}

//UpdateEpisodeMediaContent -- updates the media content to the database
func UpdateEpisodeMediaContent(db *sql.DB, episodeID int64, mediaContent string) (err error) {
	_, err = db.Exec("UPDATE episodes SET media_content = $1 WHERE id = $2", mediaContent, episodeID)
	return
}

//EpisodeExist -- Based on the title, it checks if that episode already exists
func EpisodeExist(db *sql.DB, title string) (result bool) {
	var count int64
	row := db.QueryRow("SELECT COUNT(*) FROM episodes WHERE title = $1", title)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count != 0 {
		result = true
	}
	return
}

//MarkEpisodeAsSeen -- marks and episode as seen
func MarkEpisodeAsSeen(db *sql.DB, episodeID int64) (err error) {
	_, err = db.Exec("UPDATE episodes SET seen = 1 WHERE id = $1", episodeID)
	return
}

//GetEpisode -- gets an episode from the database
func GetEpisode(db *sql.DB, episodeID int64) (url, title string, date time.Time, seen int64, rawData string, err error) {
	stmt := "SELECT uri, title, date, seen, raw_data FROM episodes WHERE id = $1"
	row := db.QueryRow(stmt, episodeID)
	err = row.Scan(&url, &title, &date, &seen, &rawData)
	if err != nil {
		return
	}
	return
}

//AddEpisode -- adds an episode to the database
func AddEpisode(db *sql.DB, feedID int64, url, title string, date *time.Time, rawData string) (int64, error) {
	stmt := "INSERT INTO episodes (feed_id, uri, title, date, raw_data) VALUES ($1, $2, $3, $4, $5)"
	var result int64

	dbResult, err := db.Exec(stmt, feedID, url, title, date, rawData)
	if err != nil {
		return result, err
	}

	result, err = dbResult.LastInsertId()
	return result, err
}
