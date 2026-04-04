# Fly.io deployment (identity service)

Deploy from this directory so the Docker build context matches the `Dockerfile` (`COPY` paths).

```bash
cd backend/cmd/identity
fly auth login   # once per machine
fly apps create zero-budget-identity   # if the name in fly.toml is not taken; otherwise pick a name and update `app` in fly.toml
fly postgres create --name zero-budget-db   # optional; or use an external Postgres
fly postgres attach --app zero-budget-identity zero-budget-db   # sets DATABASE_URL on the app
fly secrets set \
  BETTER_AUTH_SECRET='<openssl rand -base64 32>' \
  BETTER_AUTH_URL='https://zero-budget-identity.fly.dev' \
  DATABASE_URL='postgres://...' \
  INTERNAL_API_KEY='...'
fly deploy
```

If you already use `fly launch`, it can create the app and wire `fly.toml`; adjust `app` and region as needed.

## Secrets (`fly secrets set`)

Store sensitive values here (not in `fly.toml` `[env]`).

| Name | Notes |
|------|--------|
| `BETTER_AUTH_SECRET` | Required in production. |
| `BETTER_AUTH_URL` | Public HTTPS base URL of this service (e.g. `https://<app>.fly.dev`). Drives cookie `secure` and OAuth callbacks. |
| `DATABASE_URL` | Postgres connection string. Often set automatically when attaching Fly Postgres. |
| `INTERNAL_API_KEY` | Shared secret for calls to the Go API (`GO_API_URL`), if used. |
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` | Optional OAuth. |
| `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET` | Optional OAuth. |

## Non-secret env (`fly deploy` / dashboard / `[env]`)

Set with `fly secrets set` only if you prefer; these are not inherently secret:

| Name | Notes |
|------|--------|
| `GO_API_URL` | See `.env.example` — use internal URL on Fly when calling sibling apps: `http://<go-app>.internal:<port>`. |
| `TRUSTED_ORIGINS` | Comma-separated **browser** origins (scheme + host + port, no path). Used by Better Auth **and** by Hono CORS when `CORS_ORIGIN` is unset. List the SPA URL(s) users load in the browser (e.g. `https://zero-budget-web.fly.dev`). You do **not** need the Go API URL here for normal SPA requests—the browser’s `Origin` header is always the page’s site, not the API. |
| `CORS_ORIGIN` | Optional. Comma-separated allowed origins for Hono CORS only. If unset, `TRUSTED_ORIGINS` is used so CORS matches Better Auth. Set this only if CORS must differ from `TRUSTED_ORIGINS`. |
| `APP_URL` | Public frontend URL for invitation links. |

Example:

```bash
fly secrets set \
  TRUSTED_ORIGINS='https://zero-budget-web.fly.dev' \
  APP_URL='https://zero-budget-web.fly.dev' \
  GO_API_URL='http://zero-budget-api.internal:8080'
```

(`CORS_ORIGIN` omitted—Hono reuses `TRUSTED_ORIGINS`.)

(Adjust app hostnames and ports to match your Fly apps.)

## OAuth redirect URLs

In Google/GitHub developer consoles, add callback URLs that match Better Auth’s expectations for your `BETTER_AUTH_URL`.

## Scale to zero and at most one Machine (cost)

`fly.toml` sets `auto_stop_machines = "stop"`, `auto_start_machines = true`, `min_machines_running = 0`, and `[deploy] ha = false` so Fly does not add a second redundant Machine on deploy. After idle periods, Fly Proxy can stop all Machines in the region; you are not billed for CPU/RAM while they are stopped ([stopped Machines pricing](https://fly.io/docs/about/pricing/#stopped-fly-machines)). The next request triggers a **cold start** (several seconds of latency) while a Machine boots.

Keep the generous `[[http_service.checks]]` `grace_period` / `timeout` so new Machines can pass checks after autostart. Expect occasional health-check log noise while Machines are stopped or booting.

If an older deploy left two Machines in the pool, run:

```bash
fly scale count 1 -a zero-budget-identity --yes
```

Fly Postgres and other resources bill separately from this app’s Machines.

## Health checks

Fly is configured to `GET /health` (see `fly.toml`). The app listens on `PORT` (8080 on Fly) and `HOST=0.0.0.0` in the image.

### “Health check on port 8080 has failed” while the app still loads in the browser

Fly runs **per-machine** service checks. The proxy only sends traffic to machines that pass; if you have **more than one machine** and one is unhealthy, public requests can still succeed while logs show failures for the bad instance ([health checks](https://fly.io/docs/reference/health-checks/)).

Typical causes:

1. **Cold start vs grace period** — With `auto_stop_machines`, a machine that was stopped needs time to boot Node before `/health` returns 200. `fly.toml` uses a longer `grace_period` / `timeout` for that.
2. **A stuck or bad machine** — Inspect and reconcile:

   ```bash
   fly checks list -a zero-budget-identity
   fly machine list -a zero-budget-identity
   fly logs -a zero-budget-identity
   ```

   If one machine never becomes healthy, destroy it and run a single instance until stable:

   ```bash
   fly machine destroy <machine-id> -a zero-budget-identity --force
   fly scale count 1 -a zero-budget-identity
   fly deploy -a zero-budget-identity
   ```

3. **Process crash on boot** — Missing/invalid secrets (e.g. `DATABASE_URL`, `BETTER_AUTH_*`) can prevent the process from listening; the healthy sibling would still answer. Confirm `fly secrets list` and logs for startup errors.

## Replace placeholder `DATABASE_URL`

If the app was deployed with a temporary `DATABASE_URL` only to satisfy startup checks, point it at a real Postgres cluster before using auth or any DB-backed routes:

```bash
fly postgres attach --app zero-budget-identity <your-postgres-app-name>
# or
fly secrets set DATABASE_URL='postgresql://...'
fly deploy
```
