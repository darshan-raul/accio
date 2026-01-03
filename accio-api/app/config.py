from pydantic_settings import BaseSettings
from functools import lru_cache

class Settings(BaseSettings):
    API_V1_STR: str = "/api/v1"
    PROJECT_NAME: str = "Accio API"
    
    # Keycloak Configuration
    KEYCLOAK_URL: str = "http://localhost:8080"
    KEYCLOAK_REALM: str = "accio"
    KEYCLOAK_CLIENT_ID: str = "accio-cli"
    KEYCLOAK_CLIENT_SECRET: str = "" # Public client for CLI, but might be needed for admin ops
    
    # Dependencies
    REDIS_URL: str = "redis://localhost:6379/0"

    class Config:
        case_sensitive = True
        env_file = ".env"

@lru_cache()
def get_settings():
    return Settings()
