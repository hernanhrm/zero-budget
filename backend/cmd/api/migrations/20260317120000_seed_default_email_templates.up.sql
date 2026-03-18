-- A. Seed the system-level welcome email template (workspace_id = NULL).
INSERT INTO notifications.email_templates (workspace_id, event, name, description, subject, content, is_active, locale)
VALUES (
    NULL,
    'user.signed_up',
    'Welcome Email',
    'Sent to new users after sign-up',
    'Welcome to Zero Budget, {{.FirstName}}!',
    '<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f4f4f5;font-family:sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f5;padding:40px 0;">
    <tr><td align="center">
      <table width="600" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:8px;overflow:hidden;">
        <tr>
          <td style="background-color:#18181b;padding:32px;text-align:center;">
            <h1 style="margin:0;color:#ffffff;font-size:24px;">Zero Budget</h1>
          </td>
        </tr>
        <tr>
          <td style="padding:32px;">
            <h2 style="margin:0 0 16px;color:#18181b;font-size:20px;">Welcome, {{.FirstName}}!</h2>
            <p style="margin:0 0 16px;color:#3f3f46;font-size:16px;line-height:1.5;">
              Your account (<strong>{{.Email}}</strong>) is ready. You can now start tracking your finances with Zero Budget.
            </p>
            <p style="margin:0;color:#3f3f46;font-size:16px;line-height:1.5;">
              If you have any questions, just reply to this email &mdash; we are happy to help.
            </p>
          </td>
        </tr>
        <tr>
          <td style="padding:24px 32px;background-color:#f4f4f5;text-align:center;">
            <p style="margin:0;color:#71717a;font-size:12px;">&copy; Zero Budget</p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>',
    true,
    'en'
);

-- B. Function that copies all system templates (workspace_id IS NULL) into a new workspace.
CREATE OR REPLACE FUNCTION notifications.copy_system_templates_to_workspace()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
    INSERT INTO notifications.email_templates
        (workspace_id, event, name, description, subject, content, is_active, locale)
    SELECT
        NEW.id, event, name, description, subject, content, is_active, locale
    FROM notifications.email_templates
    WHERE workspace_id IS NULL;

    RETURN NEW;
END;
$$;

-- C. Fire the copy function after every new workspace is created.
CREATE TRIGGER trg_copy_system_templates
    AFTER INSERT ON auth.workspaces
    FOR EACH ROW
    EXECUTE FUNCTION notifications.copy_system_templates_to_workspace();
