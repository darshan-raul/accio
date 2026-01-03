from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from app.core import auth

router = APIRouter()

class LoginRequest(BaseModel):
    redirect_uri: str

class CallbackRequest(BaseModel):
    code: str
    redirect_uri: str

@router.post("/login")
def login(request: LoginRequest):
    """
    Returns the Keycloak login URL.
    """
    try:
        url = auth.get_auth_url(request.redirect_uri)
        return {"authorization_url": url}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/callback")
def callback(request: CallbackRequest):
    """
    Exchanges authorization code for access tokens.
    """
    try:
        token = auth.exchange_code_for_token(request.code, request.redirect_uri)
        return token
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))
