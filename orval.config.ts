import { defineConfig } from "orval";

export default defineConfig({
  api: {
    output: {
      mode: "tags-split",
      target: "frontend/packages/api/src/generated/hooks",
      schemas: "frontend/packages/api/src/generated/models",
      client: "react-query",
      httpClient: "fetch",
      mock: false,
      prettier: false,
    },
    input: {
      target: "./backend/cmd/api/docs/openapi.bundled.yaml",
    },
  },
});
