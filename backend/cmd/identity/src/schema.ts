import { relations } from "drizzle-orm";
import {
  pgSchema,
  text,
  timestamp,
  boolean,
  index,
  uniqueIndex,
} from "drizzle-orm/pg-core";

const identitySchema = pgSchema("identity");
const pgTable = identitySchema.table;

export const users = pgTable("users", {
  id: text("id").primaryKey(),
  name: text("name").notNull(),
  email: text("email").notNull().unique(),
  email_verified: boolean("email_verified").default(false).notNull(),
  image: text("image"),
  created_at: timestamp("created_at").defaultNow().notNull(),
  updated_at: timestamp("updated_at")
    .defaultNow()
    .$onUpdate(() => /* @__PURE__ */ new Date())
    .notNull(),
  two_factor_enabled: boolean("two_factor_enabled").default(false),
  role: text("role"),
  banned: boolean("banned").default(false),
  ban_reason: text("ban_reason"),
  ban_expires: timestamp("ban_expires"),
});

export const sessions = pgTable(
  "sessions",
  {
    id: text("id").primaryKey(),
    expires_at: timestamp("expires_at").notNull(),
    token: text("token").notNull().unique(),
    created_at: timestamp("created_at").defaultNow().notNull(),
    updated_at: timestamp("updated_at")
      .$onUpdate(() => /* @__PURE__ */ new Date())
      .notNull(),
    ip_address: text("ip_address"),
    user_agent: text("user_agent"),
    user_id: text("user_id")
      .notNull()
      .references(() => users.id, { onDelete: "cascade" }),
    activeOrganizationId: text("active_organization_id"),
    impersonated_by: text("impersonated_by"),
  },
  (table) => [index("sessions_user_id_idx").on(table.user_id)],
);

export const accounts = pgTable(
  "accounts",
  {
    id: text("id").primaryKey(),
    account_id: text("account_id").notNull(),
    provider_id: text("provider_id").notNull(),
    user_id: text("user_id")
      .notNull()
      .references(() => users.id, { onDelete: "cascade" }),
    access_token: text("access_token"),
    refresh_token: text("refresh_token"),
    id_token: text("id_token"),
    access_token_expires_at: timestamp("access_token_expires_at"),
    refresh_token_expires_at: timestamp("refresh_token_expires_at"),
    scope: text("scope"),
    password: text("password"),
    created_at: timestamp("created_at").defaultNow().notNull(),
    updated_at: timestamp("updated_at")
      .$onUpdate(() => /* @__PURE__ */ new Date())
      .notNull(),
  },
  (table) => [index("accounts_user_id_idx").on(table.user_id)],
);

export const verifications = pgTable(
  "verifications",
  {
    id: text("id").primaryKey(),
    identifier: text("identifier").notNull(),
    value: text("value").notNull(),
    expires_at: timestamp("expires_at").notNull(),
    created_at: timestamp("created_at").defaultNow().notNull(),
    updated_at: timestamp("updated_at")
      .defaultNow()
      .$onUpdate(() => /* @__PURE__ */ new Date())
      .notNull(),
  },
  (table) => [index("verifications_identifier_idx").on(table.identifier)],
);

export const organizations = pgTable(
  "organizations",
  {
    id: text("id").primaryKey(),
    name: text("name").notNull(),
    slug: text("slug").notNull().unique(),
    logo: text("logo"),
    created_at: timestamp("created_at").notNull(),
    metadata: text("metadata"),
  },
  (table) => [uniqueIndex("organizations_slug_uidx").on(table.slug)],
);

export const members = pgTable(
  "members",
  {
    id: text("id").primaryKey(),
    organization_id: text("organization_id")
      .notNull()
      .references(() => organizations.id, { onDelete: "cascade" }),
    user_id: text("user_id")
      .notNull()
      .references(() => users.id, { onDelete: "cascade" }),
    role: text("role").default("member").notNull(),
    created_at: timestamp("created_at").notNull(),
  },
  (table) => [
    index("members_organization_id_idx").on(table.organization_id),
    index("members_user_id_idx").on(table.user_id),
  ],
);

export const invitations = pgTable(
  "invitations",
  {
    id: text("id").primaryKey(),
    organization_id: text("organization_id")
      .notNull()
      .references(() => organizations.id, { onDelete: "cascade" }),
    email: text("email").notNull(),
    role: text("role"),
    status: text("status").default("pending").notNull(),
    expires_at: timestamp("expires_at").notNull(),
    createdAt: timestamp("created_at").defaultNow().notNull(),
    inviter_id: text("inviter_id")
      .notNull()
      .references(() => users.id, { onDelete: "cascade" }),
  },
  (table) => [
    index("invitations_organization_id_idx").on(table.organization_id),
    index("invitations_email_idx").on(table.email),
  ],
);

export const twoFactors = pgTable(
  "two_factors",
  {
    id: text("id").primaryKey(),
    secret: text("secret").notNull(),
    backupCodes: text("backup_codes").notNull(),
    user_id: text("user_id")
      .notNull()
      .references(() => users.id, { onDelete: "cascade" }),
  },
  (table) => [
    index("twoFactors_secret_idx").on(table.secret),
    index("twoFactors_user_id_idx").on(table.user_id),
  ],
);

export const usersRelations = relations(users, ({ many }) => ({
  sessions: many(sessions),
  accounts: many(accounts),
  members: many(members),
  invitations: many(invitations),
  twoFactors: many(twoFactors),
}));

export const sessionsRelations = relations(sessions, ({ one }) => ({
  users: one(users, {
    fields: [sessions.user_id],
    references: [users.id],
  }),
}));

export const accountsRelations = relations(accounts, ({ one }) => ({
  users: one(users, {
    fields: [accounts.user_id],
    references: [users.id],
  }),
}));

export const organizationsRelations = relations(organizations, ({ many }) => ({
  members: many(members),
  invitations: many(invitations),
}));

export const membersRelations = relations(members, ({ one }) => ({
  organizations: one(organizations, {
    fields: [members.organization_id],
    references: [organizations.id],
  }),
  users: one(users, {
    fields: [members.user_id],
    references: [users.id],
  }),
}));

export const invitationsRelations = relations(invitations, ({ one }) => ({
  organizations: one(organizations, {
    fields: [invitations.organization_id],
    references: [organizations.id],
  }),
  users: one(users, {
    fields: [invitations.inviter_id],
    references: [users.id],
  }),
}));

export const twoFactorsRelations = relations(twoFactors, ({ one }) => ({
  users: one(users, {
    fields: [twoFactors.user_id],
    references: [users.id],
  }),
}));
