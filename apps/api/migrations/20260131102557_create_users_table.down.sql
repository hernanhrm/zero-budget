DROP POLICY IF EXISTS api_routes_read_all ON auth.api_routes;
DROP POLICY IF EXISTS roles_read ON auth.roles;
DROP POLICY IF EXISTS workspace_members_org_owner ON auth.workspace_members;
DROP POLICY IF EXISTS workspace_members_manage ON auth.workspace_members;
DROP POLICY IF EXISTS workspace_members_read_self ON auth.workspace_members;
DROP POLICY IF EXISTS workspaces_member_read ON auth.workspaces;
DROP POLICY IF EXISTS workspaces_org_owner_all ON auth.workspaces;
DROP POLICY IF EXISTS organizations_read_member ON auth.organizations;
DROP POLICY IF EXISTS organizations_owner_all ON auth.organizations;
DROP POLICY IF EXISTS users_update_self ON auth.users;
DROP POLICY IF EXISTS users_read_self ON auth.users;

DROP FUNCTION IF EXISTS auth.check_permission;

DROP TABLE IF EXISTS auth.workspace_members;
DROP TABLE IF EXISTS auth.api_routes;
DROP TABLE IF EXISTS auth.role_permissions;
DROP TABLE IF EXISTS auth.roles;
DROP TABLE IF EXISTS auth.permissions;
DROP TABLE IF EXISTS auth.workspaces;
DROP TABLE IF EXISTS auth.organizations;
DROP TABLE IF EXISTS auth.users;

DROP EXTENSION IF EXISTS "pgcrypto";
DROP SCHEMA IF EXISTS auth;
