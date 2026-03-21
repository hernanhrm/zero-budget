-- A. Restore the original trigger function (without is_organization_template filter).
CREATE OR REPLACE FUNCTION notifications.copy_system_templates_to_organization()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
    INSERT INTO notifications.email_templates
        (organization_id, event, name, description, subject, content, is_active, locale)
    SELECT
        NEW.id, event, name, description, subject, content, is_active, locale
    FROM notifications.email_templates
    WHERE organization_id IS NULL;

    RETURN NEW;
END;
$$;

-- B. Drop the column.
ALTER TABLE notifications.email_templates DROP COLUMN is_organization_template;
