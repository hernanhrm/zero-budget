import { drizzle } from "drizzle-orm/node-postgres";
import pg from "pg";
import * as schema from "./schema.js";

if (process.env.NODE_ENV === "production" && !process.env.DATABASE_URL?.trim()) {
  throw new Error(
    "DATABASE_URL is required. Set a Worker secret: npx wrangler secret put DATABASE_URL (from backend/cmd/identity)",
  );
}

const pool = new pg.Pool({
  connectionString: process.env.DATABASE_URL,
});

export const db = drizzle({ client: pool, schema });
