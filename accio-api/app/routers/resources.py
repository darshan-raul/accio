from fastapi import APIRouter, HTTPException, Depends
from pydantic import BaseModel
from typing import List, Dict, Any

router = APIRouter()

class ResourceCreate(BaseModel):
    provider: str
    type: str
    config: Dict[str, Any]

class Resource(ResourceCreate):
    id: str
    status: str

@router.get("/", response_model=List[Resource])
def list_resources():
    """
    List all resources managed by Accio.
    """
    return []

@router.post("/", response_model=Resource)
def create_resource(resource: ResourceCreate):
    """
    Request creation of a new resource. This triggers the MCP 'create_resource' tool.
    """
    # TODO: Trigger MCP tool
    return {
        "id": "res-123",
        "provider": resource.provider,
        "type": resource.type,
        "config": resource.config,
        "status": "pending_approval"
    }

@router.get("/{resource_id}", response_model=Resource)
def get_resource(resource_id: str):
    """
    Get details of a specific resource.
    """
    return {
        "id": resource_id,
        "provider": "aws",
        "type": "s3",
        "config": {},
        "status": "deployed"
    }

@router.delete("/{resource_id}")
def delete_resource(resource_id: str):
    """
    Request deletion of a resource.
    """
    return {"status": "deletion_requested", "pr_url": "http://github.com/org/repo/pull/124"}
