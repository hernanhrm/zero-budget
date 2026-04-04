# Fly.io deployment (web app / SPA)

Static **Vite + TanStack Router** build served by **nginx** on port **8080**. Deploy from **`frontend/`** so the Docker build context matches the pnpm workspace ([`Dockerfile`](./Dockerfile)).

## First-time setup

```bash
cd frontend
fly auth login
fly apps create zero-budget-web   # or another name; update `app` in fly.toml if needed
```

## Build-time: `VITE_IDENTITY_URL`

Vite inlines this at **build** time. Pass it on every deploy that should talk to a real identity host:

```bash
fly deploy --build-arg VITE_IDENTITY_URL=https://zero-budget-identity.fly.dev
```

Replace with your actual Better Auth / identity service public URL (see [`backend/cmd/identity/FLY.md`](../backend/cmd/identity/FLY.md)).

Local image test:

```bash
docker build -f Dockerfile . \
  --build-arg VITE_IDENTITY_URL=http://localhost:8081 \
  -t zero-budget-web:local
docker run --rm -p 8080:8080 zero-budget-web:local
```

## Identity / CORS alignment

The browser loads the SPA from your Fly web hostname (e.g. `https://zero-budget-web.fly.dev`) and calls the identity API on `VITE_IDENTITY_URL`. On the **identity** app, set (e.g. `fly secrets set`):

- **`TRUSTED_ORIGINS`** — comma-separated **exact** browser origins for your SPA (e.g. `https://zero-budget-web.fly.dev`). Identity uses this for Better Auth **and** for Hono CORS unless you set `CORS_ORIGIN` separately.
- **`CORS_ORIGIN`** — optional; only if CORS must differ from `TRUSTED_ORIGINS`.
- **`APP_URL`** — same public frontend URL (invitation links, etc.).

## Optional: Cloudflare Pages

If you previously used Cloudflare Pages, remove or disconnect that project so it does not keep running failed builds. This app is intended to run on Fly only.

## Scale-to-zero

`fly.toml` matches the identity pattern: `auto_stop_machines`, `min_machines_running = 0`, `[deploy] ha = false`. Expect cold starts after idle.

## Health checks

Fly checks `GET /`. nginx serves `index.html` for `/`.
