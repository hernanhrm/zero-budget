DROP POLICY IF EXISTS email_logs_read ON notifications.email_logs;
DROP POLICY IF EXISTS email_templates_read ON notifications.email_templates;
DROP TABLE IF EXISTS notifications.email_logs;
DROP TABLE IF EXISTS notifications.email_templates;
DROP SCHEMA IF EXISTS notifications;
