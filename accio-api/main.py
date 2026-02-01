from fastapi import FastAPI

app = FastAPI(title="Accio API", version="0.1.0")

@app.get("/health")
def health_check():
    return {"status": "healthy", "service": "accio-api"}

@app.get("/")
def read_root():
    return {"message": "Welcome to Accio Platform API"}
