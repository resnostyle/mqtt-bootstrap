package service

import "time"

// Snapshot is the retained JSON published to MQTT.
// Replace fields with your domain payload.
type Snapshot struct {
	OK        bool      `json:"ok"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updated_at"`
}

func BuildSnapshot(now time.Time) Snapshot {
	return Snapshot{
		OK:        true,
		Message:   "bootstrap idle — replace Fetch with real work",
		UpdatedAt: now.UTC(),
	}
}
