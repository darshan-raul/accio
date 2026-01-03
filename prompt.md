# ACCIO: MCP-Based AI Cloud Platform - Comprehensive Build Prompt

## Role & Expertise
You are a **Senior Solutions Architect and DevOps Engineer** with deep expertise in:
- Model Context Protocol (MCP) implementation and best practices
- Multi-cloud infrastructure (AWS, Azure, GCP)
- GitOps methodologies and Infrastructure as Code (Terraform)
- Distributed systems architecture
- Identity and Access Management (Keycloak, OAuth2/OIDC)
- CLI development in Go
- AI/LLM integration patterns

## Project Overview
Build **ACCIO** - an MCP-based AI Cloud Platform that enables LLMs to manage cloud infrastructure across AWS, Azure, and GCP using a strict GitOps workflow. The AI acts as an infrastructure assistant that generates Terraform code and creates Pull Requests rather than applying changes directly.

---

## Architecture Components

### 1. **CLI Tool (Golang)**
**Requirements:**
- Interactive CLI using `cobra` and `bubbletea` for TUI components
- Support both natural language conversation mode and structured command mode
- Commands needed:
  - `accio login` - OAuth2/OIDC flow with Keycloak
  - `accio chat` - Interactive AI conversation mode
  - `accio resource create <type>` - Structured resource creation
  - `accio resource list/get/delete` - Resource management
  - `accio status` - Check platform and PR status
  - `accio config` - Manage cloud credentials and repo settings

**Implementation Details:**
```go
// Example structure
type AccioCLI struct {
    APIClient    *APIClient
    AuthProvider *KeycloakAuth
    Config       *Config
    TUI          *TUIManager
}
```

**Features:**
- Store auth tokens securely using OS keyring
- Support for multiple profiles (dev, staging, prod)
- Rich TUI for showing resource states, PR status, and approval workflows
- Streaming LLM responses in chat mode
- Syntax highlighting for generated Terraform code preview

### 2. **API Server (Python/FastAPI or Go/Gin)**
**Requirements:**
- RESTful API with WebSocket support for streaming LLM responses
- Endpoints:
  - `POST /api/v1/auth/login` - Initiate OAuth flow
  - `POST /api/v1/auth/callback` - OAuth callback handler
  - `POST /api/v1/chat` - Send message to AI (streaming response)
  - `GET/POST /api/v1/resources` - CRUD operations
  - `GET /api/v1/pull-requests` - List PRs created by platform
  - `POST /api/v1/cloud-consent` - Grant cloud provider access
  - `GET /api/v1/status` - Platform health and user status

**Architecture Improvements:**
- Add rate limiting per user
- Implement request validation middleware
- Add audit logging for all operations
- Use Redis for caching and session management
- Implement circuit breakers for external service calls

### 3. **MCP Server (Python with official SDK)**
**Core Implementation:**

```python
from mcp.server import Server, McpServerTool
from mcp.types import Tool, TextContent

class AccioMCPServer(Server):
    def __init__(self):
        super().__init__("accio-mcp-server")
        self.terraform_generator = TerraformGenerator()
        self.git_manager = GitManager()
        self.cloud_clients = {
            'aws': AWSClient(),
            'azure': AzureClient(), 
            'gcp': GCPClient()
        }
```

**Tools to Implement:**

1. **`analyze_infrastructure`**
   - Scan existing resources across clouds
   - Return current state and recommendations
   - Input: cloud_provider, resource_type, filters
   - Output: JSON with resource inventory

2. **`create_resource`**
   - Generate Terraform for new resource
   - Create feature branch
   - Commit changes and create PR
   - Input: provider, resource_type, configuration (JSON)
   - Output: PR URL, branch name, terraform preview

3. **`update_resource`**
   - Modify existing Terraform
   - Create PR with changes
   - Show diff before/after
   - Input: resource_id, changes (JSON)
   - Output: PR URL with diff

4. **`delete_resource`**
   - Generate destruction PR
   - Add safety checks and confirmations
   - Input: resource_id, confirmation_token
   - Output: PR URL for deletion

5. **`get_resource_status`**
   - Check actual cloud state vs. Terraform state
   - Detect drift
   - Input: resource_id
   - Output: state comparison, drift details

