from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from keycloak import KeycloakOpenID
from .database import get_db, SessionLocal
from .models import User
from sqlalchemy.orm import Session
import os

# Configuration
KEYCLOAK_URL = os.getenv("KEYCLOAK_URL", "http://localhost:8888/")
REALM_NAME = os.getenv("KEYCLOAK_REALM", "accio")
CLIENT_ID = os.getenv("KEYCLOAK_CLIENT_ID", "accio-cli")

oauth2_scheme = OAuth2PasswordBearer(tokenUrl=f"{KEYCLOAK_URL}realms/{REALM_NAME}/protocol/openid-connect/token")

# Initialize Keycloak Client (for validation)
keycloak_openid = KeycloakOpenID(
    server_url=KEYCLOAK_URL,
    client_id=CLIENT_ID,
    realm_name=REALM_NAME,
    verify=True
)

def get_current_user(token: str = Depends(oauth2_scheme), db: Session = Depends(get_db)):
    if token == "dummy":
        # Mock user for dev
        return get_mock_user(db)
        
    try:
        # In production use verify_token with public key
        # For dev, we might trust intospect or decode
        # keycloak_openid.decode_token would verify signature if public key is set or retrieved
        
        # Simple decode for now (Verify signature logic should be added)
        # options = {"verify_signature": False, "verify_aud": False, "exp": True}
        # token_info = keycloak_openid.decode_token(token, key=KEYCLOAK_PUBLIC_KEY, options=options)
        
        # Or use introspection
        user_info = keycloak_openid.userinfo(token)
        
        # Ensure user exists in local DB
        keycloak_sub = user_info.get("sub")
        email = user_info.get("email")
        
        user = db.query(User).filter(User.keycloak_sub == keycloak_sub).first()
        if not user:
            user = User(email=email, keycloak_sub=keycloak_sub)
            db.add(user)
            db.commit()
            db.refresh(user)
            
        return user
        
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail=f"Could not validate credentials: {str(e)}",
            headers={"WWW-Authenticate": "Bearer"},
        )

def get_mock_user(db: Session):
    keycloak_sub = "mock-sub-123"
    email = "mock@accio.dev"
    user = db.query(User).filter(User.keycloak_sub == keycloak_sub).first()
    if not user:
        user = User(email=email, keycloak_sub=keycloak_sub)
        db.add(user)
        db.commit()
        db.refresh(user)
    return user
