CREATE TABLE identity.teams (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    organization_id TEXT NOT NULL REFERENCES identity.organizations(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX teams_organization_id_idx ON identity.teams(organization_id);

CREATE TABLE identity.team_members (
    id TEXT PRIMARY KEY,
    team_id TEXT NOT NULL REFERENCES identity.teams(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX team_members_team_id_idx ON identity.team_members(team_id);
CREATE INDEX team_members_user_id_idx ON identity.team_members(user_id);

ALTER TABLE identity.sessions ADD COLUMN active_team_id TEXT;
ALTER TABLE identity.invitations ADD COLUMN team_id TEXT;
