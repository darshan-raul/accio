from keycloak import KeycloakOpenID
from app.config import get_settings

settings = get_settings()

keycloak_openid = KeycloakOpenID(
    server_url=settings.KEYCLOAK_URL,
    client_id=settings.KEYCLOAK_CLIENT_ID,
    realm_name=settings.KEYCLOAK_REALM,
    client_secret_key=settings.KEYCLOAK_CLIENT_SECRET,
    verify=True
)

def get_auth_url(redirect_uri: str):
    return keycloak_openid.auth_url(
        redirect_uri=redirect_uri,
        scope="openid profile email",
        state="some_random_state" # In prod, manage state securely
    )

def exchange_code_for_token(code: str, redirect_uri: str):
    return keycloak_openid.token(
        grant_type='authorization_code',
        code=code,
        redirect_uri=redirect_uri
    )

def refresh_token(refresh_token: str):
    return keycloak_openid.refresh_token(refresh_token)

def get_user_info(token: str):
    return keycloak_openid.userinfo(token)
