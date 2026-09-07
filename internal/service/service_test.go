package service

import (
	"testing"
	"time"

	"github.com/resnostyle/mqttkit/mqttpub"
)

func TestFromEnvDefaults(t *testing.T) {
	t.Setenv("MQTT_HOST", "mqtt.example.test")
	t.Setenv("POLL_INTERVAL_SECONDS", "60")
	s, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if s.MQTTHost != "mqtt.example.test" {
		t.Fatalf("host %q", s.MQTTHost)
	}
	if s.PollIntervalSeconds != 60 {
		t.Fatalf("poll %d", s.PollIntervalSeconds)
	}
	if s.MQTTTopicPrefix != "home/example" {
		t.Fatalf("prefix %q", s.MQTTTopicPrefix)
	}
}

func TestFromEnvRejectsShortPoll(t *testing.T) {
	t.Setenv("POLL_INTERVAL_SECONDS", "1")
	if _, err := FromEnv(); err == nil {
		t.Fatal("expected error")
	}
}

func TestBuildDiscoveryConfigs(t *testing.T) {
	configs := BuildDiscoveryConfigs("home/example")
	if len(configs) != 2 {
		t.Fatalf("got %d configs", len(configs))
	}
	if configs[0].Component != "binary_sensor" {
		t.Fatalf("component %q", configs[0].Component)
	}
	stateTopic, _ := configs[0].Payload["state_topic"].(string)
	if stateTopic != "home/example/summary" {
		t.Fatalf("state_topic %q", stateTopic)
	}
}

func TestPublishSnapshot(t *testing.T) {
	sink := &fakeSink{}
	snap := BuildSnapshot(time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC))
	if err := PublishSnapshot(sink, snap); err != nil {
		t.Fatal(err)
	}
	if sink.suffix != "summary" {
		t.Fatalf("suffix %q", sink.suffix)
	}
	if !sink.retain {
		t.Fatal("expected retain")
	}
}

type fakeSink struct {
	suffix string
	retain bool
}

func (f *fakeSink) Publish(suffix string, _ any, retain bool) error {
	f.suffix = suffix
	f.retain = retain
	return nil
}

func (f *fakeSink) PublishQuiet(suffix string, payload any, retain bool) error {
	return f.Publish(suffix, payload, retain)
}

func (f *fakeSink) PublishRaw(string, any, bool) error { return nil }

func (f *fakeSink) PublishDiscovery([]mqttpub.Config, string) error { return nil }
