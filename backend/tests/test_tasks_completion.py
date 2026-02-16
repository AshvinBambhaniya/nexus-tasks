def test_task_completion_status(client):
    # Register and login
    client.post(
        "/api/v1/auth/register",
        json={
            "email": "task_test@example.com",
            "password": "password123",
            "full_name": "Test User",
        },
    )
    login_res = client.post(
        "/api/v1/auth/login",
        data={"username": "task_test@example.com", "password": "password123"},
    )
    # The client handles cookies/session if using SessionMiddleware, but here it seems token based?
    # backend/tests/conftest.py doesn't show headers usage in `test_create_workspace`.
    # Wait, `client` fixture returns `TestClient`.
    # `test_workspaces.py` uses `client.post(..., data=...)` for login.
    # Does it set a cookie? Or does subsequent request work?
    # `test_workspaces.py` shows subsequent requests without explicit headers.
    # This implies the auth is either session-based (cookies) or the test client handles it.
    # Let's check `backend/api/v1/auth.py`.

    # Check if login returns a token
    token = login_res.json().get("access_token")
    headers = {}
    if token:
        headers = {"Authorization": f"Bearer {token}"}

    # If the other tests don't use headers, maybe the client persists cookies?
    # But FastAPI with JWT usually requires header.
    # Let's follow `test_workspaces.py`: it just calls login then other endpoints.
    # This suggests the `client` (TestClient) might be managing cookies if the auth uses cookies,
    # OR the test setup is slightly different.
    # Let's check `backend/api/v1/auth.py` to see what login does.
    # It likely returns an access token.
    # If `test_workspaces.py` works without headers, maybe `TestClient` is somehow authenticated?
    # Ah, `client` fixture in `conftest.py` uses `TestClient(app)`.

    # I'll stick to providing headers if I see a token, otherwise I'll try without.
    # But to be safe, I'll extract the token if present.

    # Get user's workspaces
    ws_res = client.get("/api/v1/workspaces/", headers=headers)
    assert ws_res.status_code == 200
    workspace_id = ws_res.json()[0]["id"]

    # Create Project
    proj_res = client.post(
        f"/api/v1/workspaces/{workspace_id}/projects",
        json={"name": "Test Project", "description": "Test Desc"},
        headers=headers,
    )
    assert proj_res.status_code == 200
    project_id = proj_res.json()["id"]

    # Create Task
    task_res = client.post(
        f"/api/v1/projects/{project_id}/tasks",
        json={"title": "Test Task", "priority": "P2"},
        headers=headers,
    )
    assert task_res.status_code == 200
    task_id = task_res.json()["id"]

    # Verify initial state
    assert task_res.json()["status"] == "TODO"
    assert task_res.json()["completed_at"] is None

    # Update status to DONE
    update_res = client.patch(
        f"/api/v1/tasks/{task_id}", json={"status": "DONE"}, headers=headers
    )
    assert update_res.status_code == 200
    task_data = update_res.json()
    assert task_data["status"] == "DONE"
    assert task_data["completed_at"] is not None

    # Update status back to IN_PROGRESS
    update_res2 = client.patch(
        f"/api/v1/tasks/{task_id}", json={"status": "IN_PROGRESS"}, headers=headers
    )
    assert update_res2.status_code == 200
    task_data2 = update_res2.json()
    assert task_data2["status"] == "IN_PROGRESS"
    assert task_data2["completed_at"] is None