6. **`validate_terraform`**
   - Run terraform validate and plan
   - Return cost estimates (using Infracost)
   - Input: terraform_code
   - Output: validation results, estimated cost

7. **`list_pull_requests`**
   - Show pending infrastructure PRs
   - Filter by status, resource type
   - Input: filters
   - Output: List of PRs with metadata

**GitOps Workflow:**
```python
class GitOpsManager:
    async def create_infrastructure_pr(
        self,
        resource_config: dict,
        terraform_code: str,
        user_context: dict
    ) -> PullRequest:
        # 1. Create feature branch from main
        branch = f"accio/create-{resource_type}-{timestamp}"
        
        # 2. Generate proper directory structure
        # infrastructure/
        #   ├── aws/
        #   ├── azure/
        #   ├── gcp/
        #   └── modules/
        
        # 3. Write Terraform files
        # 4. Run terraform fmt and validate
        # 5. Generate PR description with:
        #    - What's being created
        #    - Why (from AI context)
        #    - Cost estimate
        #    - Security scan results
        
        # 6. Create PR with labels: accio-managed, cloud:aws, etc.
        # 7. Add reviewers from team
        # 8. Post summary in PR comments
        
        return pr
```

### 4. **Keycloak Integration**
**Setup Instructions for Users:**

```markdown
## Keycloak Configuration Steps

### 1. Install Keycloak
docker run -p 8080:8080 -e KEYCLOAK_ADMIN=admin \
  -e KEYCLOAK_ADMIN_PASSWORD=admin \
  quay.io/keycloak/keycloak:latest start-dev

### 2. Create Realm "accio"
- Login to admin console
- Create new realm: "accio"

### 3. Create Client "accio-cli"
- Client ID: accio-cli
- Client Protocol: openid-connect
- Access Type: public
- Standard Flow Enabled: ON
- Valid Redirect URIs: http://localhost:8888/callback
- Web Origins: http://localhost:8888

### 4. Create Client Scopes
- cloud:aws - Manage AWS resources
- cloud:azure - Manage Azure resources
- cloud:gcp - Manage GCP resources
- git:write - Create PRs in Git repositories

### 5. Configure Consent
- Consent Required: ON
- Display Client Scope Consent: ON

This ensures users explicitly grant permission for:
- Cloud provider access
- Git repository modifications
```

**Auth Flow Implementation:**
```python
class KeycloakAuthProvider:
    def __init__(self, realm_url: str, client_id: str):
        self.realm_url = realm_url
        self.client_id = client_id
        
    async def initiate_login(self) -> str:
        """Return authorization URL for user"""
        # PKCE flow for CLI
        code_verifier = generate_code_verifier()
        code_challenge = generate_code_challenge(code_verifier)
        
        params = {
            'client_id': self.client_id,
            'response_type': 'code',
            'scope': 'openid profile cloud:aws cloud:azure cloud:gcp git:write',
            'redirect_uri': 'http://localhost:8888/callback',
            'code_challenge': code_challenge,
            'code_challenge_method': 'S256',
            'consent': 'required'
        }
        
        return f"{self.realm_url}/auth?{urlencode(params)}"
    
    async def handle_callback(self, code: str, state: str) -> Tokens:
        """Exchange code for tokens"""
        # Exchange authorization code for tokens
        # Store refresh token securely
        # Return access token for API calls
```

**Consent Screen Configuration:**
```yaml
# Keycloak should show:
Accio Platform requests permission to:
☐ Access your AWS account to manage resources
☐ Access your Azure account to manage resources  
☐ Access your GCP account to manage resources
☐ Create pull requests in your Git repositories
☐ Read your profile information

These permissions allow the AI assistant to generate and 
propose infrastructure changes on your behalf through 
GitOps pull requests.
```

### 5. **Terraform Code Generation**
**Template System:**
```python
class TerraformGenerator:
    def __init__(self):
        self.templates = {
            'aws': {
                'ec2': EC2Template(),
                's3': S3Template(),
                'rds': RDSTemplate(),
            },
            'azure': {...},
            'gcp': {...}
        }
    
    def generate(self, provider: str, resource_type: str, config: dict) -> str:
        template = self.templates[provider][resource_type]
        
        # Generate with best practices:
        # - Tags for all resources
        # - Encryption by default
        # - Least privilege IAM
        # - Network security
        
        tf_code = template.render(config)
        
        # Add to modules for reusability
        return tf_code
```

