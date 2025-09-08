package session

import (
	"database/sql"
	"time"
)

// Session represents a user session or privileged action
type Session struct {
	ID        int
	UserID    int
	Action    string
	Timestamp time.Time
}

// LogSession records a user action into the sessions table
func LogSession(db *sql.DB, userID int, action string) error {
	_, err := db.Exec(
		"INSERT INTO sessions(user_id, action, timestamp) VALUES($1, $2, $3)",
		userID, action, time.Now(),
	)
	return err
}

// GetUserSessions retrieves all sessions for a specific user
func GetUserSessions(db *sql.DB, userID int) ([]Session, error) {
	rows, err := db.Query("SELECT id, user_id, action, timestamp FROM sessions WHERE user_id=$1 ORDER BY timestamp DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.UserID, &s.Action, &s.Timestamp); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}

// GetAllSessions retrieves all sessions (for reporting)
func GetAllSessions(db *sql.DB) ([]Session, error) {
	rows, err := db.Query("SELECT id, user_id, action, timestamp FROM sessions ORDER BY timestamp DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.UserID, &s.Action, &s.Timestamp); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}
