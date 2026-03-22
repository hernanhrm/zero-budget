DELETE FROM notifications.email_templates
WHERE event = 'user.password_reset'
  AND organization_id IS NULL;
