-- Update API Routes for workspace members
-- Remove old workspace_members routes if they exist
DELETE FROM auth.api_routes WHERE path LIKE '/workspace-members%';

-- Add workspace member routes
DO $$
DECLARE
    p_mem_read UUID;
    p_mem_manage UUID;
BEGIN
    SELECT id INTO p_mem_read FROM auth.permissions WHERE slug = 'members.read';
    SELECT id INTO p_mem_manage FROM auth.permissions WHERE slug = 'members.manage';

    INSERT INTO auth.api_routes (method, path, permission_id) VALUES
    ('GET', '/workspaces/:slug/members', p_mem_read),
    ('POST', '/workspaces/:slug/members', p_mem_manage),
    ('PUT', '/workspaces/:slug/members', p_mem_manage),
    ('DELETE', '/workspaces/:slug/members', p_mem_manage)
    ON CONFLICT (method, path) DO NOTHING;
END $$;
