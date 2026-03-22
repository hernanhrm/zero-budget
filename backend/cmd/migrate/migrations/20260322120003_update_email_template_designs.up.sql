-- Update all 3 email templates to match the new dark-theme Pencil designs.

-- 1. Welcome Email (user.signed_up)
UPDATE notifications.email_templates
SET subject = 'Welcome to Zero Budget, {{.Name}}!',
    content = '<!DOCTYPE html>
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
                    <td align="center" style="background-color:#FFD600;width:64px;height:64px;border-radius:50%;font-size:28px;line-height:64px;color:#1A1A1A;">
                      &#10022;
                    </td>
                  </tr></table>
                </td>
              </tr>
              <!-- Heading -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <h1 style="margin:0;color:#F5F5F0;font-family:''Space Grotesk'',sans-serif;font-size:24px;font-weight:700;letter-spacing:1px;">WELCOME TO ZERO BUDGET!</h1>
                </td>
              </tr>
              <!-- Description -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <p style="margin:0;color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:13px;letter-spacing:1px;line-height:1.5;text-align:center;">
                    Your account (<a href="mailto:{{.Email}}" style="color:#6B6B6B;text-decoration:none;font-weight:700;">{{.Email}}</a>) is ready. You can now start tracking your finances with Zero Budget.
                  </p>
                </td>
              </tr>
              <!-- Button -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <table cellpadding="0" cellspacing="0"><tr>
                    <td align="center" style="background-color:#FFD600;width:220px;height:48px;border-radius:0;">
                      <a href="https://zerobudget.app" style="display:block;line-height:48px;color:#1A1A1A;font-family:''Space Grotesk'',sans-serif;font-size:11px;font-weight:700;letter-spacing:1px;text-decoration:none;">GET STARTED</a>
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
              <!-- Help Text -->
              <tr>
                <td align="center">
                  <p style="margin:0;color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:12px;letter-spacing:1px;line-height:1.5;text-align:center;">
                    If you have any questions, just reply to this email &mdash; we are happy to help.
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
    updated_at = NOW()
WHERE event = 'user.signed_up';

-- 2. Verify Email (user.verification_email)
UPDATE notifications.email_templates
SET subject = 'Verify your email address',
    content = '<!DOCTYPE html>
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
                    <td align="center" style="background-color:#FFD600;width:64px;height:64px;border-radius:50%;font-size:28px;line-height:64px;color:#1A1A1A;">
                      &#10022;
                    </td>
                  </tr></table>
                </td>
              </tr>
              <!-- Heading -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <h1 style="margin:0;color:#F5F5F0;font-family:''Space Grotesk'',sans-serif;font-size:24px;font-weight:700;letter-spacing:1px;">VERIFY YOUR EMAIL</h1>
                </td>
              </tr>
              <!-- Description -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <p style="margin:0;color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:13px;letter-spacing:1px;line-height:1.5;text-align:center;">
                    Thanks for signing up for Zero Budget! Please verify your email address by clicking the button below.
                  </p>
                </td>
              </tr>
              <!-- Button -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <table cellpadding="0" cellspacing="0"><tr>
                    <td align="center" style="background-color:#FFD600;width:220px;height:48px;border-radius:0;">
                      <a href="{{.VerificationURL}}" style="display:block;line-height:48px;color:#1A1A1A;font-family:''Space Grotesk'',sans-serif;font-size:11px;font-weight:700;letter-spacing:1px;text-decoration:none;">VERIFY EMAIL ADDRESS</a>
                    </td>
                  </tr></table>
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
                      <a href="{{.VerificationURL}}" style="color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:11px;letter-spacing:1px;word-break:break-all;text-decoration:none;">{{.VerificationURL}}</a>
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
                    This link will expire in 24 hours. If you didn''t create an account, you can safely ignore this email.
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
    updated_at = NOW()
WHERE event = 'user.verification_email';

-- 3. Reset Password (user.password_reset)
UPDATE notifications.email_templates
SET subject = 'Reset your password',
    content = '<!DOCTYPE html>
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
                    <td align="center" style="background-color:#FFD600;width:64px;height:64px;border-radius:50%;font-size:28px;line-height:64px;color:#1A1A1A;">
                      &#10022;
                    </td>
                  </tr></table>
                </td>
              </tr>
              <!-- Heading -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <h1 style="margin:0;color:#F5F5F0;font-family:''Space Grotesk'',sans-serif;font-size:24px;font-weight:700;letter-spacing:1px;">RESET YOUR PASSWORD</h1>
                </td>
              </tr>
              <!-- Description -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <p style="margin:0;color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:13px;letter-spacing:1px;line-height:1.5;text-align:center;">
                    We received a request to reset the password for your account (<a href="mailto:{{.Email}}" style="color:#6B6B6B;text-decoration:none;font-weight:700;">{{.Email}}</a>). Click the button below to set a new password.
                  </p>
                </td>
              </tr>
              <!-- Button -->
              <tr>
                <td align="center" style="padding-bottom:24px;">
                  <table cellpadding="0" cellspacing="0"><tr>
                    <td align="center" style="background-color:#FFD600;width:220px;height:48px;border-radius:0;">
                      <a href="{{.ResetURL}}" style="display:block;line-height:48px;color:#1A1A1A;font-family:''Space Grotesk'',sans-serif;font-size:11px;font-weight:700;letter-spacing:1px;text-decoration:none;">RESET PASSWORD</a>
                    </td>
                  </tr></table>
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
                      <a href="{{.ResetURL}}" style="color:#6B6B6B;font-family:''IBM Plex Mono'',monospace;font-size:11px;letter-spacing:1px;word-break:break-all;text-decoration:none;">{{.ResetURL}}</a>
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
                    This link will expire in 24 hours. If you didn''t request a password reset, you can safely ignore this email.
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
    updated_at = NOW()
WHERE event = 'user.password_reset';
