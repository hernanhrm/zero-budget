-- Drop the trigger first, then the function.
DROP TRIGGER IF EXISTS trg_copy_system_templates ON identity.organizations;
DROP FUNCTION IF EXISTS notifications.copy_system_templates_to_organization();

-- Delete organization-level copies that were created from system templates.
DELETE FROM notifications.email_templates
WHERE organization_id IS NOT NULL
  AND event IN (
      SELECT event FROM notifications.email_templates WHERE organization_id IS NULL
  );

-- Delete the system-level seed template(s).
DELETE FROM notifications.email_templates WHERE organization_id IS NULL;
