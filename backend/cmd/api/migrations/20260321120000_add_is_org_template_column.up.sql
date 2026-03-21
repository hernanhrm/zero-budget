-- A. Add the is_organization_template column to distinguish templates that
--    should be copied into every new organization from system-only templates.
ALTER TABLE notifications.email_templates
    ADD COLUMN is_organization_template BOOLEAN NOT NULL DEFAULT false;

-- B. Mark the existing welcome email as an organization template.
UPDATE notifications.email_templates
SET is_organization_template = true
WHERE event = 'user.signed_up'
  AND organization_id IS NULL;

-- C. Replace the trigger function so it only copies templates flagged as
--    organization templates.
CREATE OR REPLACE FUNCTION notifications.copy_system_templates_to_organization()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
    INSERT INTO notifications.email_templates
        (organization_id, event, name, description, subject, content, is_active, locale, is_organization_template)
    SELECT
        NEW.id, event, name, description, subject, content, is_active, locale, is_organization_template
    FROM notifications.email_templates
    WHERE organization_id IS NULL
      AND is_organization_template = true;

    RETURN NEW;
END;
$$;
