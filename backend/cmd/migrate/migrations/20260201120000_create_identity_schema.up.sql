CREATE SCHEMA IF NOT EXISTS identity;

-- Users
CREATE TABLE "identity"."users" (
	"id" text PRIMARY KEY NOT NULL,
	"name" text NOT NULL,
	"email" text NOT NULL,
	"email_verified" boolean DEFAULT false NOT NULL,
	"image" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	"two_factor_enabled" boolean DEFAULT false,
	"role" text,
	"banned" boolean DEFAULT false,
	"ban_reason" text,
	"ban_expires" timestamp,
	CONSTRAINT "users_email_unique" UNIQUE("email")
);

-- Accounts
CREATE TABLE "identity"."accounts" (
	"id" text PRIMARY KEY NOT NULL,
	"account_id" text NOT NULL,
	"provider_id" text NOT NULL,
	"user_id" text NOT NULL REFERENCES "identity"."users"("id") ON DELETE CASCADE,
	"access_token" text,
	"refresh_token" text,
	"id_token" text,
	"access_token_expires_at" timestamp,
	"refresh_token_expires_at" timestamp,
	"scope" text,
	"password" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp NOT NULL
);
CREATE INDEX "accounts_user_id_idx" ON "identity"."accounts" USING btree ("user_id");

-- Sessions
CREATE TABLE "identity"."sessions" (
	"id" text PRIMARY KEY NOT NULL,
	"expires_at" timestamp NOT NULL,
	"token" text NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp NOT NULL,
	"ip_address" text,
	"user_agent" text,
	"user_id" text NOT NULL REFERENCES "identity"."users"("id") ON DELETE CASCADE,
	"active_organization_id" text,
	"impersonated_by" text,
	CONSTRAINT "sessions_token_unique" UNIQUE("token")
);
CREATE INDEX "sessions_user_id_idx" ON "identity"."sessions" USING btree ("user_id");

-- Verifications
CREATE TABLE "identity"."verifications" (
	"id" text PRIMARY KEY NOT NULL,
	"identifier" text NOT NULL,
	"value" text NOT NULL,
	"expires_at" timestamp NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
CREATE INDEX "verifications_identifier_idx" ON "identity"."verifications" USING btree ("identifier");

-- Organizations
CREATE TABLE "identity"."organizations" (
	"id" text PRIMARY KEY NOT NULL,
	"name" text NOT NULL,
	"slug" text NOT NULL,
	"logo" text,
	"created_at" timestamp NOT NULL,
	"metadata" text,
	CONSTRAINT "organizations_slug_unique" UNIQUE("slug")
);
CREATE UNIQUE INDEX "organizations_slug_uidx" ON "identity"."organizations" USING btree ("slug");

-- Members
CREATE TABLE "identity"."members" (
	"id" text PRIMARY KEY NOT NULL,
	"organization_id" text NOT NULL REFERENCES "identity"."organizations"("id") ON DELETE CASCADE,
	"user_id" text NOT NULL REFERENCES "identity"."users"("id") ON DELETE CASCADE,
	"role" text DEFAULT 'member' NOT NULL,
	"created_at" timestamp NOT NULL
);
CREATE INDEX "members_organization_id_idx" ON "identity"."members" USING btree ("organization_id");
CREATE INDEX "members_user_id_idx" ON "identity"."members" USING btree ("user_id");

-- Invitations
CREATE TABLE "identity"."invitations" (
	"id" text PRIMARY KEY NOT NULL,
	"organization_id" text NOT NULL REFERENCES "identity"."organizations"("id") ON DELETE CASCADE,
	"email" text NOT NULL,
	"role" text,
	"status" text DEFAULT 'pending' NOT NULL,
	"expires_at" timestamp NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"inviter_id" text NOT NULL REFERENCES "identity"."users"("id") ON DELETE CASCADE
);
CREATE INDEX "invitations_organization_id_idx" ON "identity"."invitations" USING btree ("organization_id");
CREATE INDEX "invitations_email_idx" ON "identity"."invitations" USING btree ("email");

-- Organization Roles (dynamic access control)
CREATE TABLE "identity"."organization_roles" (
	"id" text PRIMARY KEY NOT NULL,
	"organization_id" text NOT NULL REFERENCES "identity"."organizations"("id") ON DELETE CASCADE,
	"role" text NOT NULL,
	"permission" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
CREATE INDEX "organization_roles_organization_id_idx" ON "identity"."organization_roles" USING btree ("organization_id");

-- Two Factors
CREATE TABLE "identity"."two_factors" (
	"id" text PRIMARY KEY NOT NULL,
	"secret" text NOT NULL,
	"backup_codes" text NOT NULL,
	"user_id" text NOT NULL REFERENCES "identity"."users"("id") ON DELETE CASCADE
);
CREATE INDEX "twoFactors_secret_idx" ON "identity"."two_factors" USING btree ("secret");
CREATE INDEX "twoFactors_user_id_idx" ON "identity"."two_factors" USING btree ("user_id");

-- API Keys
CREATE TABLE "identity"."api_keys" (
	"id" text PRIMARY KEY NOT NULL,
	"config_id" text DEFAULT 'default' NOT NULL,
	"name" text,
	"start" text,
	"reference_id" text NOT NULL,
	"prefix" text,
	"key" text NOT NULL,
	"refill_interval" integer,
	"refill_amount" integer,
	"last_refill_at" timestamp,
	"enabled" boolean DEFAULT true,
	"rate_limit_enabled" boolean DEFAULT true,
	"rate_limit_time_window" integer DEFAULT 86400000,
	"rate_limit_max" integer DEFAULT 10,
	"request_count" integer DEFAULT 0,
	"remaining" integer,
	"last_request" timestamp,
	"expires_at" timestamp,
	"created_at" timestamp NOT NULL,
	"updated_at" timestamp NOT NULL,
	"permissions" text,
	"metadata" text
);
CREATE INDEX "api_keys_config_id_idx" ON "identity"."api_keys" USING btree ("config_id");
CREATE INDEX "api_keys_reference_id_idx" ON "identity"."api_keys" USING btree ("reference_id");
CREATE INDEX "api_keys_key_idx" ON "identity"."api_keys" USING btree ("key");
