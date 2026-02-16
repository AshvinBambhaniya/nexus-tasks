import os


class Settings:
    PROJECT_NAME: str = "Nexus Tasks"
    PROJECT_VERSION: str = "1.0.0"
    FRONTEND_URL: str = os.getenv("FRONTEND_URL", "http://localhost:3000")

    POSTGRES_USER: str = os.getenv("POSTGRES_USER", "postgres")
    POSTGRES_PASSWORD: str = os.getenv("POSTGRES_PASSWORD", "postgres")
    POSTGRES_DB: str = os.getenv("POSTGRES_DB", "nexus_tasks")
    POSTGRES_HOST: str = os.getenv("POSTGRES_HOST", "localhost")
    POSTGRES_PORT: str = os.getenv("POSTGRES_PORT", "5432")

    DATABASE_URL = f"postgresql+psycopg://{POSTGRES_USER}:{POSTGRES_PASSWORD}@{POSTGRES_HOST}:{POSTGRES_PORT}/{POSTGRES_DB}"

    SECRET_KEY: str = os.getenv("SECRET_KEY", "supersecretkey")
    ALGORITHM: str = "HS256"
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 30

    # Celery
    CELERY_BROKER_URL: str = os.getenv("CELERY_BROKER_URL", "redis://localhost:6379/0")
    CELERY_RESULT_BACKEND: str = os.getenv(
        "CELERY_RESULT_BACKEND", "redis://localhost:6379/0"
    )

    # Mail
    MAIL_USERNAME: str = os.getenv("MAIL_USERNAME", "")
    MAIL_PASSWORD: str = os.getenv("MAIL_PASSWORD", "")
    MAIL_FROM: str = os.getenv("MAIL_FROM", "noreply@example.com")
    MAIL_PORT: int = int(os.getenv("MAIL_PORT", 1025))
    MAIL_SERVER: str = os.getenv("MAIL_SERVER", "localhost")
    MAIL_FROM_NAME: str = os.getenv("MAIL_FROM_NAME", "Nexus Tasks")
    MAIL_STARTTLS: bool = os.getenv("MAIL_STARTTLS", "False").lower() == "true"
    MAIL_SSL_TLS: bool = os.getenv("MAIL_SSL_TLS", "False").lower() == "true"
    USE_CREDENTIALS: bool = os.getenv("USE_CREDENTIALS", "False").lower() == "true"
    VALIDATE_CERTS: bool = os.getenv("VALIDATE_CERTS", "False").lower() == "true"


settings = Settings()
