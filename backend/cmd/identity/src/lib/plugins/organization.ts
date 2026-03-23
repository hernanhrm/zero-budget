import { organization } from "better-auth/plugins";
import { ac, owner, admin, member } from "../../permissions.js";
import { publishOrganizationInvitationCreated } from "../events.js";

export const organizationPlugin = organization({
  ac,
  roles: {
    owner,
    admin,
    member,
  },
  allowUserToCreateOrganization: true,
  sendInvitationEmail: async (data) => {
    const inviterName = data.inviter.user.name || data.inviter.user.email;
    const baseUrl = process.env.APP_URL || "http://localhost:3000";
    await publishOrganizationInvitationCreated({
      email: data.email,
      inviterName,
      inviterEmail: data.inviter.user.email,
      inviterInitial: inviterName.charAt(0).toUpperCase(),
      organizationName: data.organization.name,
      acceptUrl: `${baseUrl}/invite/accept/${data.id}`,
      declineUrl: `${baseUrl}/invite/decline/${data.id}`,
    });
  },
  teams: {
    enabled: true,
  },
  dynamicAccessControl: {
    enabled: true,
  },
  schema: {
    session: {
      fields: {
        activeOrganizationId: "active_organization_id",
        activeTeamId: "active_team_id",
      },
    },
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
        createdAt: "created_at",
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
