package service

import (
	"fmt"

	"github.com/resnostyle/mqttkit/env"
)

// Settings is the runtime configuration for this poller.
// Add domain-specific fields here when copying the bootstrap.
type Settings struct {
	env.MQTT
	PollIntervalSeconds int
}

func FromEnv() (Settings, error) {
	mqtt, err := env.LoadMQTT("home/example", "mqtt-bootstrap")
	if err != nil {
		return Settings{}, err
	}
	poll, err := env.Int("POLL_INTERVAL_SECONDS", 300)
	if err != nil {
		return Settings{}, err
	}
	if poll < 5 {
		return Settings{}, fmt.Errorf("POLL_INTERVAL_SECONDS must be >= 5")
	}
	return Settings{
		MQTT:                mqtt,
		PollIntervalSeconds: poll,
	}, nil
}
