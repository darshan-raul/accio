from fastapi import APIRouter, WebSocket, WebSocketDisconnect
from pydantic import BaseModel
from typing import List

router = APIRouter()

class ChatMessage(BaseModel):
    message: str
    conversation_id: str = None

@router.post("/")
async def chat(request: ChatMessage):
    """
    Send a message to the AI and get a streaming response (simulated for now).
    """
    # TODO: Integrate with LLM and MCP
    return {"response": f"Echo: {request.message}", "conversation_id": request.conversation_id or "new_id"}

@router.websocket("/ws")
async def websocket_chat(websocket: WebSocket):
    await websocket.accept()
    try:
        while True:
            data = await websocket.receive_text()
            # Simulate streaming response
            await websocket.send_text(f"Received: {data}")
            await websocket.send_text("Thinking...")
            await websocket.send_text(f"Done processing: {data}")
    except WebSocketDisconnect:
        print("Client disconnected")
