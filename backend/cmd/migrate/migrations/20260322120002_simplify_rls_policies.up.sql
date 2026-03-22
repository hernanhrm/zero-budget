-- Replace permission-based RLS policies with simple organization scope filters.
-- Authorization is handled by Better Auth's hasPermission via the API middleware.

DROP POLICY IF EXISTS email_templates_read ON notifications.email_templates;
DROP POLICY IF EXISTS email_logs_read ON notifications.email_logs;

CREATE POLICY email_templates_org_scope ON notifications.email_templates
  FOR ALL USING (
    organization_id IS NULL OR
    organization_id = current_setting('app.current_organization_id', true)
  );

CREATE POLICY email_logs_org_scope ON notifications.email_logs
  FOR ALL USING (
    organization_id = current_setting('app.current_organization_id', true)
  );
