package auth

import (
	"time"
)

// Session represents the session of a user in a client.
type Session struct {
	ID int `json:"id,omitempty"`

	UserID int `json:"user_id,omitempty"`

	// First time that the user has been sign.
	LoggedAt time.Time `json:"logged_at,omitempty"`

	// Last time that the user perform a auth-required route call.
	LastSeenAt time.Time `json:"last_seen_at,omitempty"`

	// Platform which the user has been sign.
	LoggedWith string `json:"logged_with,omitempty"`

	// Session status
	Actived bool `json:"actived,omitempty"`
}

// NewSession initializes a new session instance
//  @param userID int: user id unique identifier.
//  @param loggedWith string: platform which the user has been sign.
//  @return session Session: new Session instance.
//	@return err error: session encryptation error.
func NewSession(userID int, loggedWith string) (session Session, err error) {
	now := time.Now()

	session = Session{
		UserID:     userID,
		LoggedWith: loggedWith,
		LoggedAt:   now,
		LastSeenAt: now,
		Actived:    true,
	}
	return
}
