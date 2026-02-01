-- Revert API Routes changes
DELETE FROM auth.api_routes WHERE path LIKE '/workspaces/:slug/members';