**Example Generated Terraform:**
```hcl
# Generated by ACCIO AI Platform
# Request: "Create a PostgreSQL database for production app"
# Generated: 2025-01-03 10:30:00 UTC

terraform {
  required_version = ">= 1.6.0"
  
  backend "s3" {
    bucket = "accio-terraform-state"
    key    = "prod/database/postgres.tfstate"
    region = "us-east-1"
  }
}

resource "aws_db_instance" "prod_postgres" {
  identifier     = "prod-app-postgres"
  engine         = "postgres"
  engine_version = "15.4"
  instance_class = "db.t3.medium"
  
  allocated_storage     = 100
  max_allocated_storage = 500
  storage_encrypted     = true
  
  db_name  = "prodapp"
  username = "dbadmin"
  password = data.aws_secretsmanager_secret_version.db_password.secret_string
  
  vpc_security_group_ids = [aws_security_group.db_sg.id]
  db_subnet_group_name   = aws_db_subnet_group.private.name
  
  backup_retention_period = 30
  backup_window          = "03:00-04:00"
  maintenance_window     = "mon:04:00-mon:05:00"
  
  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]
  
  tags = {
    ManagedBy   = "accio"
    Environment = "production"
    CreatedBy   = "ai-assistant"
    CreatedAt   = "2025-01-03"
  }
}
```

### 6. **Cloud Provider Clients**
**Implementation Pattern:**
```python
class CloudProviderClient(ABC):
    @abstractmethod
    async def list_resources(self, resource_type: str) -> List[Resource]:
        pass
    
    @abstractmethod
    async def get_resource_details(self, resource_id: str) -> Resource:
        pass
    
    @abstractmethod
    async def validate_permissions(self, user_context: dict) -> bool:
        pass

class AWSClient(CloudProviderClient):
    def __init__(self, credentials: AWSCredentials):
        # Use user's credentials from OAuth consent
        # Never store credentials in code
        self.session = boto3.Session(
            aws_access_key_id=credentials.access_key,
            aws_secret_access_key=credentials.secret_key,
            region_name=credentials.region
        )
```

**Credential Management:**
- After Keycloak consent, user provides cloud credentials
- Encrypted at rest using Vault or KMS
- Credentials scoped to specific permissions
- Rotation policy enforced
- Never logged or exposed in Git

---

## Improved Architecture Features

### 1. **Policy Engine**
Add a policy validation layer before creating PRs:
```python
class PolicyEngine:
    async def validate_change(self, terraform_code: str) -> PolicyResult:
        checks = [
            self.check_cost_limits(),
            self.check_security_compliance(),
            self.check_naming_conventions(),
            self.check_tagging_requirements(),
            self.check_network_rules()
        ]
        return await asyncio.gather(*checks)
```

### 2. **Cost Estimation**
Integrate Infracost:
```python
async def estimate_cost(terraform_code: str) -> CostEstimate:
    # Use Infracost API
    # Show monthly cost breakdown
    # Compare to current spending
    # Alert if over budget threshold
```

### 3. **Drift Detection**
```python
async def detect_drift() -> DriftReport:
    # Compare Terraform state vs actual cloud state
    # Schedule periodic drift detection
    # Auto-create PRs to reconcile drift
```

### 4. **Multi-Tenancy**
- Each user/team has isolated:
  - Git branches (user/team namespace)
  - Cloud accounts
  - State files
  - PR workflows

### 5. **Approval Workflows**
```yaml
# .github/workflows/accio-pr-approval.yml
name: Accio Infrastructure Approval

on:
  pull_request:
    types: [opened, synchronize]
    paths:
      - 'infrastructure/**'

jobs:
  validate:
    - terraform validate
    - terraform plan
    - security scan (tfsec)
    - cost estimate
    - require approval from team leads
```

---

## Implementation Checklist

### Phase 1: Foundation (Week 1-2)
- [ ] Set up project structure
- [ ] Implement Keycloak auth flow
- [ ] Create basic CLI with login
- [ ] Build API server skeleton
- [ ] Set up MCP server with 1-2 basic tools

### Phase 2: Core Features (Week 3-4)
- [ ] Implement all MCP tools
- [ ] Terraform generation for common resources
- [ ] GitOps workflow (branch, commit, PR)
- [ ] Cloud provider clients (AWS first)
- [ ] CLI chat interface with TUI

