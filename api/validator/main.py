from fastapi import FastAPI
import httpx
import os

app = FastAPI()

# this will be the url of the core go api
core_api = os.environ['CORE_API']

@app.get("/ping")
async def root():
    r = httpx.get(f'http://{core_api}:3000/ping')
    return {"message": r.text}