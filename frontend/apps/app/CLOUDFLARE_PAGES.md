# Cloudflare Pages (web app)

This package is a **Vite + TanStack Router** SPA. Deploy it as a static site on [Cloudflare Pages](https://developers.cloudflare.com/pages/).

## Dashboard (Git integration)

In Cloudflare **Workers & Pages** → **Create** → **Pages** → connect your Git repository, then set:

| Setting | Value |
|--------|--------|
| **Root directory** | `frontend` |
| **Build command** | `pnpm install --frozen-lockfile && pnpm --filter app build` (**required** — if this is empty, `dist` is never created) |
| **Build output directory** | `apps/app/dist` |
| **Deploy command** | **Leave empty** (recommended). Pages uploads **Build output directory** after a successful build ([build configuration](https://developers.cloudflare.com/pages/configuration/build-configuration/)). |

**Wrangler and this monorepo**

Two config files are intentional:

| File | Role |
|------|------|
| [`../../wrangler.jsonc`](../../wrangler.jsonc) (`frontend/`) | `assets.directory` → `./apps/app/dist` so **`npx wrangler versions upload`** works when the dashboard cwd is `frontend/` (Workers-style upload). |
| [`./wrangler.jsonc`](./wrangler.jsonc) (`apps/app/`) | `pages_build_output_dir` → `./dist` for **`wrangler pages deploy --cwd apps/app`** and dashboard alignment. **No** `assets` key here (Pages validation rejects it in that context). |

- **Do not** run `npx wrangler deploy` for a **Pages** project. Wrangler blocks it and tells you to use `wrangler pages deploy` instead.
- **`npx wrangler deploy`** from `frontend/` also hits monorepo workspace detection unless you use **`--cwd apps/app`** (still wrong for Pages—use `pages deploy`).

**Build command must run before any Wrangler step.** If logs show only `pnpm install` and then deploy, `apps/app/dist` does not exist and uploads will fail or ship an empty tree.

If you insist on a **custom deploy command** (unusual for Git-connected Pages), prefer **Pages** deploy from `frontend/`:

```bash
npx wrangler pages deploy --cwd apps/app
```

Add `--project-name=<your-pages-project-name>` if needed. Or:

```bash
npx wrangler pages deploy apps/app/dist --project-name=<your-pages-project-name>
```

If the dashboard is set to **`npx wrangler versions upload`** (preview/production), the repo root [`../../wrangler.jsonc`](../../wrangler.jsonc) supplies the assets directory. Alternative without that file: `npx wrangler versions upload --cwd apps/app --assets=./dist`.

### Environment variables

Add under **Settings → Environment variables** (Production and Preview as needed):

| Name | Purpose |
|------|--------|
| `VITE_IDENTITY_URL` | Public base URL of the Better Auth / identity service (e.g. `https://zero-budget-identity.fly.dev`). Inlined at **build** time by Vite. |
| `NODE_VERSION` | Optional but recommended: `22` (or `20`) so the Pages build uses a current Node ([build image](https://developers.cloudflare.com/pages/configuration/build-image/)). |

Changing `VITE_IDENTITY_URL` requires a **new deployment** (rebuild), not a runtime toggle.

### SPA routing

[`public/_redirects`](./public/_redirects) is copied into `dist` so client-side routes (e.g. `/sign-in`) resolve on refresh:

```txt
/*    /index.html   200
```

## Identity / CORS alignment

The browser loads the SPA from your Pages hostname and calls the identity API on `VITE_IDENTITY_URL`. On the **identity** service (see [`backend/cmd/identity/FLY.md`](../../../backend/cmd/identity/FLY.md)) set:

- **`TRUSTED_ORIGINS`** — include your Pages origin(s), e.g. `https://<project>.pages.dev` and any custom domain.
- **`CORS_ORIGIN`** — typically the same primary frontend origin used by the Hono CORS config.
- **`APP_URL`** — same public frontend base URL (invitation links, etc.).

## Optional: Wrangler CLI deploy

If you prefer uploading a folder without Git:

1. Build locally or in CI from the monorepo:

   ```bash
   cd frontend
   pnpm install --frozen-lockfile
   pnpm --filter app build
   ```

2. Deploy the output directory:

   ```bash
   cd frontend
   npx wrangler pages deploy apps/app/dist --project-name=<your-pages-project-name>
   ```

   Use the same `VITE_IDENTITY_URL` when building (export in the shell or use a `.env` file that Vite reads during `pnpm --filter app build`).

3. **GitHub Action** (optional): use [cloudflare/pages-action](https://github.com/cloudflare/pages-action) with `workingDirectory: frontend`, build command and output path as above, and pass secrets for `CLOUDFLARE_API_TOKEN` / `CLOUDFLARE_ACCOUNT_ID`.

## Verification

- After build: `apps/app/dist` should contain `index.html` and `_redirects`.
- After deploy: open a deep link (e.g. `/sign-in`) in a new tab; you should not get a 404 from Pages.
- Test sign-in against the configured identity URL once CORS/trusted origins include your Pages URL.
