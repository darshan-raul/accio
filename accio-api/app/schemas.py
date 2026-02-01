from pydantic import BaseModel, UUID4, EmailStr
from typing import Optional, List, Dict, Any
from datetime import datetime
from enum import Enum

class StackStatus(str, Enum):
    PENDING_GENERATION = "pending_generation"
    SYNCING = "syncing"
    RUNNING = "running"
    UPDATING = "updating"
    DESTROYING = "destroying"
    DESTROYED = "destroyed"
    ERROR = "error"

class CredentialType(str, Enum):
    AWS = "aws"
    GIT = "git"
    GCP = "gcp"
    AZURE = "azure"

class UserBase(BaseModel):
    email: EmailStr

class UserCreate(UserBase):
    keycloak_sub: str

class User(UserBase):
    id: UUID4
    keycloak_sub: str
    created_at: datetime

    class Config:
        from_attributes = True

class CredentialBase(BaseModel):
    id: str
    type: CredentialType
    encrypted_payload: str
    kms_key_id: str

class CredentialCreate(CredentialBase):
    pass

class Credential(CredentialBase):
    owner_id: UUID4
    created_at: datetime
    
    class Config:
        from_attributes = True

class StackBase(BaseModel):
    name: str
    analysis_id: str
    spec: Dict[str, Any]
    target_region: str
    aws_credential_id: Optional[str] = None
    git_credential_id: Optional[str] = None

class StackCreate(StackBase):
    pass

class Stack(StackBase):
    id: UUID4
    owner_id: UUID4
    status: StackStatus
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True
