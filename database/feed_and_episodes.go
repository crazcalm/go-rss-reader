package database

import (
	"database/sql"
)

//GetEpisodeIDByFeedIDAndTitle -- Gets the episodes using feed id and episode title
func GetEpisodeIDByFeedIDAndTitle(db *sql.DB, feedID int64, episodeTitle string) (id int64, err error) {
	stmt := "SELECT id FROM episodes WHERE feed_id = $1 AND title = $2"
	row := db.QueryRow(stmt, feedID, episodeTitle)
	err = row.Scan(&id)
	if err != nil {
		return
	}
	return
}

//GetFeedEpisodeIDs -- return ???
func GetFeedEpisodeIDs(db *sql.DB, feedID int64) (ids []int64, err error) {
	rows, err := db.Query("SELECT id FROM episodes WHERE feed_id = $1", feedID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	return
}

//GetFeedEpisodeSeenRatio -- return the seen over total episode seen for a feed.
func GetFeedEpisodeSeenRatio(db *sql.DB, feedID int64) (seen, total int64, err error) {
	row := db.QueryRow("SELECT COUNT(*) FROM episodes WHERE feed_id = $1", feedID)
	err = row.Scan(&total)
	if err != nil {
		return
	}
	row = db.QueryRow("SELECT COUNT(*) FROM episodes WHERE feed_id = $1 AND seen = 1", feedID)
	err = row.Scan(&seen)
	if err != nil {
		return
	}
	return
}
