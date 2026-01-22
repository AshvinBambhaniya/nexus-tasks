from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from api.v1 import auth, workspaces
import os

app = FastAPI(title="Nexus Tasks API")

# CORS Configuration
origins = [
    "http://localhost:3000",  # Next.js Frontend
    os.getenv("FRONTEND_URL", ""), # Production URL from env
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=[origin for origin in origins if origin],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(auth.router, prefix="/api/v1/auth", tags=["Authentication"])
app.include_router(workspaces.router, prefix="/api/v1/workspaces", tags=["Workspaces"])

@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "backend"}
