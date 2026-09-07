package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/resnostyle/mqttkit/mqttpub"
	"github.com/resnostyle/mqttkit/poll"
)

// Fetch gathers domain data. Replace this with a real client call.
func Fetch(_ context.Context) (Snapshot, error) {
	return BuildSnapshot(time.Now()), nil
}

// Run connects MQTT, publishes discovery once, then polls until cancelled.
func Run(ctx context.Context, settings Settings) error {
	mqtt, err := mqttpub.New(
		settings.MQTTHost,
		settings.MQTTPort,
		settings.MQTTClientID,
		settings.MQTTUsername,
		settings.MQTTPassword,
		settings.MQTTTopicPrefix,
	)
	if err != nil {
		return err
	}
	defer mqtt.Close()

	if err := PublishDiscovery(settings, mqtt); err != nil {
		slog.Error("mqtt discovery publish failed", "err", err)
	}

	interval := time.Duration(settings.PollIntervalSeconds) * time.Second
	for {
		snap, err := Fetch(ctx)
		if err != nil {
			slog.Error("fetch failed", "err", err)
			snap = Snapshot{
				OK:        false,
				Message:   err.Error(),
				UpdatedAt: time.Now().UTC(),
			}
		}
		if err := PublishSnapshot(mqtt, snap); err != nil {
			slog.Error("publish failed", "err", err)
		}

		poll.Wait(ctx, interval)
		if ctx.Err() != nil {
			return nil
		}
	}
}
