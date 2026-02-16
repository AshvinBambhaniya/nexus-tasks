from unittest.mock import AsyncMock, patch

from tasks.email_tasks import send_workspace_invitation_email


def test_send_workspace_invitation_email():
    """
    Test that the email task correctly constructs the email and calls FastMail.
    """
    email_to = "test@example.com"
    workspace_name = "Test Workspace"

    # Mock FastMail to avoid sending real emails
    with patch("tasks.email_tasks.FastMail") as MockFastMail:
        # Create a mock instance
        mock_fm_instance = MockFastMail.return_value
        # Mock the send_message method (it's async)
        mock_fm_instance.send_message = AsyncMock()

        # Call the task directly (synchronously as it's just a python function decorated with @task)
        # Note: When testing Celery tasks directly, we can just call the function if we don't care about the Celery wrapper
        # or use .apply() for synchronous execution.
        # Since the task function itself creates an event loop to run the async send_message,
        # we can call it as a normal function.

        result = send_workspace_invitation_email(email_to, workspace_name)

        assert result == f"Email sent to {email_to}"

        # Verify FastMail was initialized
        MockFastMail.assert_called_once()

        # Verify send_message was called
        mock_fm_instance.send_message.assert_called_once()

        # Inspect the arguments passed to send_message
        call_args = mock_fm_instance.send_message.call_args
        message = call_args[0][0]  # The MessageSchema object

        assert message.subject == f"You've been added to {workspace_name}"
        assert message.recipients == [email_to]
        assert "Welcome to Nexus Tasks" in message.body
        assert workspace_name in message.body
