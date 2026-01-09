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
    # Using LangChain for LLM integration
    # For now, defaulting to OpenAI, but structure allows easy swapping
    # for Bedrock, Claude, etc. via LangChain's abstraction.
    from langchain_openai import ChatOpenAI
    from langchain.schema import HumanMessage, SystemMessage

    # In a real scenario, API keys should be loaded from environment variables
    llm = ChatOpenAI(temperature=0.7)

    messages = [
        SystemMessage(content="You are a helpful AI assistant."),
        HumanMessage(content=request.message),
    ]

    response = llm.invoke(messages)
    
    # For now, just return the content. In the future, this might stream or handle tools (MCP).
    return {"response": response.content, "conversation_id": request.conversation_id or "new_id"}
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
