package main

import (
	"log/slog"
	"os"

	"github.com/resnostyle/mqtt-bootstrap/internal/service"
	"github.com/resnostyle/mqttkit/logx"
	"github.com/resnostyle/mqttkit/poll"
)

func main() {
	settings, err := service.FromEnv()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	logx.Configure(settings.LogLevel, true)

	slog.Info("starting mqtt-bootstrap",
		"poll_interval", settings.PollIntervalSeconds,
		"mqtt", settings.MQTTHost,
		"port", settings.MQTTPort,
		"prefix", settings.MQTTTopicPrefix,
		"discovery", settings.MQTTDiscoveryEnabled,
	)

	ctx, cancel := poll.NotifyContext()
	defer cancel()

	if err := service.Run(ctx, settings); err != nil {
		slog.Error("service stopped", "err", err)
		os.Exit(1)
	}
}
