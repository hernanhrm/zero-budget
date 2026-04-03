import { Container, getRandom } from "@cloudflare/containers";
import type { DurableObject } from "cloudflare:workers";

/**
 * Keys forwarded from Worker env into the Go API container (see localconfig.GetConfig).
 * Set each via: cd backend/cmd/api && npx wrangler secret put <KEY>
 */
const containerEnvKeys = [
  "DATABASE_URL",
  "SERVICE_PORT",
  "SERVICE_NAME",
  "DOCS_PATH",
  "RESEND_API_KEY",
  "RESEND_FROM_ADDRESS",
  "IDENTITY_URL",
  "INTERNAL_API_KEY",
  "CONFIG_ENV_PATH",
  "CONFIG_ENV_FILENAME",
] as const;

type ContainerEnvKey = (typeof containerEnvKeys)[number];

export type Env = {
  API_CONTAINER: DurableObjectNamespace<ApiContainer>;
} & Partial<Record<ContainerEnvKey, string>>;

/** Load-balanced pool size (must be <= containers.max_instances). */
const loadBalanceInstances = 4;

function containerEnvFromWorkerEnv(env: Env): Record<string, string> {
  const out: Record<string, string> = {};
  for (const key of containerEnvKeys) {
    const value = env[key];
    if (typeof value === "string" && value.length > 0) {
      out[key] = value;
    }
  }
  return out;
}

export class ApiContainer extends Container<Env> {
  defaultPort = 8080;
  sleepAfter = "30m";
  enableInternet = true;
  /** /ping always 200; /health can be 503 when DB checks fail. */
  pingEndpoint = "/ping";

  constructor(ctx: DurableObject<Env>["ctx"], env: Env) {
    super(ctx, env);
    this.envVars = {
      ...containerEnvFromWorkerEnv(env),
      SERVICE_PORT: env.SERVICE_PORT?.trim() || "8080",
    };
  }
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const container = await getRandom(env.API_CONTAINER, loadBalanceInstances);
    await container.startAndWaitForPorts({
      startOptions: {
        envVars: {
          ...containerEnvFromWorkerEnv(env),
          SERVICE_PORT: env.SERVICE_PORT?.trim() || "8080",
        },
        enableInternet: true,
      },
    });
    return container.fetch(request);
  },
};
