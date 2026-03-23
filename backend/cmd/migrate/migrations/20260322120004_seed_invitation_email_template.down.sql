DELETE FROM notifications.email_templates
WHERE event = 'organization.invitation_created' AND organization_id IS NULL;
