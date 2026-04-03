import { Container, getContainer } from "@cloudflare/containers";

// Each key must exist as a Worker secret (or plain var) so Wrangler injects `env.*`
// and it is forwarded into the container. Example:
//   cd backend/cmd/identity && npx wrangler secret put BETTER_AUTH_SECRET
//   npx wrangler secret put BETTER_AUTH_URL
// BETTER_AUTH_URL must be the public Worker URL (e.g. https://zero-budget-identity.<subdomain>.workers.dev).
const containerEnvKeys = [
  "DATABASE_URL",
  "TRUSTED_ORIGINS",
  "CORS_ORIGIN",
  "APP_URL",
  "GO_API_URL",
  "INTERNAL_API_KEY",
  "GOOGLE_CLIENT_ID",
  "GOOGLE_CLIENT_SECRET",
  "GITHUB_CLIENT_ID",
  "GITHUB_CLIENT_SECRET",
  "BETTER_AUTH_SECRET",
  "BETTER_AUTH_URL",
] as const;

type ContainerEnvKey = (typeof containerEnvKeys)[number];

export type Env = {
  IDENTITY_CONTAINER: DurableObjectNamespace<IdentityContainer>;
} & Partial<Record<ContainerEnvKey, string>>;

function containerEnvFromWorkerEnv(env: Env): Record<string, string> {
  const out: Record<string, string> = {
    NODE_ENV: "production",
    PORT: "8081",
  };
  for (const key of containerEnvKeys) {
    const value = env[key];
    if (typeof value === "string" && value.length > 0) {
      out[key] = value;
    }
  }
  return out;
}

export class IdentityContainer extends Container<Env> {
  defaultPort = 8081;
  sleepAfter = "30m";
  enableInternet = true;
  pingEndpoint = "/health";

  constructor(ctx: DurableObjectState, env: Env) {
    super(ctx, env);
    this.envVars = containerEnvFromWorkerEnv(env);
  }
}

// One Durable Object / one container. getRandom(N) cold-starts a different instance per request and looks like a restart loop.
export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const container = getContainer(env.IDENTITY_CONTAINER);
    await container.startAndWaitForPorts({
      startOptions: {
        envVars: containerEnvFromWorkerEnv(env),
        enableInternet: true,
      },
    });
    return container.fetch(request);
  },
};
