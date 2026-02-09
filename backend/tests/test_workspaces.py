def test_create_workspace(client):
    # Register and login
    client.post(
        "/api/v1/auth/register",
        json={
            "email": "ws@example.com",
            "password": "password123",
            "full_name": "WS User",
        },
    )
    client.post(
        "/api/v1/auth/login",
        data={"username": "ws@example.com", "password": "password123"},
    )

    # Check default workspace
    response = client.get("/api/v1/workspaces/")
    assert response.status_code == 200
    workspaces = response.json()
    assert len(workspaces) == 1
    assert workspaces[0]["name"] == "Personal"
    assert workspaces[0]["type"] == "PERSONAL"

    # Create new team workspace
    response = client.post(
        "/api/v1/workspaces/",
        json={"name": "New Team WS"},
    )
    assert response.status_code == 200
    data = response.json()
    assert data["name"] == "New Team WS"
    assert data["type"] == "TEAM"

    # List again
    response = client.get("/api/v1/workspaces/")
    assert len(response.json()) == 2


def test_list_workspace_members(client):
    # Register two users
    client.post(
        "/api/v1/auth/register",
        json={
            "email": "user1@example.com",
            "password": "password123",
            "full_name": "User One",
        },
    )
    client.post(
        "/api/v1/auth/register",
        json={
            "email": "user2@example.com",
            "password": "password123",
            "full_name": "User Two",
        },
    )

    # Login as User One
    client.post(
        "/api/v1/auth/login",
        data={"username": "user1@example.com", "password": "password123"},
    )

    # Create a workspace
    response = client.post(
        "/api/v1/workspaces/",
        json={"name": "Shared WS"},
    )
    workspace_id = response.json()["id"]

    # Invite User Two
    response = client.post(
        f"/api/v1/workspaces/{workspace_id}/members",
        json={"email": "user2@example.com"},
    )
    assert response.status_code == 200

    # List members
    response = client.get(f"/api/v1/workspaces/{workspace_id}/members")
    assert response.status_code == 200
    members = response.json()
    assert len(members) == 2
    emails = [m["user"]["email"] for m in members]
    assert "user1@example.com" in emails
    assert "user2@example.com" in emails
