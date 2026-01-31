CREATE SCHEMA IF NOT EXISTS auth;
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. Users
CREATE TABLE auth.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT,
    image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Organizations
CREATE TABLE auth.organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    owner_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Workspaces
CREATE TABLE auth.workspaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES auth.organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(organization_id, slug)
);

-- 4. Dynamic RBAC (Roles & Permissions)
CREATE TABLE auth.permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(100) NOT NULL UNIQUE, -- e.g., 'workspace.read', 'member.invite'
    description TEXT
);

CREATE TABLE auth.roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES auth.workspaces(id) ON DELETE CASCADE, -- NULL implies a System/Global Template Role
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE auth.role_permissions (
    role_id UUID NOT NULL REFERENCES auth.roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES auth.permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- 5. API Route Permissions (Mapping Endpoints to Permissions)
CREATE TABLE auth.api_routes (
    method VARCHAR(10) NOT NULL, -- GET, POST, PUT, DELETE, PATCH
    path VARCHAR(255) NOT NULL, -- e.g. '/workspaces/:slug'
    permission_id UUID NOT NULL REFERENCES auth.permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (method, path)
);

-- 6. Workspace Members (Linked to Roles)
CREATE TABLE auth.workspace_members (
    workspace_id UUID NOT NULL REFERENCES auth.workspaces(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES auth.roles(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (workspace_id, user_id)
);

-- Indexes
CREATE INDEX idx_users_email ON auth.users(email);
CREATE INDEX idx_workspaces_organization_id ON auth.workspaces(organization_id);
CREATE INDEX idx_roles_workspace_id ON auth.roles(workspace_id);
CREATE INDEX idx_workspace_members_user_id ON auth.workspace_members(user_id);
CREATE INDEX idx_workspace_members_role_id ON auth.workspace_members(role_id);
CREATE INDEX idx_api_routes_permission_id ON auth.api_routes(permission_id);

-- Seed Basic Permissions
INSERT INTO auth.permissions (slug, description) VALUES
    ('workspace.read', 'View workspace details and content'),
    ('workspace.create', 'Create a new workspace'),
    ('workspace.update', 'Update workspace settings'),
    ('workspace.delete', 'Delete workspace'),
    ('members.read', 'View members list'),
    ('members.manage', 'Invite or remove members'),
    ('roles.manage', 'Create, update, and delete custom roles');

-- Seed API Routes Mapping
DO $$
DECLARE
    p_read UUID;
    p_create UUID;
    p_update UUID;
    p_delete UUID;
    p_mem_read UUID;
    p_mem_manage UUID;
    p_roles_manage UUID;
BEGIN
    SELECT id INTO p_read FROM auth.permissions WHERE slug = 'workspace.read';
    SELECT id INTO p_create FROM auth.permissions WHERE slug = 'workspace.create';
    SELECT id INTO p_update FROM auth.permissions WHERE slug = 'workspace.update';
    SELECT id INTO p_delete FROM auth.permissions WHERE slug = 'workspace.delete';
    SELECT id INTO p_mem_read FROM auth.permissions WHERE slug = 'members.read';
    SELECT id INTO p_mem_manage FROM auth.permissions WHERE slug = 'members.manage';
    SELECT id INTO p_roles_manage FROM auth.permissions WHERE slug = 'roles.manage';

    INSERT INTO auth.api_routes (method, path, permission_id) VALUES
    ('GET', '/workspaces', p_read),
    ('GET', '/workspaces/:slug', p_read),
    ('POST', '/workspaces', p_create),
    ('PUT', '/workspaces/:slug', p_update),
    ('DELETE', '/workspaces/:slug', p_delete),
    ('GET', '/workspaces/:slug/members', p_mem_read),
    ('POST', '/workspaces/:slug/members', p_mem_manage),
    ('DELETE', '/workspaces/:slug/members/:user-id', p_mem_manage),
    ('POST', '/workspaces/:slug/roles', p_roles_manage),
    ('PUT', '/workspaces/:slug/roles/:role-id', p_roles_manage),
    ('DELETE', '/workspaces/:slug/roles/:role-id', p_roles_manage);
END $$;

-- Seed System Roles (Templates)
DO $$
DECLARE
    p_read UUID;
    p_create UUID;
    p_update UUID;
    p_delete UUID;
    p_mem_read UUID;
    p_mem_manage UUID;
    p_roles_manage UUID;
    r_admin UUID;
    r_viewer UUID;
BEGIN
    -- Get Permission IDs
    SELECT id INTO p_read FROM auth.permissions WHERE slug = 'workspace.read';
    SELECT id INTO p_create FROM auth.permissions WHERE slug = 'workspace.create';
    SELECT id INTO p_update FROM auth.permissions WHERE slug = 'workspace.update';
    SELECT id INTO p_delete FROM auth.permissions WHERE slug = 'workspace.delete';
    SELECT id INTO p_mem_read FROM auth.permissions WHERE slug = 'members.read';
    SELECT id INTO p_mem_manage FROM auth.permissions WHERE slug = 'members.manage';
    SELECT id INTO p_roles_manage FROM auth.permissions WHERE slug = 'roles.manage';

    -- Create 'Admin' Role (System level)
    INSERT INTO auth.roles (name, description, workspace_id) 
    VALUES ('Admin', 'Full access', NULL) RETURNING id INTO r_admin;

    INSERT INTO auth.role_permissions (role_id, permission_id) VALUES 
    (r_admin, p_read), (r_admin, p_create), (r_admin, p_update), (r_admin, p_delete), 
    (r_admin, p_mem_read), (r_admin, p_mem_manage), (r_admin, p_roles_manage);

    -- Create 'Viewer' Role (System level)
    INSERT INTO auth.roles (name, description, workspace_id) 
    VALUES ('Viewer', 'Read only', NULL) RETURNING id INTO r_viewer;

    INSERT INTO auth.role_permissions (role_id, permission_id) VALUES 
    (r_viewer, p_read), (r_viewer, p_mem_read);
END $$;


-- RLS Enablement
ALTER TABLE auth.users ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.organizations ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.workspaces ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.workspace_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.role_permissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.permissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.api_routes ENABLE ROW LEVEL SECURITY;

-- Helper function to check permission by ID
CREATE OR REPLACE FUNCTION auth.check_permission(p_user_id UUID, p_workspace_id UUID, p_perm_slug TEXT)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 
        FROM auth.workspace_members wm
        JOIN auth.role_permissions rp ON wm.role_id = rp.role_id
        JOIN auth.permissions p ON rp.permission_id = p.id
        WHERE wm.user_id = p_user_id
        AND wm.workspace_id = p_workspace_id
        AND p.slug = p_perm_slug
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;


-- RLS Policies

-- Users: Read/Update self
CREATE POLICY users_read_self ON auth.users FOR SELECT USING (id = current_setting('app.current_user_id', true)::UUID);
CREATE POLICY users_update_self ON auth.users FOR UPDATE USING (id = current_setting('app.current_user_id', true)::UUID);

-- Organizations: Owner Full Access
CREATE POLICY organizations_owner_all ON auth.organizations USING (owner_id = current_setting('app.current_user_id', true)::UUID);
-- Organizations: Read if member of any workspace
CREATE POLICY organizations_read_member ON auth.organizations FOR SELECT USING (
    EXISTS (
        SELECT 1 FROM auth.workspaces w
        JOIN auth.workspace_members wm ON w.id = wm.workspace_id
        WHERE w.organization_id = auth.organizations.id
        AND wm.user_id = current_setting('app.current_user_id', true)::UUID
    )
);

-- Workspaces: Org Owner Full Access
CREATE POLICY workspaces_org_owner_all ON auth.workspaces USING (
    EXISTS (
        SELECT 1 FROM auth.organizations o
        WHERE o.id = auth.workspaces.organization_id
        AND o.owner_id = current_setting('app.current_user_id', true)::UUID
    )
);
-- Workspaces: Read Access (Requires 'workspace.read' permission)
CREATE POLICY workspaces_member_read ON auth.workspaces FOR SELECT USING (
    auth.check_permission(current_setting('app.current_user_id', true)::UUID, id, 'workspace.read')
);
-- Workspaces: Update Access (Requires 'workspace.update' permission)
CREATE POLICY workspaces_member_update ON auth.workspaces FOR UPDATE USING (
    auth.check_permission(current_setting('app.current_user_id', true)::UUID, id, 'workspace.update')
);
-- Workspaces: Delete Access (Requires 'workspace.delete' permission)
CREATE POLICY workspaces_member_delete ON auth.workspaces FOR DELETE USING (
    auth.check_permission(current_setting('app.current_user_id', true)::UUID, id, 'workspace.delete')
);


-- Members: Read Self
CREATE POLICY workspace_members_read_self ON auth.workspace_members FOR SELECT USING (user_id = current_setting('app.current_user_id', true)::UUID);
-- Members: Manage (Requires 'members.manage' permission)
CREATE POLICY workspace_members_manage ON auth.workspace_members USING (
    auth.check_permission(current_setting('app.current_user_id', true)::UUID, workspace_id, 'members.manage')
);
-- Members: Org Owner Full Access
CREATE POLICY workspace_members_org_owner ON auth.workspace_members USING (
    EXISTS (
        SELECT 1 FROM auth.workspaces w
        JOIN auth.organizations o ON w.organization_id = o.id
        WHERE w.id = auth.workspace_members.workspace_id
        AND o.owner_id = current_setting('app.current_user_id', true)::UUID
    )
);

-- Roles: Read Access (System Roles OR Workspace Roles user is part of)
CREATE POLICY roles_read ON auth.roles FOR SELECT USING (
    workspace_id IS NULL OR -- System Role
    EXISTS ( -- User is member of the workspace the role belongs to
        SELECT 1 FROM auth.workspace_members wm
        WHERE wm.workspace_id = auth.roles.workspace_id
        AND wm.user_id = current_setting('app.current_user_id', true)::UUID
    )
);
-- Roles: Manage Custom Roles (Requires 'roles.manage')
CREATE POLICY roles_manage ON auth.roles USING (
    workspace_id IS NOT NULL AND -- Only custom roles
    auth.check_permission(current_setting('app.current_user_id', true)::UUID, workspace_id, 'roles.manage')
);


-- Role Permissions: Read Access (If user can see the Role, they can see its perms)
CREATE POLICY role_permissions_read ON auth.role_permissions FOR SELECT USING (
    EXISTS (
        SELECT 1 FROM auth.roles r
        WHERE r.id = auth.role_permissions.role_id
        AND (
            r.workspace_id IS NULL OR
            EXISTS (
                SELECT 1 FROM auth.workspace_members wm
                WHERE wm.workspace_id = r.workspace_id
                AND wm.user_id = current_setting('app.current_user_id', true)::UUID
            )
        )
    )
);
-- Role Permissions: Manage (Requires 'roles.manage' on the role's workspace)
CREATE POLICY role_permissions_manage ON auth.role_permissions USING (
    EXISTS (
        SELECT 1 FROM auth.roles r
        WHERE r.id = auth.role_permissions.role_id
        AND r.workspace_id IS NOT NULL -- Only custom roles
        AND auth.check_permission(current_setting('app.current_user_id', true)::UUID, r.workspace_id, 'roles.manage')
    )
);


-- Permissions: Public Read (Required for system to function)
CREATE POLICY permissions_read_all ON auth.permissions FOR SELECT USING (true);

-- API Routes: Public Read
CREATE POLICY api_routes_read_all ON auth.api_routes FOR SELECT USING (true);
