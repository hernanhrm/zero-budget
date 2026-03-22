import { createAccessControl } from "better-auth/plugins/access";

const statement = {
  organization: ["update", "delete"],
  member: ["create", "read", "update", "delete"],
  invitation: ["create", "read", "update", "delete"],
  team: ["create", "read", "update", "delete"],
  emailTemplate: ["create", "read", "update", "delete"],
  emailLog: ["read"],
} as const;

export const ac = createAccessControl(statement);

export const owner = ac.newRole({
  organization: ["update", "delete"],
  member: ["create", "read", "update", "delete"],
  invitation: ["create", "read", "update", "delete"],
  team: ["create", "read", "update", "delete"],
  emailTemplate: ["create", "read", "update", "delete"],
  emailLog: ["read"],
});

export const admin = ac.newRole({
  organization: ["update"],
  member: ["create", "read", "update", "delete"],
  invitation: ["create", "read", "update", "delete"],
  team: ["create", "read", "update", "delete"],
  emailTemplate: ["create", "read", "update", "delete"],
  emailLog: ["read"],
});

export const member = ac.newRole({
  member: ["read"],
  invitation: ["read"],
  team: ["read"],
});
