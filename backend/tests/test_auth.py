def test_register_user(client):
    response = client.post(
        "/api/v1/auth/register",
        json={"email": "test@example.com", "password": "password123"},
    )
    assert response.status_code == 200
    data = response.json()
    assert data["email"] == "test@example.com"
    assert "id" in data


def test_register_duplicate_email(client):
    client.post(
        "/api/v1/auth/register",
        json={"email": "duplicate@example.com", "password": "password123"},
    )
    response = client.post(
        "/api/v1/auth/register",
        json={"email": "duplicate@example.com", "password": "newpassword"},
    )
    assert response.status_code == 400
    assert response.json()["detail"] == "Email already registered"


def test_login_user(client):
    # Register
    client.post(
        "/api/v1/auth/register",
        json={
            "email": "login@example.com",
            "password": "password123",
            "full_name": "Test User",
        },
    )

    # Login (OAuth2 uses form data)
    response = client.post(
        "/api/v1/auth/login",
        data={"username": "login@example.com", "password": "password123"},
    )
    assert response.status_code == 200
    assert response.json() == {"message": "Login successful"}
    assert "access_token" in response.cookies


def test_login_invalid_credentials(client):
    client.post(
        "/api/v1/auth/register",
        json={"email": "wrong@example.com", "password": "password123"},
    )
    response = client.post(
        "/api/v1/auth/login",
        data={"username": "wrong@example.com", "password": "wrongpassword"},
    )
    assert response.status_code == 401


def test_get_me(client):
    # Register and login
    client.post(
        "/api/v1/auth/register",
        json={
            "email": "me@example.com",
            "password": "password123",
            "full_name": "Me User",
        },
    )
    client.post(
        "/api/v1/auth/login",
        data={"username": "me@example.com", "password": "password123"},
    )

    response = client.get("/api/v1/auth/me")
    assert response.status_code == 200
    data = response.json()
    assert data["email"] == "me@example.com"
    assert data["full_name"] == "Me User"


def test_logout(client):
    # Register and login
    client.post(
        "/api/v1/auth/register",
        json={"email": "logout@example.com", "password": "password123"},
    )
    client.post(
        "/api/v1/auth/login",
        data={"username": "logout@example.com", "password": "password123"},
    )
    assert "access_token" in client.cookies

    response = client.post("/api/v1/auth/logout")
    assert response.status_code == 200
    assert response.json() == {"message": "Logout successful"}
    assert "access_token" not in client.cookies


def test_update_me(client):
    # Register and login
    client.post(
        "/api/v1/auth/register",
        json={
            "email": "update@example.com",
            "password": "password123",
            "full_name": "Old Name",
        },
    )
    client.post(
        "/api/v1/auth/login",
        data={"username": "update@example.com", "password": "password123"},
    )

    response = client.patch(
        "/api/v1/auth/me",
        json={"full_name": "New Name"},
    )
    assert response.status_code == 200
    assert response.json()["full_name"] == "New Name"

    # Update email
    response = client.patch(
        "/api/v1/auth/me",
        json={"email": "newemail@example.com"},
    )
    assert response.status_code == 200
    assert response.json()["email"] == "newemail@example.com"
