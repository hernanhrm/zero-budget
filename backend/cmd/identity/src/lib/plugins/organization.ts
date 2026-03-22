import { organization } from "better-auth/plugins";
import { ac, owner, admin, member } from "../../permissions.js";

export const organizationPlugin = organization({
  ac,
  roles: {
    owner,
    admin,
    member,
  },
  allowUserToCreateOrganization: true,
  teams: {
    enabled: true,
  },
  dynamicAccessControl: {
    enabled: true,
  },
  schema: {
    organization: {
      fields: {
        createdAt: "created_at",
      },
    },
    member: {
      fields: {
        organizationId: "organization_id",
        userId: "user_id",
        createdAt: "created_at",
      },
    },
    invitation: {
      fields: {
        organizationId: "organization_id",
        inviterId: "inviter_id",
        expiresAt: "expires_at",
        teamId: "team_id",
      },
    },
    team: {
      fields: {
        organizationId: "organization_id",
        createdAt: "created_at",
        updatedAt: "updated_at",
      },
    },
    teamMember: {
      fields: {
        teamId: "team_id",
        userId: "user_id",
        createdAt: "created_at",
      },
    },
    organizationRole: {
      fields: {
        organizationId: "organization_id",
        createdAt: "created_at",
        updatedAt: "updated_at",
      },
    },
  },
});
