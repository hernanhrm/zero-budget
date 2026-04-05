import { drizzle } from "drizzle-orm/node-postgres";
import pg from "pg";
import * as schema from "./schema.js";

/** libpq allows sslrootcert=system (use OS trust store). node-postgres parses it as a file path and calls readFileSync("system"). Strip it; Node TLS already uses system CAs when verifying. */
function postgresUrlForNodePg(connectionString: string): string {
  try {
    const u = new URL(connectionString);
    const rootCert = u.searchParams.get("sslrootcert");
    if (rootCert != null && rootCert.toLowerCase() === "system") {
      u.searchParams.delete("sslrootcert");
    }
    return u.toString();
  } catch {
    return connectionString;
  }
}

if (process.env.NODE_ENV === "production" && !process.env.DATABASE_URL?.trim()) {
  throw new Error("DATABASE_URL is required in production.");
}

const databaseUrl = process.env.DATABASE_URL?.trim()
  ? postgresUrlForNodePg(process.env.DATABASE_URL.trim())
  : undefined;

const pool = new pg.Pool({
  connectionString: databaseUrl,
});

export const db = drizzle({ client: pool, schema });
