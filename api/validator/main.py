from fastapi import FastAPI
import httpx

app = FastAPI()


@app.get("/ping")
async def root():
    r = httpx.get('http://localhost:3000/ping')
    return {"message": r.text}