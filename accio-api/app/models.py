import uuid
from sqlalchemy import Column, String, ForeignKey, DateTime, JSON, Enum as SAEnum
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import relationship
from datetime import datetime
import enum
from .database import Base

class StackStatus(str, enum.Enum):
    PENDING_GENERATION = "pending_generation"
    SYNCING = "syncing"
    RUNNING = "running"
    UPDATING = "updating"
    DESTROYING = "destroying"
    DESTROYED = "destroyed"
    ERROR = "error"

class CredentialType(str, enum.Enum):
    AWS = "aws"
    GIT = "git"
    GCP = "gcp"
    AZURE = "azure"

class User(Base):
    __tablename__ = "users"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    email = Column(String, unique=True, index=True)
    keycloak_sub = Column(String, unique=True, index=True)  # Subject ID from Keycloak
    created_at = Column(DateTime, default=datetime.utcnow)

    credentials = relationship("Credential", back_populates="owner")
    stacks = relationship("Stack", back_populates="owner")

class Credential(Base):
    __tablename__ = "credentials"

    id = Column(String, primary_key=True)  # User facing ID e.g. aws-cred-xyz
    owner_id = Column(UUID(as_uuid=True), ForeignKey("users.id"))
    type = Column(SAEnum(CredentialType))
    encrypted_payload = Column(String) # For simplicity using String here, eventually BYTEA
    kms_key_id = Column(String)
    created_at = Column(DateTime, default=datetime.utcnow)

    owner = relationship("User", back_populates="credentials")

class Stack(Base):
    __tablename__ = "stacks"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    name = Column(String, index=True)
    owner_id = Column(UUID(as_uuid=True), ForeignKey("users.id"))
    analysis_id = Column(String)
    spec = Column(JSON)
    status = Column(SAEnum(StackStatus), default=StackStatus.PENDING_GENERATION)
    aws_credential_id = Column(String, ForeignKey("credentials.id"), nullable=True)
    git_credential_id = Column(String, ForeignKey("credentials.id"), nullable=True)
    target_region = Column(String)
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)

    owner = relationship("User", back_populates="stacks")
    versions = relationship("StackVersion", back_populates="stack")

class StackVersion(Base):
    __tablename__ = "stack_versions"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    stack_id = Column(UUID(as_uuid=True), ForeignKey("stacks.id"))
    version_number = Column(String) # Integer ideally
    spec_snapshot = Column(JSON)
    change_reason = Column(String)
    created_at = Column(DateTime, default=datetime.utcnow)

    stack = relationship("Stack", back_populates="versions")
