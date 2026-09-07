# mqtt-bootstrap

Go poller template for Home Assistant MQTT discovery services.

Use this repo as a starting point for a new MQTT bridge. Shared infrastructure comes from [`mqttkit`](https://github.com/resnostyle/mqttkit):

- env loading (`mqttkit/env`)
- MQTT publish + HA discovery (`mqttkit/mqttpub`)
- interruptible poll loop (`mqttkit/poll`)
- slog setup (`mqttkit/logx`)
- payload helpers + HA device blocks (`payload`, `hadisc`)

Also includes Docker / Compose / CI / mise tasks.

No Home Assistant token required — only an MQTT broker (Mosquitto, EMQX, etc.). HA just needs MQTT discovery enabled.

## Layout

```
cmd/service/                 entrypoint
internal/service/            config, payload, discovery, poll loop
```

## Topics (defaults)

| Topic | Contents |
|-------|----------|
| `home/example/summary` | `{ ok, message, updated_at }` |

With discovery enabled, HA gets:

- `binary_sensor.mqtt_bootstrap_ok`
- `sensor.mqtt_bootstrap_last_update`

## Quick start

```bash
cp .env.example .env
go test ./...
mise run service    # loads .env
```

## Create a new service from this bootstrap

```bash
cp -a mqtt-bootstrap mything-mqtt
cd mything-mqtt
rm -rf .git   # if you init a fresh repo
```

Then rename identity:

| Placeholder | Replace with |
|-------------|--------------|
| folder / image `mqtt-bootstrap` | `mything-mqtt` |
| Go module path in `go.mod` | your module path |
| `MQTT_TOPIC_PREFIX=home/example` | `home/mything` |
| `MQTT_CLIENT_ID=mqtt-bootstrap` | `mything-mqtt` |
| discovery `deviceUID` / names in `publish.go` | your device id |
| mise tasks `service` | `mything` (optional) |

Implement domain work in:

1. `internal/lib/<api>/` — external client
2. `internal/service/payload.go` — JSON shape
3. `internal/service/publish.go` — topics + discovery entities
4. `Fetch()` in `run.go` — call the client each poll

Improve shared MQTT/env/poll behavior in **mqttkit**, not by copying packages into the new service.

## Configuration

See [`.env.example`](.env.example).

```bash
MQTT_HOST=127.0.0.1
MQTT_PORT=1883
MQTT_TOPIC_PREFIX=home/example
MQTT_CLIENT_ID=mqtt-bootstrap
MQTT_DISCOVERY_ENABLED=true
POLL_INTERVAL_SECONDS=300
LOG_LEVEL=INFO
```

## Docker

```bash
docker compose up -d --build
```

Update the image name in `docker-compose.yml` (and CI) when you fork this template.
