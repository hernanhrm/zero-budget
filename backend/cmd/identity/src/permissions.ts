import { createAccessControl } from "better-auth/plugins/access";

const statement = {
  organization: ["update", "delete"],
  member: ["create", "read", "update", "delete"],
  invitation: ["create", "read", "update", "delete", "cancel"],
  team: ["create", "read", "update", "delete"],
  emailTemplate: ["create", "read", "update", "delete"],
  emailLog: ["read"],
  currency: ["read"],
  organizationCurrency: ["create", "read", "update", "delete"],
  account: ["create", "read", "update", "delete"],
  category: ["create", "read", "update", "delete"],
  budget: ["create", "read", "update", "delete"],
  transaction: ["create", "read", "update", "delete"],
} as const;

export const ac = createAccessControl(statement);

export const owner = ac.newRole({
  organization: ["update", "delete"],
  member: ["create", "read", "update", "delete"],
  invitation: ["create", "read", "update", "delete", "cancel"],
  team: ["create", "read", "update", "delete"],
  emailTemplate: ["create", "read", "update", "delete"],
  emailLog: ["read"],
  currency: ["read"],
  organizationCurrency: ["create", "read", "update", "delete"],
  account: ["create", "read", "update", "delete"],
  category: ["create", "read", "update", "delete"],
  budget: ["create", "read", "update", "delete"],
  transaction: ["create", "read", "update", "delete"],
});

export const admin = ac.newRole({
  organization: ["update"],
  member: ["create", "read", "update", "delete"],
  invitation: ["create", "read", "update", "delete", "cancel"],
  team: ["create", "read", "update", "delete"],
  emailTemplate: ["create", "read", "update", "delete"],
  emailLog: ["read"],
  currency: ["read"],
  organizationCurrency: ["create", "read", "update", "delete"],
  account: ["create", "read", "update", "delete"],
  category: ["create", "read", "update", "delete"],
  budget: ["create", "read", "update", "delete"],
  transaction: ["create", "read", "update", "delete"],
});

export const member = ac.newRole({
  member: ["read"],
  invitation: ["read"],
  team: ["read"],
  currency: ["read"],
  organizationCurrency: ["read"],
  account: ["read"],
  category: ["read"],
  budget: ["read"],
  transaction: ["read"],
});
