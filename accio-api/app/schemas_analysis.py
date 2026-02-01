from pydantic import BaseModel
from typing import List, Optional
from .schemas import StackStatus

class ProjectAnalysis(BaseModel):
    language: Optional[str] = None
    framework: Optional[str] = None
    dependencies: List[str] = []
    database: Optional[str] = None
    cache: Optional[str] = None
    files: List[str] = []

class AnalysisResponse(BaseModel):
    analysis_id: str
    summary: str
    recommendation_ready: bool
