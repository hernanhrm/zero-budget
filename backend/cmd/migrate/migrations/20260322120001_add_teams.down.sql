ALTER TABLE identity.invitations DROP COLUMN IF EXISTS team_id;
ALTER TABLE identity.sessions DROP COLUMN IF EXISTS active_team_id;

DROP TABLE IF EXISTS identity.team_members;
DROP TABLE IF EXISTS identity.teams;
