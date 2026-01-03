# ACCIO: AI Cloud Platform

Accio is an MCP-based AI Cloud Platform that enables LLMs to manage cloud infrastructure across AWS, Azure, and GCP using a strict GitOps workflow.

## Components

- **accio-cli**: Go-based CLI tool for interaction.
- **accio-api**: FastAPI server for platform logic and auth.
- **accio-mcp**: MCP Server providing tools for the AI.
- **accio-terraform**: Terraform modules and templates.

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for CLI)
- Python 3.10+ (for API/MCP)
- Terraform

### Setup

1. **Start Infrastructure**:
   ```bash
   make up
   ```
   This spins up Keycloak, Redis, and Postgres.

2. **Run API Server**:
   ```bash
   cd accio-api
   pip install -r requirements.txt
   uvicorn app.main:app --reload
   ```

3. **Run MCP Server**:
   ```bash
   cd accio-mcp
   pip install -r requirements.txt
   python server.py
   ```

4. **Build CLI**:
   ```bash
   cd accio-cli
   go build -o accio
   ./accio login
   ```

## Documentation

See generic documentation in `docs/`.