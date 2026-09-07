# mqtt-bootstrap

Generic Go poller template for Home Assistant MQTT services.

Copy this folder when starting a new bridge (grades, bus, whatever). Shared infrastructure comes from [`mqttkit`](../mqttkit):

- env loading (`github.com/resnostyle/mqttkit/env`)
- MQTT publish + HA discovery (`github.com/resnostyle/mqttkit/mqttpub`)
- interruptible poll loop (`github.com/resnostyle/mqttkit/poll`)
- slog setup (`github.com/resnostyle/mqttkit/logx`)
- payload helpers + HA device blocks (`payload`, `hadisc`)

Also includes Docker / Compose / CI / mise tasks.

No Home Assistant token required — only Mosquitto (or EMQX). HA just needs MQTT discovery enabled.

## Layout

```
cmd/service/                 entrypoint
internal/service/            config, payload, discovery, poll loop
```

Shared libs live in the sibling `mqttkit` module (resolved via root `go.work` or the `replace` in `go.mod`).

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
# from mqtt-services/
cp -a mqtt-bootstrap mything-mqtt
cd mything-mqtt
rm -rf .git   # if you init a fresh repo
```

Then rename identity:

| Placeholder | Replace with |
|-------------|--------------|
| folder `mqtt-bootstrap` | `mything-mqtt` |
| module `github.com/resnostyle/mqtt-bootstrap` | `github.com/resnostyle/mything-mqtt` |
| `cmd/service` / package paths | keep or rename to `cmd/mything` |
| `MQTT_TOPIC_PREFIX=home/example` | `home/mything` |
| `MQTT_CLIENT_ID=mqtt-bootstrap` | `mything-mqtt` |
| discovery `deviceUID` / names in `publish.go` | your device id |
| GHCR image `mqtt-bootstrap` | `mything-mqtt` |
| mise tasks `service` | `mything` |

Add the new module to the workspace root `go.work`, and keep:

```
require github.com/resnostyle/mqttkit v0.0.0
replace github.com/resnostyle/mqttkit => ../mqttkit
```

until `mqttkit` is published as its own module.

Implement domain work in:

1. `internal/lib/<api>/` — external client
2. `internal/service/payload.go` — JSON shape
3. `internal/service/publish.go` — topics + discovery entities
4. `Fetch()` in `run.go` — call the client each poll

Improve shared MQTT/env/poll behavior in **`mqttkit`**, not by copying packages into the new service.

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

Image: `ghcr.io/resnostyle/mqtt-bootstrap:latest`

## Kubernetes

Clone an existing chart under `k8s-gitops/.../apps/automation/` (e.g. `weather-mqtt` or `pinger-mqtt`), point the image at your new GHCR repo, and put secrets in Vault.
