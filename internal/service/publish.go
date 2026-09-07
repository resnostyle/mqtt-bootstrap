package service

import (
	"log/slog"

	"github.com/resnostyle/mqttkit/hadisc"
	"github.com/resnostyle/mqttkit/mqttpub"
)

const (
	deviceManufacturer = "mqtt-bootstrap"
	deviceUID          = "mqtt_bootstrap"
)

func deviceBlock() map[string]any {
	return hadisc.Device([]string{deviceUID}, "MQTT Bootstrap", deviceManufacturer, "Poller")
}

// BuildDiscoveryConfigs returns HA MQTT discovery configs for the summary topic.
func BuildDiscoveryConfigs(topicPrefix string) []mqttpub.Config {
	summary := mqttpub.Join(topicPrefix, "summary")
	device := deviceBlock()

	ok := map[string]any{
		"name":                  "OK",
		"unique_id":             deviceUID + "_ok",
		"state_topic":           summary,
		"value_template":        "{{ value_json.ok }}",
		"payload_on":            "true",
		"payload_off":           "false",
		"device":                device,
		"object_id":             deviceUID + "_ok",
		"icon":                  "mdi:check-circle",
		"json_attributes_topic": summary,
	}

	updated := map[string]any{
		"name":                  "Last update",
		"unique_id":             deviceUID + "_last_update",
		"state_topic":           summary,
		"value_template":        "{{ value_json.updated_at }}",
		"device":                device,
		"object_id":             deviceUID + "_last_update",
		"device_class":          "timestamp",
		"icon":                  "mdi:clock-outline",
		"json_attributes_topic": summary,
	}

	return []mqttpub.Config{
		{ObjectID: deviceUID + "_ok", Component: "binary_sensor", Payload: ok},
		{ObjectID: deviceUID + "_last_update", Component: "sensor", Payload: updated},
	}
}

func PublishDiscovery(settings Settings, mqtt mqttpub.Sink) error {
	if !settings.MQTTDiscoveryEnabled {
		return nil
	}
	configs := BuildDiscoveryConfigs(settings.MQTTTopicPrefix)
	if err := mqtt.PublishDiscovery(configs, settings.MQTTDiscoveryPrefix); err != nil {
		return err
	}
	slog.Info("published mqtt discovery", "entities", len(configs))
	return nil
}

func PublishSnapshot(mqtt mqttpub.Sink, snap Snapshot) error {
	return mqtt.Publish("summary", snap, true)
}
