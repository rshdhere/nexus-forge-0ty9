# Nexus Chat

Realtime chat app in Go — shared room, SSE live updates, in-memory message history.

## Getting started

```bash
go run .
# or
bun run start
```

Open `http://localhost:3000`. Health check: `GET /health` → `{ "ok": true }`.

## API

- `GET /api/messages` — message history
- `POST /api/messages` — `{ "user", "text" }` send a message
- `GET /api/stream` — Server-Sent Events for live messages + presence
- `POST /api/presence` — heartbeat with `{ "user" }`
