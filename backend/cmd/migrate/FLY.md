# Fly.io one-shot database migrations

The [`migrate`](./main.go) binary embeds SQL from [`migrations/`](./migrations/) and exits after **`up`**, **`down`**, or **`version`**. It is **not** a long-running HTTP service.

Use **`fly machine run`** from the **repository `backend/` directory** so Docker can `COPY cmd/migrate/` (see [`Dockerfile`](./Dockerfile)). The image builds **`cmd/migrate`** as a **standalone module** (no `go.work` in the image), because a full workspace copy would require every `go.work` entry to be present for `go mod download`.

Related apps: [Go API](../api/FLY.md), [identity](../identity/FLY.md) — use the **same Postgres** (`DATABASE_URL`) as the API and identity services.

## One-time setup

```bash
cd backend
fly auth login
fly apps create zero-budget-migrate   # or change `app` in [fly.toml](./fly.toml)
```

## Networking

The migration Machine must **reach Postgres**:

- **Public URL** — same `DATABASE_URL` you use locally or from CI (often includes `?sslmode=require`).
- **Fly private network** — if the database is only reachable on Fly’s private network, use a hostname your org can resolve from a Machine (e.g. Fly Postgres attachment strings or internal hostnames). The migrate Machine runs in your org like other apps.

## `DATABASE_URL` as a Fly secret (recommended)

Set the connection string on the **migrate** app once; Fly injects app [secrets](https://fly.io/docs/apps/secrets/) as **environment variables** when each Machine starts (including Machines created with **`fly machine run`**):

```bash
fly secrets set DATABASE_URL='postgresql://...' -a zero-budget-migrate
```

Use the **same** URL (or an equivalent reachable URL) as your API/identity apps. After changing a staged secret (`--stage`), new Machines pick it up; existing ones need a restart or `fly secrets deploy` per Fly docs.

You do **not** need **`--env DATABASE_URL=...`** on `fly machine run` if this secret is set—`/migrate` reads `DATABASE_URL` from the environment. You may still pass **`--env DATABASE_URL=...`** to override for a single run (e.g. pointing at a staging DB).

## Run migrations (production)

From **`backend/`** (with **`DATABASE_URL`** secret set on `zero-budget-migrate`, or pass **`--env`** as above):

```bash
fly machine run . --dockerfile cmd/migrate/Dockerfile --app zero-budget-migrate \
  --region iad --rm
```

- **`--rm`** — destroy the Machine when `/migrate` exits (one-shot; restart policy `no`).
- **`--region`** — align with your Postgres / primary region (e.g. `iad`).

Override the command after the image argument (replaces default `up`):

```bash
# Current schema version
fly machine run . --dockerfile cmd/migrate/Dockerfile --app zero-budget-migrate \
  --region iad --rm -- version

# Roll back one migration (use with care)
fly machine run . --dockerfile cmd/migrate/Dockerfile --app zero-budget-migrate \
  --region iad --rm -- down 1
```

`fly machine run` **disregards** most of `fly.toml` except the default **`app`** name when you run from this directory; you can still pass **`--app`** explicitly as above.

## Local / CI

From **`backend/`**:

```bash
docker build -f cmd/migrate/Dockerfile . -t zb-migrate:local
docker run --rm -e DATABASE_URL='postgresql://...' zb-migrate:local up
```

## Commands (binary)

| Args | Effect |
|------|--------|
| *(none)* / `up` | Apply all pending migrations |
| `down` | Roll back one step |
| `down N` | Roll back N steps |
| `version` | Print current version and dirty flag |
