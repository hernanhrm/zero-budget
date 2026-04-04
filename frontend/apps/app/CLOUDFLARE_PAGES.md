# Cloudflare Pages (web app)

This package is a **Vite + TanStack Router** SPA. Deploy it as a static site on [Cloudflare Pages](https://developers.cloudflare.com/pages/).

## Dashboard (Git integration)

In Cloudflare **Workers & Pages** → **Create** → **Pages** → connect your Git repository, then set:

| Setting | Value |
|--------|--------|
| **Root directory** | `frontend` |
| **Build command** | `pnpm install --frozen-lockfile && pnpm --filter app build` |
| **Build output directory** | `apps/app/dist` |
| **Deploy command** | **Leave empty.** After a successful build, Pages uploads the build output directory automatically ([build configuration](https://developers.cloudflare.com/pages/configuration/build-configuration/)). |

**Do not** set the deploy command to `npx wrangler versions upload` unless you mean to publish via **Workers** [static assets](https://developers.cloudflare.com/workers/static-assets/). That command does **not** read `pages_build_output_dir`; it needs a Worker `main` script or an `assets.directory` in [`wrangler.jsonc`](../../wrangler.jsonc). This repo sets both `pages_build_output_dir` and `assets.directory` to `./apps/app/dist` so `versions upload` can resolve the folder—but **clearing the deploy command** is still the right setup for a normal **Pages** Git project (Pages uploads the build output for you).

If you truly need Wrangler in the deploy step (unusual for Git-connected Pages), use **Pages** deploy instead, from the repo root `frontend/`:

```bash
npx wrangler pages deploy apps/app/dist --project-name=<your-pages-project-name>
```

Or rely on [`../../wrangler.jsonc`](../../wrangler.jsonc) (`pages_build_output_dir` + `name`) and run `npx wrangler pages deploy` with no positional path—after setting `"name"` to match your dashboard project.

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
