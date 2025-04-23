import os
from google.adk.agents import LlmAgent
from google.adk.models.lite_llm import LiteLlm
from google.adk.tools.mcp_tool.mcp_toolset import MCPToolset, SseServerParams
from google.adk.runners import Runner
from google.adk.sessions import InMemorySessionService
from google.genai import types
import asyncio

async def get_tools_async(server_url):
    # Правильный адрес для подключения как клиент
    
    params = SseServerParams(
        url=f"{server_url}/events",            # SSE-поток
        headers={"Accept": "text/event-stream"},
        message_url=f"{server_url}/messages"   # HTTP-POST endpoint
    )
    tools, exit_stack = await MCPToolset.from_server(connection_params=params)
    return tools, exit_stack

async def main():
    tools, exit_stack = await get_tools_async()
    # Get the MCP server URL from environment variable, defaulting to localhost:8082
    # server_url = "http://0.0.0.0:8082"
    server_url = "http://localhost:8082"

    
    try:
        tools, exit_stack = await get_tools_async(server_url)
        agent = LlmAgent(
            name="wb_data_analyzer",
            model=LiteLlm(model="openai/gpt-4"),  # Fixed model name
            tools=tools
        )
        session_service = InMemorySessionService()
        session_service.create_session(app_name="wb_data_app", user_id="user_1", session_id="sess1")
        runner = Runner(app_name="wb_data_app", agent=agent, session_service=session_service)

        content = types.Content(role="user", parts=[types.Part(text="Analyze sales data for the last week")])
        async for event in runner.run_async(user_id="user_1", session_id="sess1", new_message=content):
            if event.is_final_response():
                print(event.content.parts[0].text)
                break
    finally:
        if 'exit_stack' in locals():
            await exit_stack.aclose()

if __name__ == "__main__":
    asyncio.run(main())
