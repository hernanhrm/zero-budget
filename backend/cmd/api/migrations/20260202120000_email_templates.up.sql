CREATE SCHEMA IF NOT EXISTS notifications;

CREATE TABLE notifications.email_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id TEXT REFERENCES identity.organizations(id) ON DELETE CASCADE,
    event VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    subject TEXT NOT NULL,
    content TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    locale VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE notifications.email_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID REFERENCES notifications.email_templates(id),
    organization_id TEXT REFERENCES identity.organizations(id),
    recipient_email TEXT NOT NULL,
    event VARCHAR(100) NOT NULL,
    subject TEXT NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_email_templates_organization_event
  ON notifications.email_templates(organization_id, event);
CREATE INDEX idx_email_templates_is_active
  ON notifications.email_templates(is_active);
CREATE INDEX idx_email_logs_organization_id
  ON notifications.email_logs(organization_id);
CREATE INDEX idx_email_logs_sent_at
  ON notifications.email_logs(sent_at DESC);

ALTER TABLE notifications.email_templates ENABLE ROW LEVEL SECURITY;
ALTER TABLE notifications.email_logs ENABLE ROW LEVEL SECURITY;

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
