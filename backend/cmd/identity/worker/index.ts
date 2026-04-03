import { Container, getRandom } from "@cloudflare/containers";

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

const loadBalancedInstances = 10;

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const container = await getRandom(env.IDENTITY_CONTAINER, loadBalancedInstances);
    await container.startAndWaitForPorts({
      startOptions: {
        envVars: containerEnvFromWorkerEnv(env),
        enableInternet: true,
      },
    });
    return container.fetch(request);
  },
};
