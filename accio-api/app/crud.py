from sqlalchemy.orm import Session
from . import models, schemas

def get_user_by_sub(db: Session, keycloak_sub: str):
    return db.query(models.User).filter(models.User.keycloak_sub == keycloak_sub).first()

def create_user(db: Session, user: schemas.UserCreate):
    db_user = models.User(email=user.email, keycloak_sub=user.keycloak_sub)
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    return db_user

def create_credential(db: Session, credential: schemas.CredentialCreate, user_id: str):
    db_cred = models.Credential(**credential.dict(), owner_id=user_id)
    db.add(db_cred)
    db.commit()
    db.refresh(db_cred)
    return db_cred

def get_credentials(db: Session, user_id: str, limit: int = 100):
    return db.query(models.Credential).filter(models.Credential.owner_id == user_id).limit(limit).all()

def create_stack(db: Session, stack: schemas.StackCreate, user_id: str):
    db_stack = models.Stack(**stack.dict(), owner_id=user_id)
    db.add(db_stack)
    db.commit()
    db.refresh(db_stack)
    return db_stack

def get_stacks(db: Session, user_id: str, limit: int = 100):
    return db.query(models.Stack).filter(models.Stack.owner_id == user_id).limit(limit).all()

def get_stack(db: Session, stack_id: str, user_id: str):
    return db.query(models.Stack).filter(models.Stack.id == stack_id, models.Stack.owner_id == user_id).first()
