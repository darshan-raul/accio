from fastapi import FastAPI, Depends, HTTPException
from sqlalchemy.orm import Session
from typing import List
from . import models, schemas, crud, auth, database, schemas_analysis

# Initialize Database
models.Base.metadata.create_all(bind=database.engine)

app = FastAPI(title="Accio API", version="0.1.0")

@app.get("/health")
def health_check():
    return {"status": "healthy", "service": "accio-api"}

@app.get("/api/v1/me", response_model=schemas.User)
def read_users_me(current_user: models.User = Depends(auth.get_current_user)):
    return current_user

@app.post("/api/v1/credentials", response_model=schemas.Credential)
def create_credential(
    credential: schemas.CredentialCreate, 
    db: Session = Depends(database.get_db),
    current_user: models.User = Depends(auth.get_current_user)
):
    return crud.create_credential(db=db, credential=credential, user_id=current_user.id)

@app.get("/api/v1/credentials", response_model=List[schemas.Credential])
def read_credentials(
    skip: int = 0, 
    limit: int = 100, 
    db: Session = Depends(database.get_db),
    current_user: models.User = Depends(auth.get_current_user)
):
    return crud.get_credentials(db=db, user_id=current_user.id, limit=limit)

@app.post("/api/v1/stacks", response_model=schemas.Stack)
def create_stack(
    stack: schemas.StackCreate, 
    db: Session = Depends(database.get_db),
    current_user: models.User = Depends(auth.get_current_user)
):
    # TODO: Trigger analysis logic or template generation here via microservices
    return crud.create_stack(db=db, stack=stack, user_id=current_user.id)

@app.get("/api/v1/stacks", response_model=List[schemas.Stack])
def read_stacks(
    skip: int = 0, 
    limit: int = 100, 
    db: Session = Depends(database.get_db),
    current_user: models.User = Depends(auth.get_current_user)
):
    return crud.get_stacks(db=db, user_id=current_user.id, limit=limit)

@app.get("/api/v1/stacks/{stack_id}", response_model=schemas.Stack)
def read_stack(
    stack_id: str, 
    db: Session = Depends(database.get_db),
    current_user: models.User = Depends(auth.get_current_user)
):
    db_stack = crud.get_stack(db, stack_id=stack_id, user_id=current_user.id)
    if db_stack is None:
        raise HTTPException(status_code=404, detail="Stack not found")
    return db_stack

@app.post("/api/v1/analyze", response_model=schemas_analysis.AnalysisResponse)
def analyze_project(
    analysis: schemas_analysis.ProjectAnalysis,
    db: Session = Depends(database.get_db),
    current_user: models.User = Depends(auth.get_current_user)
):
    # In a real implementation:
    # 1. Store analysis in DB (Create Analysis table)
    # 2. Trigger LLM/RAG pipeline to generate recommendations
    # 3. Return an ID
    
    import uuid
    analysis_id = str(uuid.uuid4())
    
    # Mock response for now
    return {
        "analysis_id": analysis_id,
        "summary": f"Detected {analysis.language} project using {analysis.framework}",
        "recommendation_ready": True
    }
