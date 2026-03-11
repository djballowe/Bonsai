# Bonsai

A minimal web dashboard for monitoring a Bambu Lab A1 Mini 3D printer in real time. Connects to the printer over MQTT (SSL/TLS), parses status messages, and streams live updates to the browser.

## Environment Variables

| Variable    | Description           |
| ----------- | --------------------- |
| `IP`        | Printer IP address    |
| `SERIAL`    | Printer serial number |
| `MQTT_USER` | MQTT username         |
| `PASS`      | MQTT password         |

Copy `.env.example` to `.env` and fill in your printer's values.

## Running

```sh
go run .
```

The server starts at `http://localhost:3100`. The MQTT client connects to `ssl://<IP>:8883` and subscribes to `device/<serial>/report`.
