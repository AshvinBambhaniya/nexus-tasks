from fastapi import FastAPI

app = FastAPI(title="Nexus Tasks API")

@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "backend"}
