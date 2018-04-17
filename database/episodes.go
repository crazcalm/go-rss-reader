package database

import (
	"database/sql"
	"time"
)

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
func AddEpisode(db *sql.DB, feedID int64, url, title string, date time.Time, rawData string) (int64, error) {
	stmt := "INSERT INTO episodes (feed_id, uri, title, date, raw_data) VALUES ($1, $2, $3, $4, $5)"
	var result int64

	dbResult, err := db.Exec(stmt, feedID, url, title, date, rawData)
	if err != nil {
		return result, err
	}

	result, err = dbResult.LastInsertId()
	if err != nil {
		return result, err
	}

	return result, nil
}
