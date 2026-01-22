from fastapi import FastAPI
from api.v1 import auth, workspaces

app = FastAPI(title="Nexus Tasks API")

app.include_router(auth.router, prefix="/auth", tags=["Authentication"])
app.include_router(workspaces.router, prefix="/workspaces", tags=["Workspaces"])

@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "backend"}
