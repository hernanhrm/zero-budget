DELETE FROM notifications.email_templates
WHERE event = 'user.verification_email'
  AND organization_id IS NULL;
