package auth

import "time"

type Sudo struct {
	ID             int
	SessionID      int
	DurationInSecs int
	CreatedAt      time.Time
}

func NewSudo(sessionID int, durationInSecs int) (sudo Sudo) {
	return Sudo{
		SessionID:      sessionID,
		DurationInSecs: durationInSecs,
		CreatedAt:      time.Now(),
	}
}
