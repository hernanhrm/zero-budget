-- Seed the system-level organization invitation email template.
INSERT INTO notifications.email_templates (organization_id, event, name, description, subject, content, is_active, locale, is_organization_template)
VALUES (
    NULL,
    'organization.invitation_created',
    'Organization Invitation',
    'Sent when a user is invited to join an organization',
    'You''ve been invited to join {{.OrganizationName}} on Zero Budget',
    '<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@700&family=IBM+Plex+Mono:wght@400&display=swap" rel="stylesheet">
</head>
<body style="margin:0;padding:0;background-color:#1A1A1A;font-family:''IBM Plex Mono'',monospace;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#1A1A1A;padding:40px 0;">
    <tr><td align="center">
      <table width="600" cellpadding="0" cellspacing="0" style="padding:0 20px;">
        <!-- Header -->
        <tr>
          <td align="center" style="background-color:#232323;height:56px;">
            <span style="color:#FFD600;font-family:''Space Grotesk'',sans-serif;font-size:18px;font-weight:700;letter-spacing:2px;">ZERO BUDGET</span>
          </td>
        </tr>
        <!-- Card -->
        <tr>
          <td style="background-color:#232323;border:1px solid #2d2d2d;border-top:none;padding:40px 32px;">
            <table width="100%" cellpadding="0" cellspacing="0">
              <!-- Icon Circle -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <table cellpadding="0" cellspacing="0"><tr>
                    <td align="center" style="background-color:#FFD600;width:64px;height:64px;border-radius:50%;font-family:''Space Grotesk'',sans-serif;font-size:28px;font-weight:700;line-height:64px;color:#1A1A1A;">
                      {{.InviterInitial}}
                    </td>
                  </tr></table>
                </td>
              </tr>
              <!-- Heading -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <h1 style="margin:0;color:#F5F5F0;font-family:''Space Grotesk'',sans-serif;font-size:24px;font-weight:700;letter-spacing:1px;">YOU''VE BEEN INVITED!</h1>
                </td>
              </tr>
              <!-- Description -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <p style="margin:0;color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:13px;letter-spacing:1px;line-height:1.5;text-align:center;">
                    <strong style="color:#F5F5F0;">{{.InviterName}}</strong> has invited you to join <strong style="color:#F5F5F0;">{{.OrganizationName}}</strong> on Zero Budget.
                  </p>
                </td>
              </tr>
              <!-- Inviter Info Block -->
              <tr>
                <td style="padding-bottom:24px;">
                  <table width="100%" cellpadding="0" cellspacing="0"><tr>
                    <td style="background-color:#2d2d2d;padding:16px;">
                      <table cellpadding="0" cellspacing="0"><tr>
                        <td style="width:40px;vertical-align:middle;">
                          <table cellpadding="0" cellspacing="0"><tr>
                            <td align="center" style="background-color:#FFD600;width:36px;height:36px;border-radius:50%;color:#1A1A1A;font-family:''Space Grotesk'',sans-serif;font-size:14px;font-weight:700;line-height:36px;">
                              {{.InviterInitial}}
                            </td>
                          </tr></table>
                        </td>
                        <td style="padding-left:12px;vertical-align:middle;">
                          <span style="display:block;color:#F5F5F0;font-family:''Space Grotesk'',sans-serif;font-size:13px;font-weight:700;letter-spacing:1px;">{{.InviterName}}</span>
                          <a style="display:block;color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:11px;letter-spacing:1px;text-decoration:none;">{{.InviterEmail}}</a>
                        </td>
                      </tr></table>
                    </td>
                  </tr></table>
                </td>
              </tr>
              <!-- Accept Button -->
              <tr>
                <td align="center" style="padding-bottom:16px;">
                  <table cellpadding="0" cellspacing="0"><tr>
                    <td align="center" style="background-color:#FFD600;width:220px;height:48px;border-radius:0;">
                      <a href="{{.AcceptURL}}" style="display:block;line-height:48px;color:#1A1A1A;font-family:''Space Grotesk'',sans-serif;font-size:11px;font-weight:700;letter-spacing:1px;text-decoration:none;">ACCEPT INVITATION</a>
                    </td>
                  </tr></table>
                </td>
              </tr>
              <!-- Decline Link -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <a href="{{.DeclineURL}}" style="color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:11px;letter-spacing:1px;text-decoration:underline;">DECLINE INVITATION</a>
                </td>
              </tr>
              <!-- Alt Text -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <p style="margin:0;color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:12px;letter-spacing:1px;text-align:center;">
                    Or copy and paste this link into your browser:
                  </p>
                </td>
              </tr>
              <!-- Link Box -->
              <tr>
                <td style="padding-bottom:24px;">
                  <table width="100%" cellpadding="0" cellspacing="0"><tr>
                    <td align="center" style="background-color:#2d2d2d;height:40px;padding:0 12px;">
                      <a href="{{.AcceptURL}}" style="color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:11px;letter-spacing:1px;word-break:break-all;text-decoration:none;">{{.AcceptURL}}</a>
                    </td>
                  </tr></table>
                </td>
              </tr>
              <!-- Divider -->
              <tr>
                <td style="padding-bottom:24px;">
                  <table width="100%" cellpadding="0" cellspacing="0"><tr>
                    <td style="background-color:#2d2d2d;height:1px;font-size:0;line-height:0;">&nbsp;</td>
                  </tr></table>
                </td>
              </tr>
              <!-- Expiry Note -->
              <tr>
                <td align="center">
                  <p style="margin:0;color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:12px;letter-spacing:1px;line-height:1.5;text-align:center;">
                    This invitation will expire in 7 days. If you didn''t expect this invitation, you can safely ignore this email.
                  </p>
                </td>
              </tr>
            </table>
          </td>
        </tr>
        <!-- Footer -->
        <tr>
          <td style="background-color:#232323;padding:20px 32px;">
            <table width="100%" cellpadding="0" cellspacing="0">
              <tr>
                <td align="center" style="padding-bottom:8px;">
                  <span style="color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:11px;letter-spacing:1px;">&copy; 2026 ZERO BUDGET. ALL RIGHTS RESERVED.</span>
                </td>
              </tr>
              <tr>
                <td align="center">
                  <span style="color:#F5F5F0;font-family:''IBM Plex Mono'',monospace;font-size:11px;letter-spacing:1px;">PRIVACY POLICY &nbsp;&bull;&nbsp; TERMS OF SERVICE &nbsp;&bull;&nbsp; UNSUBSCRIBE</span>
                </td>
              </tr>
            </table>
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