### Phase 3: Advanced Features (Week 5-6)
- [ ] Policy engine
- [ ] Cost estimation
- [ ] Drift detection
- [ ] Multi-cloud support (Azure, GCP)
- [ ] PR automation and webhooks

### Phase 4: Production Ready (Week 7-8)
- [ ] Comprehensive testing
- [ ] Documentation
- [ ] Security hardening
- [ ] Monitoring and observability
- [ ] Performance optimization

---

## User Journey Example

### First Time Setup
```bash
# 1. Install CLI
curl -fsSL https://accio.dev/install.sh | sh

# 2. Login
accio login
# Opens browser for Keycloak consent
# User grants cloud and git permissions
# CLI stores encrypted tokens

# 3. Configure cloud providers
accio cloud add aws --profile production
accio cloud add azure --subscription my-sub
accio cloud add gcp --project my-project

# 4. Configure git repository
accio git set-repo github.com/myorg/infrastructure
```

### Creating Infrastructure
```bash
# Natural language mode
accio chat
> I need a PostgreSQL database for our production API

AI: I'll create a production-ready PostgreSQL instance on AWS. 
    Let me generate the Terraform configuration for you.
    
    [Shows preview of Terraform code]
    
    This will create:
    - RDS PostgreSQL 15.4 instance (db.t3.medium)
    - Encrypted storage with 100GB (auto-scaling to 500GB)
    - Automated backups (30 day retention)
    - Estimated cost: $145/month
    
    Should I create a pull request for this? (yes/no)

> yes

AI: Created PR #123: "Add production PostgreSQL database"
    https://github.com/myorg/infrastructure/pull/123
    
    The PR includes:
    ✓ Terraform code
    ✓ Cost estimate
    ✓ Security scan results
    ✓ Required approvals: 2
    
    Once approved and merged, the GitHub Actions workflow 
    will apply the changes to your AWS account.
```

### Structured Command Mode
```bash
accio resource create \
  --provider aws \
  --type rds \
  --name prod-postgres \
  --engine postgres \
  --instance-class db.t3.medium
```

---

## Security Considerations

1. **Credential Storage**: Use OS keyring, never plaintext
2. **Least Privilege**: Request minimum IAM permissions
3. **Audit Logging**: Log all operations with user context
4. **Secret Scanning**: Prevent secrets in generated code
5. **Network Security**: TLS for all communications
6. **Code Review**: All changes go through PR review
7. **Rollback**: Easy rollback mechanism for bad changes

---

## Testing Strategy

```python
# MCP Tool Tests
async def test_create_resource_tool():
    result = await mcp_server.call_tool(
        "create_resource",
        {
            "provider": "aws",
            "resource_type": "s3",
            "config": {"bucket_name": "test-bucket"}
        }
    )
    
    assert "pr_url" in result
    assert "terraform_preview" in result
    
# Integration Tests
async def test_end_to_end_workflow():
    # 1. CLI login
    # 2. Send chat message
    # 3. Verify PR created
    # 4. Mock approval
    # 5. Verify resource created
    
# Load Tests
async def test_concurrent_users():
    # Simulate 100 concurrent users
    # Verify rate limiting
    # Check response times
```

---

## Deliverables

### Code Repositories
1. **accio-cli** (Go) - CLI tool
2. **accio-api** (Python/FastAPI) - API server
3. **accio-mcp** (Python) - MCP server
4. **accio-terraform** (HCL) - Terraform modules
5. **accio-docs** (Markdown) - Documentation

### Documentation
- Architecture diagrams
- User guide
- Developer guide
- API reference
- Security documentation
- Runbook for operators

### Configuration Files
- Docker Compose for local dev
- Kubernetes manifests for production
- CI/CD pipelines
- Keycloak realm export
- Example Terraform modules

---

## Success Metrics

- **User Experience**: Can create infrastructure in <2 minutes
- **Safety**: 100% of changes go through PR review
- **Reliability**: 99.9% uptime
- **Performance**: API latency <200ms (p95)
- **Security**: Zero credential leaks, all audited

---

## Generate the complete, production-ready code for this platform with:
- Detailed comments explaining each component
- Error handling and logging
- Configuration management
- Testing coverage
- Deployment instructions
- Security best practices implemented