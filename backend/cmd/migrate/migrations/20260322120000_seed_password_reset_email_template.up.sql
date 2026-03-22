-- Seed the system-level password reset email template (not duplicated per org).
INSERT INTO notifications.email_templates (organization_id, event, name, description, subject, content, is_active, locale, is_organization_template)
VALUES (
    NULL,
    'user.password_reset',
    'Password Reset',
    'Sent to users to reset their password',
    'Reset your password',
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
            <h2 style="margin:0 0 16px;color:#18181b;font-size:20px;">Reset your password</h2>
            <p style="margin:0 0 16px;color:#3f3f46;font-size:16px;line-height:1.5;">
              Hi {{.Name}}, we received a request to reset the password for your account (<strong>{{.Email}}</strong>). Click the button below to set a new password.
            </p>
            <table cellpadding="0" cellspacing="0" style="margin:24px 0;">
              <tr><td style="background-color:#18181b;border-radius:6px;padding:12px 24px;">
                <a href="{{.ResetURL}}" style="color:#ffffff;text-decoration:none;font-size:16px;font-weight:600;">Reset Password</a>
              </td></tr>
            </table>
            <p style="margin:0 0 16px;color:#71717a;font-size:14px;line-height:1.5;">
              If the button doesn''t work, copy and paste this link into your browser:
            </p>
            <p style="margin:0 0 16px;color:#71717a;font-size:14px;word-break:break-all;">{{.ResetURL}}</p>
            <p style="margin:0;color:#71717a;font-size:14px;line-height:1.5;">
              If you didn''t request a password reset, you can safely ignore this email.
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
    'en',
    false
);
