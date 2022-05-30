package auth

import "time"

type Sudo struct {
	ID             int
	SessionID      string
	DurationInSecs int
	CreatedAt      time.Time
}

func NewSudo(sessionID string, durationInSecs int) (sudo Sudo) {
	return Sudo{
		SessionID:      sessionID,
		DurationInSecs: durationInSecs,
		CreatedAt:      time.Now(),
	}
}
