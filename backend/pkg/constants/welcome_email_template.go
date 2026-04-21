package constants

const WelcomeEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <title>Welcome to {{.CompanyName}}</title>
  <style>
    body { font-family: Arial, sans-serif; background: #f4f6f8; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 40px auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
    .header { background: #1d4ed8; padding: 32px; text-align: center; }
    .header h1 { color: #ffffff; margin: 0; font-size: 24px; }
    .header p { color: #bfdbfe; margin: 8px 0 0; font-size: 14px; }
    .body { padding: 32px; }
    .body h2 { color: #1e293b; font-size: 20px; margin-top: 0; }
    .body p { color: #475569; line-height: 1.7; font-size: 15px; }
    .info-box { background: #f1f5f9; border-left: 4px solid #1d4ed8; border-radius: 4px; padding: 16px 20px; margin: 24px 0; }
    .info-box table { width: 100%; border-collapse: collapse; }
    .info-box td { padding: 4px 0; color: #334155; font-size: 14px; }
    .info-box td:first-child { font-weight: bold; width: 130px; color: #1e293b; }
    .footer { background: #f1f5f9; padding: 20px 32px; text-align: center; color: #94a3b8; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Welcome to {{.CompanyName}}!</h1>
      <p>We are thrilled to have you on board.</p>
    </div>
    <div class="body">
      <h2>Dear {{.FullName}},</h2>
      <p>
        We are delighted to welcome you to the <strong>{{.CompanyName}}</strong> family.
        Your journey with us starts on <strong>{{.StartDate}}</strong>, and we look forward to achieving great things together.
      </p>
      <div class="info-box">
        <table>
          <tr><td>Name</td><td>{{.FullName}}</td></tr>
          <tr><td>Position</td><td>{{.Position}}</td></tr>
          <tr><td>Department</td><td>{{.Department}}</td></tr>
          <tr><td>Start Date</td><td>{{.StartDate}}</td></tr>
        </table>
      </div>
      <p>
        Our HR and IT teams will be in touch to help you get set up and settled in.
        If you have any questions before your first day, feel free to reach out.
      </p>
      <p>Looking forward to working with you!</p>
      <p>Warm regards,<br/><strong>{{.CompanyName}} HR Team</strong></p>
    </div>
    <div class="footer">
      &copy; {{.CompanyName}}. This email was sent automatically by the HR system.
    </div>
  </div>
</body>
</html>`
