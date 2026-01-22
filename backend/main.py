from fastapi import FastAPI
from api.v1 import auth

app = FastAPI(title="Nexus Tasks API")

app.include_router(auth.router, prefix="/auth", tags=["Authentication"])

@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "backend"}
