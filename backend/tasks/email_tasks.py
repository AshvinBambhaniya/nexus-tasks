import asyncio

from fastapi_mail import FastMail, MessageSchema, MessageType

from core.celery_app import celery_app
from core.email import conf


@celery_app.task
def send_workspace_invitation_email(email_to: str, workspace_name: str):
    """
    Celery task to send a workspace invitation email.
    """

    html = f"""
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="utf-8">
        <style>
            body {{
                font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
                line-height: 1.6;
                color: #333;
                background-color: #f9fafb;
                margin: 0;
                padding: 0;
            }}
            .container {{
                max_width: 600px;
                margin: 40px auto;
                background: #ffffff;
                border-radius: 8px;
                box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
                overflow: hidden;
            }}
            .header {{
                background-color: #0f172a; /* Dark Slate */
                color: white;
                padding: 30px;
                text-align: center;
            }}
            .header h1 {{
                margin: 0;
                font-size: 24px;
                font-weight: 600;
            }}
            .content {{
                padding: 40px 30px;
            }}
            .content p {{
                margin-bottom: 20px;
                font-size: 16px;
            }}
            .btn {{
                display: inline-block;
                background-color: #2563eb; /* Blue 600 */
                color: white;
                text-decoration: none;
                padding: 12px 24px;
                border-radius: 6px;
                font-weight: 600;
                margin-top: 10px;
            }}
            .btn:hover {{
                background-color: #1d4ed8;
            }}
            .footer {{
                background-color: #f1f5f9;
                padding: 20px;
                text-align: center;
                font-size: 14px;
                color: #64748b;
            }}
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>Welcome to Nexus Tasks</h1>
            </div>
            <div class="content">
                <p>Hello,</p>
                <p>You have been invited to join the workspace <strong>{workspace_name}</strong>.</p>
                <p>Collaborate with your team, track projects, and get things done efficiently.</p>
                <div style="text-align: center; margin: 30px 0;">
                    <a href="http://localhost:3000/dashboard" class="btn">Go to Dashboard</a>
                </div>
                <p>If you didn't expect this invitation, you can safely ignore this email.</p>
            </div>
            <div class="footer">
                &copy; 2026 Nexus Tasks. All rights reserved.
            </div>
        </div>
    </body>
    </html>
    """

    message = MessageSchema(
        subject=f"You've been added to {workspace_name}",
        recipients=[email_to],
        body=html,
        subtype=MessageType.html,
    )

    fm = FastMail(conf)

    # fastapi-mail is async, so we need to run it in the event loop
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        loop.run_until_complete(fm.send_message(message))
    finally:
        loop.close()

    return f"Email sent to {email_to}"
