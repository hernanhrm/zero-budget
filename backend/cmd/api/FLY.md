# Fly.io deployment (Go API)

The Docker image is built from the **Go workspace root** [`backend/`](../../), not from this directory. The API [`Dockerfile`](Dockerfile) copies `go.work`, `cmd/`, `internal/`, and `pkg/`.

```bash
cd backend
fly auth login
fly apps create zero-budget-api   # if the name in ../../fly.toml is free; otherwise edit `app` there
fly postgres attach --app zero-budget-api <postgres-app-name>   # or set DATABASE_URL manually
fly secrets set \
  DATABASE_URL='postgresql://...' \
  IDENTITY_URL='https://zero-budget-identity.fly.dev' \
  INTERNAL_API_KEY='...' \
  RESEND_API_KEY='...' \
  RESEND_FROM_ADDRESS='...'
fly deploy
```

From the monorepo root:

```bash
fly deploy ./backend
```

See [Fly monorepo paths](https://fly.io/docs/reference/monorepo/).

## Secrets (`fly secrets set`)

| Name | Notes |
|------|--------|
| `DATABASE_URL` | **Required.** The API calls `Ping` at startup; deploy fails without a reachable Postgres. |
| `INTERNAL_API_KEY` | Must match the identity service’s `INTERNAL_API_KEY` (events + permission client). |
| `RESEND_API_KEY` / `RESEND_FROM_ADDRESS` | If you use Resend for email from the API stack. |

## Non-secret env (`fly.toml` `[env]` or `fly secrets set`)

| Name | Notes |
|------|--------|
| `IDENTITY_URL` | Base URL of the Better Auth / identity app (e.g. `https://<identity>.fly.dev`). Used by the permission client and must be reachable from this app. |
| `SERVICE_PORT` | Defaults to 8080 in code; set in `backend/fly.toml` for clarity. |

`DOCS_PATH` is set in the **Dockerfile** (`/app/docs`) so `/v1/docs` can load OpenAPI files baked into the image.

## Health checks vs Fly

- **`GET /ping`** — Always 200. Fly `[[http_service.checks]]` use this path so transient DB issues do not mark the Machine unhealthy.
- **`GET /health`** — Runs DB (and other) checkers; may return **503** when degraded. Use for operators/monitoring, not for Fly’s service check.

## Database migrations

This image only runs the API binary. Run migrations separately (e.g. `cmd/migrate` or `migrate` CLI against the same `DATABASE_URL`). See repo `AGENTS.md` / migrations docs.

## Scale to zero and at most one Machine

[`backend/fly.toml`](../../fly.toml) sets `min_machines_running = 0` (scale to zero when idle), autostop/autostart, and `[deploy] ha = false` so Fly does not keep a second redundant Machine in the pool.

If you still see two Machines after an older deploy, run:

```bash
fly scale count 1 -a zero-budget-api --yes
```

Expect cold starts and occasional health-check noise when Machines are stopped or booting.

## Adding another `cmd/*` service later

You usually **do not** change this app’s `fly.toml`. A new deployable gets its own Fly app and config (e.g. `fly.worker.toml`) with `[build] dockerfile = "cmd/<name>/Dockerfile"`, still deploying with context `backend/`. Update [`backend/.dockerignore`](../../.dockerignore) if the new tree adds huge artifacts (e.g. `node_modules`).
