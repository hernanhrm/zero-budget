DROP POLICY IF EXISTS email_templates_org_scope ON notifications.email_templates;
DROP POLICY IF EXISTS email_logs_org_scope ON notifications.email_logs;

CREATE POLICY email_templates_read ON notifications.email_templates FOR SELECT USING (
    organization_id IS NULL OR
    EXISTS (
        SELECT 1 FROM identity.members m
        WHERE m.organization_id = notifications.email_templates.organization_id
        AND m.user_id = current_setting('app.current_user_id', true)
    )
);

CREATE POLICY email_logs_read ON notifications.email_logs FOR SELECT USING (
    organization_id IS NOT NULL AND
    EXISTS (
        SELECT 1 FROM identity.members m
        WHERE m.organization_id = notifications.email_logs.organization_id
        AND m.user_id = current_setting('app.current_user_id', true)
    )
);
