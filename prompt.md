# Accio

> **CLI-first, AI-assisted, GitOps-enabled cloud infrastructure platform.**
> Analyze your code, generate optimal cloud infrastructure, and let GitOps do the rest.

---

## Table of Contents

1. [Philosophy & Design Principles](#1-philosophy--design-principles)
2. [Architecture Overview](#2-architecture-overview)
3. [Prerequisites & Initial Setup](#3-prerequisites--initial-setup)
4. [Phase 1 — Authentication & Credential Binding](#phase-1--authentication--credential-binding)
5. [Phase 2 — Project Analysis](#phase-2--project-analysis)
6. [Phase 3 — Infrastructure Recommendation](#phase-3--infrastructure-recommendation)
7. [Phase 4 — Spec Editing & Interactive Refinement](#phase-4--spec-editing--interactive-refinement)
8. [Phase 5 — Stack Validation & Persistence](#phase-5--stack-validation--persistence)
9. [Phase 6 — Template Generation & Git Commit](#phase-6--template-generation--git-commit)
10. [Phase 7 — GitOps Sync (ArgoCD + Crossplane)](#phase-7--gitops-sync-argocd--crossplane)
11. [Phase 8 — Runtime Interaction & Observability](#phase-8--runtime-interaction--observability)
12. [Data Model](#data-model)
13. [Git Repository Structure](#git-repository-structure)
14. [Crossplane Composition Design](#crossplane-composition-design)
15. [Error Handling & Recovery](#error-handling--recovery)
16. [Security Model](#security-model)
17. [CLI Command Reference](#cli-command-reference)

---

## 1. Philosophy & Design Principles

Accio is built around a core set of beliefs about how cloud infrastructure should work for developers:

**Code is the source of truth for what you need. Git is the source of truth for what exists.** You never manually click through a cloud console. You never run `terraform apply` from a laptop. Every infrastructure state change is a commit, every commit is a sync, every sync is reconciled by the platform.

**AI reduces the gap between intent and implementation.** Most developers know what their application does but not what infrastructure it needs. Accio's job is to close that gap — not by guessing, but by analyzing your actual code and reasoning about it with context from its knowledge base.

**The platform should get out of your way after the initial setup.** Once your stack is running, you interact with it through the same CLI — asking questions, checking costs, monitoring health — without needing to know where the resources actually live or how they are wired together.

**Everything is namespaced. Nothing is shared by accident.** Every managed resource Crossplane creates is namespaced within a Kubernetes cluster. Teams cannot accidentally clobber each other's infrastructure.

---

## 2. Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              Developer                                   │
│                         runs `accio` CLI                                 │
└───────────────────────────────┬─────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                           Accio CLI                                      │
│  ┌──────────┐  ┌──────────────┐  ┌────────────┐  ┌──────────────────┐  │
│  │  Auth &  │  │   Code       │  │   AI       │  │  Interactive     │  │
│  │  Cred    │→ │   Analyzer   │→ │   Engine   │→ │  Spec Editor     │  │
│  │  Manager │  │  (AST/Grep)  │  │ (RAG+MCP)  │  │  (REPL loop)    │  │
│  └──────────┘  └──────────────┘  └────────────┘  └────────┬─────────┘  │
└─────────────────────────────────────────────────────────────┬───────────┘
                                                               │
                                                               ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          Accio API (Backend)                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────────┐  │
│  │   Stack      │  │  Validation  │  │     Template Generator       │  │
│  │   Service    │→ │  Engine      │→ │  (Crossplane Compositions)   │  │
│  └──────┬───────┘  └──────────────┘  └──────────────┬───────────────┘  │
│         │                                            │                   │
│         ▼                                            ▼                   │
│  ┌──────────────┐                         ┌──────────────────┐         │
│  │  PostgreSQL  │                         │   Git Client     │         │
│  │  (Stacks DB) │                         │  (commits YAML)  │         │
│  └──────────────┘                         └────────┬─────────┘         │
└─────────────────────────────────────────────────────┼───────────────────┘
                                                       │ push
                                                       ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         Git Repository                                   │
│                    (infra compositions YAML)                             │
└───────────────────────────────┬─────────────────────────────────────────┘
                                │ watched by
                                ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                                     │
│  ┌─────────────┐       ┌──────────────────┐                            │
│  │   ArgoCD    │──────▶│  Crossplane      │                            │
│  │  (GitOps    │ sync  │  (Resource       │                            │
│  │   sync)     │       │   Orchestration) │                            │
│  └─────────────┘       └────────┬─────────┘                            │
│                                 │ creates/manages                       │
│                                 ▼                                        │
│                    ┌────────────────────────┐                           │
│                    │   Cloud Provider       │                           │
│                    │  (AWS / GCP / Azure)   │                           │
│                    └────────────────────────┘                           │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Prerequisites & Initial Setup

Before Accio can do anything, the following must exist:

**Local machine requirements:**
- A terminal (any modern shell: bash, zsh, fish)
- `accio` CLI installed (see install instructions in your package manager or `go install`)
- Network access to your cloud provider and to your git hosting provider

**External accounts and services (set up once, referenced throughout):**
- A **Keycloak** instance (or an OIDC-compatible identity provider) configured with at least one Accio client application. This is what the CLI authenticates against. Accio does not manage users itself — it delegates entirely to your identity provider.
- An **AWS account** (or GCP / Azure — this document uses AWS as the reference provider). You will need an IAM user or role with an Access Key ID and Secret Access Key scoped to the permissions Accio needs (see [Security Model](#security-model)).
- A **Git hosting account** (GitHub, GitLab, Bitbucket, or self-hosted Gitea/Gogs). You will need a Personal Access Token (PAT) with permissions to create repositories and push commits.
- A **Kubernetes cluster** with ArgoCD and Crossplane installed and configured. This is the cluster that will actually reconcile your infrastructure. It can be a managed cluster (EKS, GKE, AKS) or self-managed. Crossplane must have the appropriate provider packages installed (e.g., `crossplane-contrib/provider-aws`) and configured with credentials that allow it to create cloud resources.
- A **PostgreSQL database** accessible by the Accio API backend. This stores stack objects and their history.

---

## Phase 1 — Authentication & Credential Binding

**What happens here:** The CLI establishes your identity and securely associates cloud and git credentials with your session.

### 1.1 — Login via Keycloak

```bash
accio login
```

The CLI opens a browser window (or prints a device code if running headless) pointed at your Keycloak authorization endpoint. You authenticate with your username and password (or SSO if configured). Keycloak issues an OIDC token. The CLI stores the token locally in an encrypted credential store (e.g., the OS keychain or a local encrypted file at `~/.accio/tokens`).

All subsequent API calls from the CLI include this token as a Bearer token in the Authorization header. The Accio API validates it against the Keycloak JWKS endpoint. Token refresh happens transparently when the access token expires.

```
$ accio login
Opening browser for authentication...
Waiting for authentication...
✓ Authenticated as sarah@acme.dev
✓ Token valid until 2025-08-15T14:32:00Z
```

### 1.2 — Credential Registration

Credentials are registered once per account and stored encrypted in the Accio API backend (not on the local machine beyond the initial input).

```bash
accio credentials add aws --access-key-id AKIA... --secret-access-key ...
accio credentials add git --pat ghp_...
```

The CLI sends these credentials to the Accio API over TLS. The API encrypts them at rest using an envelope encryption scheme (a data key encrypted by a KMS key, the data key then used to encrypt the credential payload). The CLI itself never stores the secret after submission.

```
$ accio credentials add aws --access-key-id AKIAIOSFODNN7EXAMPLE
? Secret Access Key: (hidden input)
✓ AWS credentials registered and encrypted (credential ID: aws-cred-a1b2c3)

$ accio credentials add git --pat ghp_example1234567890
✓ Git PAT registered and encrypted (credential ID: git-cred-d4e5f6)
```

You can list, rotate, or revoke credentials at any time:

```bash
accio credentials list
accio credentials rotate aws-cred-a1b2c3
accio credentials revoke git-cred-d4e5f6
```

---

## Phase 2 — Project Analysis

**What happens here:** Accio reads your source code and builds an internal model of what your application is, what it depends on, and what kind of infrastructure it likely needs.

### 2.1 — Granting Permission

Accio does not automatically read your code. You must explicitly opt in per project. This is a deliberate design choice — your source code may contain business logic, internal APIs, or sensitive patterns that should not be sent anywhere without your knowledge.

```bash
cd /path/to/my-project
accio analyze
```

The CLI detects that this directory has not been previously authorized and prompts:

```
? This command will analyze the source code in /path/to/my-project.
  Accio will read file contents, dependency manifests, and configuration files.
  No source code is stored permanently — it is used only for the duration of this analysis.
  Do you grant permission? (y/N)
```

If you answer yes, the analysis proceeds. If you answer no, Accio exits cleanly.

### 2.2 — What the Analyzer Reads

The analyzer is not arbitrary — it follows a defined priority and scope:

**Tier 1 — Project structure and metadata (always read):**
- Directory tree (depth-limited, ignores `node_modules`, `.git`, build artifacts)
- `package.json`, `go.mod`, `pom.xml`, `Pipfile`, `Cargo.toml`, or equivalent dependency manifests
- `Dockerfile` or `docker-compose.yml` if present
- `.env.example` or `.env.sample` (not `.env` itself — Accio never reads actual secret files)
- `README.md` at the project root

**Tier 2 — Application entry points and routing (read if Tier 1 indicates a web application):**
- Main/entry files (e.g., `main.go`, `app.py`, `server.js`, `Program.cs`)
- Route definitions and API handler registrations
- Configuration files that declare ports, hosts, database connection strings (with secrets redacted)

**Tier 3 — Dependency and runtime patterns (read selectively based on Tier 1 and 2 findings):**
- Import/require statements across the codebase to identify external service dependencies (e.g., AWS SDK calls, Redis clients, message queue producers/consumers)
- Database migration files or ORM model definitions to understand data layer requirements
- Any infrastructure-as-code files already present (`*.tf`, `*.yaml` in a `k8s/` or `infra/` directory) — these are treated as hints, not overridden without explicit user confirmation

### 2.3 — Analysis Output

The analyzer produces a structured summary that it feeds into the AI engine. This summary is also displayed to you so you can see exactly what Accio understood about your project:

```
$ accio analyze
Analyzing /path/to/my-project...

✓ Language: Go (primary), JavaScript (tooling)
✓ Framework: Gin (HTTP router)
✓ Entry point: main.go → listens on :8080
✓ Database dependency: PostgreSQL (via lib/pq driver)
  - Migration files found in /migrations (12 migrations, latest: 012_add_indexes.sql)
✓ Cache dependency: Redis (via go-redis/v9)
✓ External API calls: AWS S3 (bucket operations), SendGrid (email)
✓ Message queue: None detected
✓ File storage: Local filesystem (/tmp/uploads) — likely needs externalization
✓ Auth: JWT validation middleware on /api/* routes
✗ No existing infrastructure-as-code found

Analysis ID: analysis-abc123
Proceed to infrastructure recommendation? (y/N)
```

The analysis ID is stored so that if you re-run the recommendation step later, Accio can reuse the analysis without re-reading code (unless the code has changed, which it detects via a content hash).

---

## Phase 3 — Infrastructure Recommendation

**What happens here:** The AI engine takes the analysis output and reasons about what cloud infrastructure would be appropriate, then presents you with ranked options.

### 3.1 — How the AI Engine Works

This is the core intelligence layer of Accio. It operates in two modes simultaneously:

**RAG (Retrieval-Augmented Generation):** Before generating a recommendation, the engine queries a curated knowledge base of infrastructure patterns. This knowledge base contains vetted, opinionated templates for common scenarios — things like "a Go web application with PostgreSQL and Redis on AWS" or "a Python ML serving endpoint with GPU instances." The retrieval step finds the most relevant patterns based on the analysis output, and these patterns are injected as context for the generation step. This prevents the AI from hallucinating infrastructure configurations that sound plausible but are architecturally wrong.

**MCP (Model Context Protocol):** The engine also has access to live data through MCP-connected tools. These tools can query things like current AWS pricing for specific instance types in specific regions, current availability of reserved instances, or known quotas and limits in your AWS account. This means recommendations are grounded in real-world constraints, not just theoretical best practices.

### 3.2 — Recommendation Output

The AI engine produces between 2 and 4 ranked options. Each option is a complete infrastructure blueprint, not a partial suggestion:

```
$ accio recommend

Based on your project analysis, here are the recommended infrastructure configurations:

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 Option 1 (Recommended) — "Balanced Production"                   
 Estimated monthly cost: $312 - $410                              
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Compute:   2x ECS Fargate tasks (1 vCPU, 2GB) behind an ALB
  Database:  RDS PostgreSQL db.t3.medium, Multi-AZ, 100GB gp3
  Cache:     ElastiCache Redis (cache.t3.micro), single node
  Storage:   S3 bucket (replaces local /tmp/uploads)
  Networking: VPC with public + private subnets, NAT Gateway
  Reasoning: Your app is stateless with moderate memory needs.
             Fargate removes instance management. Multi-AZ RDS
             protects against AZ failure for your primary data store.
             Redis single-node is sufficient until you have >1000 
             concurrent sessions.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 Option 2 — "Cost-Optimized"
 Estimated monthly cost: $145 - $195
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Compute:   1x ECS Fargate task (0.5 vCPU, 1GB) behind an ALB
  Database:  RDS PostgreSQL db.t3.micro, single-AZ, 50GB gp3
  Cache:     ElastiCache Redis (cache.t3.micro), single node
  Storage:   S3 bucket
  Networking: VPC with public + private subnets, no NAT Gateway
               (Fargate tasks in public subnet with security groups)
  Reasoning: Appropriate for dev/staging. Single-AZ and smaller
             instances cut cost significantly. No NAT Gateway saves
             ~$32/month + data transfer costs. Not recommended for
             production due to single-AZ database.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 Option 3 — "High Availability"
 Estimated monthly cost: $890 - $1,150
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Compute:   ECS Fargate with auto-scaling (2-10 tasks), ALB
  Database:  Aurora PostgreSQL Serverless v2, Multi-AZ
  Cache:     ElastiCache Redis cluster mode (3 shards, 1 replica each)
  Storage:   S3 bucket with replication to a second region
  Networking: VPC, 3 AZs, NAT Gateways (one per AZ), VPN Gateway
  Reasoning: Aurora Serverless v2 scales with load and is cost-
             effective at high utilization. Redis cluster mode provides
             both high availability and horizontal read scaling.
             Multi-region S3 protects against regional outages for
             your uploaded files.

? Select an option (1/2/3) or type 'custom' to start from scratch:
```

### 3.3 — Selection and Continuation

You select an option. That option becomes the base spec that enters the editing phase. If you type `custom`, you get a blank spec template with comments explaining each field, and the AI engine is available to answer questions as you fill it in.

---

## Phase 4 — Spec Editing & Interactive Refinement

**What happens here:** You take the generated spec and modify it to match your actual requirements. Accio guides you through missing details interactively.

### 4.1 — The Spec Format

The infrastructure spec is a YAML document. It is designed to be human-readable and declarative — you describe *what* you want, not *how* to build it. Accio translates the "what" into the "how" (Crossplane compositions) in a later phase.

```yaml
# accio-spec.yaml — generated by Accio for analysis-abc123
# Option 1: Balanced Production
# Generated: 2025-08-15T10:00:00Z

metadata:
  stack_name: my-project-prod
  environment: production
  team: backend-platform
  region: us-east-1

networking:
  vpc:
    cidr: ""                        # REQUIRED: e.g. 10.0.0.0/16
  subnets:
    public:
      pattern: ""                   # REQUIRED: e.g. 10.0.0.0/24, 10.0.1.0/24
      count: 2
    private:
      pattern: ""                   # REQUIRED: e.g. 10.0.128.0/24, 10.0.129.0/24
      count: 2
  nat_gateway:
    enabled: true
    count: 1                        # 1 = single NAT (cost-optimized), 2+ = per-AZ

compute:
  type: ecs_fargate
  task_count:
    min: 2
    max: 2                          # Set min != max to enable auto-scaling
  cpu: 1024                         # Units: 256, 512, 1024, 2048, 4096
  memory: 2048                      # MB
  load_balancer:
    type: alb
    scheme: internet-facing         # internet-facing | internal
  ssh_keypair:
    public_key: ""                  # REQUIRED if using EC2-backed compute

database:
  engine: postgresql
  instance_class: db.t3.medium
  storage_gb: 100
  storage_type: gp3
  multi_az: true
  backup_retention_days: 7
  maintenance_window: ""           # OPTIONAL: e.g. "sun:05:00-sun:05:30"

cache:
  engine: redis
  instance_class: cache.t3.micro
  cluster_mode: false
  replicas: 0                       # 0 = single node, 1+ = with replicas

storage:
  - type: s3
    name: uploads
    versioning: true
    encryption: AES256
    lifecycle_rules: []             # OPTIONAL: e.g. transition to Glacier after 90 days
```

### 4.2 — Interactive Prompting for Missing Fields

After you save and close the spec file (the CLI opens your `$EDITOR`), Accio validates it. Any required fields that are still empty trigger an interactive prompt:

```
$ accio spec validate

Validating spec...
✗ 3 required fields missing:

  1. networking.vpc.cidr
     What CIDR block should the VPC use? (e.g. 10.0.0.0/16)
     > 10.1.0.0/16

  2. networking.subnets.public.pattern
     VPC CIDR is 10.1.0.0/16. Suggest public subnet CIDRs (comma-separated):
     > 10.1.0.0/24, 10.1.1.0/24
     ✓ Both CIDRs are within VPC range and non-overlapping.

  3. networking.subnets.private.pattern
     Suggest private subnet CIDRs (comma-separated):
     > 10.1.128.0/24, 10.1.129.0/24
     ✓ Both CIDRs are within VPC range and non-overlapping with public subnets.

  All required fields populated.
  compute.ssh_keypair.public_key is empty but not required for Fargate — skipping.

✓ Spec is valid. Proceed to stack creation? (y/N)
```

The validation engine is not just syntax checking. It performs semantic validation too:

- Subnet CIDRs must be within the VPC CIDR and must not overlap each other.
- If `nat_gateway.count` is greater than 1, it must not exceed the number of public subnets or AZs.
- `task_count.max` must be greater than or equal to `task_count.min`.
- If `cluster_mode: true` for Redis, `replicas` must be at least 1.
- If `compute.type` is `ec2` instead of `ecs_fargate`, then `ssh_keypair.public_key` becomes required.

Warnings (non-blocking) are also surfaced:

```
⚠ networking.nat_gateway.count is 1 but you have 2 AZs.
  If a single NAT Gateway's AZ fails, outbound traffic from the other AZ's
  private subnets will be disrupted. Consider setting count: 2 for production.
```

### 4.3 — AI-Assisted Editing

At any point during spec editing, you can ask Accio questions inline:

```bash
accio spec ask "Should I use gp3 or io2 for my RDS storage given 12 migrations and moderate write load?"
```

The AI engine reasons about your specific situation (it still has the analysis context) and responds:

```
gp3 is the right choice here. Your migration history suggests a schema that has
stabilized — 12 migrations with the latest being index additions. That pattern
indicates read-heavy workloads more than write-heavy ones. gp3 gives you a
baseline of 3,000 IOPS at a lower cost than io2. If you later find that writes
are a bottleneck, gp3 allows you to provision additional IOPS incrementally
without switching storage types. io2 would only make sense if you're seeing
sustained >3,000 IOPS on writes, which is unlikely given the current profile.
```

---

## Phase 5 — Stack Validation & Persistence

**What happens here:** The validated spec is sent to the Accio API, checked against global constraints, persisted as a Stack object in PostgreSQL, and assigned a unique stack ID.

### 5.1 — API Submission

```bash
accio stack create --spec accio-spec.yaml
```

The CLI serializes the spec and sends it as a `POST /api/v1/stacks` request to the Accio API. The request body includes:

- The full spec YAML
- The analysis ID it was derived from
- The authenticated user's identity (from the token)
- The selected credential IDs for AWS and Git

### 5.2 — Server-Side Validation

The API performs a second round of validation that the CLI cannot do locally:

- **Cloud account constraint checks:** Using the registered AWS credentials, the API queries the AWS account for current quotas. For example, it checks whether the requested VPC CIDR conflicts with any existing VPC, whether the region has capacity for the requested instance types, and whether the account's service quotas allow the number of resources being requested.
- **Cost estimation:** The API calls AWS Pricing APIs (or a cached pricing index) to produce a more precise cost estimate than the one shown during recommendation. This estimate is returned to the user before the stack is actually created.
- **Naming conflict checks:** The stack name must be unique within the user's account. The API enforces this.

### 5.3 — Stack Persistence

If all validations pass, the API persists the stack:

```sql
INSERT INTO stacks (
  id,                   -- UUID, globally unique
  name,                 -- user-provided, e.g. "my-project-prod"
  owner_id,             -- user ID from the Keycloak token
  analysis_id,          -- links back to the code analysis
  spec,                 -- the full validated YAML, stored as JSONB
  status,               -- "pending_generation"
  aws_credential_id,    -- encrypted reference
  git_credential_id,    -- encrypted reference
  target_region,        -- e.g. "us-east-1"
  created_at,
  updated_at
) VALUES (...);
```

A stack version history table tracks every change:

```sql
INSERT INTO stack_versions (
  id,
  stack_id,
  version_number,       -- monotonically increasing integer
  spec_snapshot,        -- full spec at this version
  changed_by,           -- user ID or "system"
  change_reason,        -- e.g. "initial creation", "user updated cache config"
  created_at
) VALUES (...);
```

### 5.4 — Confirmation

```
$ accio stack create --spec accio-spec.yaml

Submitting stack to Accio API...
✓ Cloud account validation passed
✓ No VPC CIDR conflicts detected
✓ All requested instance types available in us-east-1
✓ Service quotas sufficient

  Stack created successfully.
  Stack ID:   stack-7f3a2b1c
  Stack Name: my-project-prod
  Version:    1
  Status:     pending_generation
  Est. Cost:  $318/month (on-demand)

  Next step: Accio will now generate Crossplane compositions and push to git.
  Run `accio stack status stack-7f3a2b1c` to monitor progress.
```

---

## Phase 6 — Template Generation & Git Commit

**What happens here:** The API's template generator takes the stack spec and produces Crossplane Composition and CompositeResource YAML files, then commits them to a git repository.

### 6.1 — Template Generation

The template generator is a deterministic, spec-driven engine. It is not AI-generated at this stage — it follows strict mapping rules from Accio's spec schema to Crossplane resource types. This is intentional: the infrastructure that actually gets created must be predictable and auditable.

For the spec from Phase 4, the generator produces:

```
compositions/
  stack-7f3a2b1c/
    networking/
      vpc.yaml                    # Crossplane Composition for VPC
      subnets.yaml                # Public + private subnets
      nat-gateway.yaml            # NAT Gateway + Elastic IP
      route-tables.yaml           # Route tables and associations
      security-groups.yaml        # Security groups for each component
    compute/
      ecs-cluster.yaml            # ECS cluster definition
      ecs-task-definition.yaml    # Task definition (CPU, memory, image)
      ecs-service.yaml            # ECS service (desired count, placement)
      alb.yaml                    # Application Load Balancer + listeners + target groups
    database/
      rds-instance.yaml           # RDS PostgreSQL instance
      rds-subnet-group.yaml       # Subnet group for RDS (private subnets)
      rds-parameter-group.yaml    # Parameter group (optional, for custom settings)
    cache/
      elasticache-redis.yaml      # ElastiCache Redis cluster
      elasticache-subnet-group.yaml
    storage/
      s3-bucket.yaml              # S3 bucket with versioning and encryption config
    meta/
      stack-composite.yaml        # Top-level CompositeResource that references all others
```

Each file is a valid Crossplane `Composition` or `XR` (CompositeResource) manifest. The generator ensures:

- All resources are **namespaced** under a Kubernetes namespace derived from the stack ID (e.g., `ns-stack-7f3a2b1c`). This is enforced at the Composition level.
- Resource names are deterministic and derived from the stack ID, so re-running generation produces identical output (idempotent).
- Dependency ordering is encoded via Crossplane's `readinessCheck` and `status.atProvider` field references — for example, the ECS service references the ALB's DNS name, which is only available after the ALB is created.
- All secrets (RDS passwords, etc.) are generated as Kubernetes Secrets by Crossplane and are never written to the YAML files in plain text.

### 6.2 — Git Commit

The API uses the registered Git PAT to:

1. **Check if the target repository exists.** The target repo is either specified by the user or auto-created with a name derived from the stack name (e.g., `accio-infra-my-project-prod`).
2. **If the repo does not exist, create it.** The API calls the Git hosting provider's API (e.g., GitHub's `POST /user/repos`) using the PAT.
3. **Clone the repo (or pull latest).**
4. **Write the generated files** into the repo under `compositions/stack-<id>/`.
5. **Commit with a structured message:**
   ```
   feat(stack-7f3a2b1c): initial infrastructure generation

   Stack: my-project-prod
   Version: 1
   Environment: production
   Region: us-east-1
   Resources: VPC, 2 subnets (pub+priv), NAT GW, ECS Fargate (2 tasks),
              ALB, RDS PostgreSQL (Multi-AZ), ElastiCache Redis, S3 bucket
   Analysis ID: analysis-abc123
   Generated by: Accio API v0.4.2
   ```
6. **Push to the repository's main branch.**
7. **Update the stack status in PostgreSQL** to `syncing`.

```
$ accio stack status stack-7f3a2b1c

Stack ID:    stack-7f3a2b1c
Status:      syncing
Version:     1
Git Repo:    https://github.com/acme/accio-infra-my-project-prod
Last Commit: feat(stack-7f3a2b1c): initial infrastructure generation
Pushed:      2025-08-15T10:05:12Z
```

---

## Phase 7 — GitOps Sync (ArgoCD + Crossplane)

**What happens here:** ArgoCD detects the new commit, syncs the manifests into the Kubernetes cluster, and Crossplane reconciles the desired state into actual cloud resources. This phase is largely automatic — Accio does not orchestrate it, but it monitors it.

### 7.1 — ArgoCD Application

An ArgoCD `Application` resource must already exist that points to the infra git repository. This is a one-time setup per repository. The Application is configured to:

- Watch the repository's main branch.
- Sync automatically on new commits (or on a manual trigger, depending on your policy).
- Apply manifests into the Kubernetes cluster.

When the commit from Phase 6 lands, ArgoCD detects the change and applies all YAML files under `compositions/stack-7f3a2b1c/` to the cluster.

### 7.2 — Crossplane Reconciliation

Crossplane controllers pick up the applied Compositions and CompositeResources and begin creating cloud resources. The order of creation is governed by the dependency graph encoded in the compositions:

```
stack-composite.yaml
  └── vpc.yaml                    → creates AWS VPC
        └── subnets.yaml          → waits for VPC ID, then creates subnets
              └── route-tables.yaml    → waits for subnet IDs
              └── nat-gateway.yaml     → waits for public subnet ID
        └── security-groups.yaml  → waits for VPC ID
  └── rds-instance.yaml           → waits for private subnet IDs + security group
  └── elasticache-redis.yaml      → waits for private subnet IDs + security group
  └── ecs-cluster.yaml            → no dependencies (can start immediately)
        └── ecs-task-definition.yaml  → waits for cluster
              └── ecs-service.yaml    → waits for task definition + ALB
  └── alb.yaml                    → waits for public subnet IDs + security group
  └── s3-bucket.yaml              → no dependencies (can start immediately)
```

Each Crossplane managed resource reports its status back. When all resources reach a `Ready` state, the stack is fully provisioned.

### 7.3 — Status Monitoring

Accio polls the Crossplane resource statuses (via the Kubernetes API) and surfaces them through the CLI:

```
$ accio stack status stack-7f3a2b1c --watch

Stack ID:    stack-7f3a2b1c
Status:      provisioning
Version:     1

Resource Status:
  ✓ VPC (vpc-0a1b2c3d)                          Ready
  ✓ Public Subnet 1 (subnet-aaaa)               Ready
  ✓ Public Subnet 2 (subnet-bbbb)               Ready
  ✓ Private Subnet 1 (subnet-cccc)              Ready
  ✓ Private Subnet 2 (subnet-dddd)              Ready
  ✓ NAT Gateway (nat-0x1y2z)                    Ready
  ✓ Route Tables                                Ready
  ✓ Security Groups                             Ready
  ✓ ALB (arn:aws:elasticloadbalancing:...)       Ready
  ✓ ECS Cluster                                 Ready
  ✓ ECS Task Definition                         Ready
  ✓ ECS Service (2/2 tasks running)             Ready
  ✓ RDS PostgreSQL (myproject-prod-db)          Ready
  ✓ ElastiCache Redis (myproject-prod-cache)    Ready
  ✓ S3 Bucket (myproject-prod-uploads)          Ready

  All 14 resources Ready. Stack is fully provisioned.
  Status updated to: running
```

If any resource fails, Accio surfaces the error immediately:

```
  ✗ RDS PostgreSQL                              Failed
    Error: DBInstanceAlreadyExists — an RDS instance with this name
    already exists in the account. This may indicate a naming conflict.
    Run `accio stack diagnose stack-7f3a2b1c` for resolution options.
```

---

## Phase 8 — Runtime Interaction & Observability

**What happens here:** Your stack is running. You interact with it, query it, and monitor it — all through the CLI, powered by AI and MCP tools that have live access to cloud APIs.

### 8.1 — Querying Your Infrastructure

The CLI exposes a conversational interface for asking questions about your running stack:

```bash
accio stack ask stack-7f3a2b1c "What is the current CPU utilization of my ECS tasks?"
```

Under the hood, the AI engine uses MCP-connected tools to:
1. Identify which ECS cluster and service belong to `stack-7f3a2b1c` (from the stack's metadata).
2. Query CloudWatch for the `CpuUtilization` metric on those tasks over the last 5 minutes (or a configurable window).
3. Format and present the results.

```
$ accio stack ask stack-7f3a2b1c "What is the current CPU utilization of my ECS tasks?"

Current ECS CPU Utilization (last 5 minutes):
  Task 1: 34% average, peak 41%
  Task 2: 29% average, peak 38%
  Service average: 31.5%
  Auto-scaling threshold: 70% (not triggered)

Your tasks are comfortably below the auto-scaling threshold. Current load
does not warrant scaling action.
```

### 8.2 — Cost Monitoring

```bash
accio stack cost stack-7f3a2b1c
accio stack cost stack-7f3a2b1c --period last-30-days
accio stack cost stack-7f3a2b1c --breakdown
```

```
$ accio stack cost stack-7f3a2b1c --breakdown

Cost Report: my-project-prod (stack-7f3a2b1c)
Period: 2025-07-16 to 2025-08-15

  ECS Fargate (compute)         $  42.18    (13.3%)
  ALB                           $  18.03    ( 5.7%)
  RDS PostgreSQL (Multi-AZ)     $ 139.50    (43.9%)
  ElastiCache Redis             $  13.68    ( 4.3%)
  NAT Gateway                   $  34.20    (10.8%)
  S3 (storage + requests)       $   2.11    ( 0.7%)
  Data Transfer                 $  12.44    ( 3.9%)
  Elastic IP                    $   0.00    ( 0.0%)  ← in use, no charge
  ─────────────────────────────────────
  Total                         $ 262.14    (100%)
  Projected next month          $ 271.00    (assuming same usage pattern)

💡 Optimization suggestion: RDS is your largest cost center. Your database
   shows average utilization of 22%. Consider downsizing to db.t3.small
   ($70/month savings) if performance remains acceptable after a 48-hour
   monitoring window. Run `accio stack ask` to simulate the impact.
```

### 8.3 — Infrastructure Modifications

If you need to change a running stack (e.g., scale up, add a new resource, change an instance type), the process mirrors the initial creation:

```bash
accio stack edit stack-7f3a2b1c
```

This opens the current spec in your `$EDITOR`. You make changes (e.g., change `instance_class: db.t3.medium` to `instance_class: db.t3.large`). On save:

1. The CLI re-validates the modified spec.
2. The CLI sends a `PUT /api/v1/stacks/stack-7f3a2b1c` request with the new spec.
3. The API persists a new version in `stack_versions`.
4. The template generator diffs the old and new specs and generates only the changed Crossplane manifests.
5. The API commits the changes to git with a message like:
   ```
   feat(stack-7f3a2b1c): scale RDS to db.t3.large

   Version: 2
   Changed: database.instance_class db.t3.medium → db.t3.large
   ```
6. ArgoCD syncs. Crossplane updates the RDS instance in place (or triggers a replacement, depending on the resource type and the nature of the change — Crossplane handles this).

### 8.4 — Tearing Down a Stack

```bash
accio stack destroy stack-7f3a2b1c
```

This is a destructive operation. The CLI requires explicit confirmation and prints exactly what will be deleted:

```
$ accio stack destroy stack-7f3a2b1c

⚠ WARNING: This will permanently destroy all cloud resources in stack-7f3a2b1c.

  Resources to be destroyed:
    - ECS Service (2 running tasks)
    - ECS Cluster
    - ALB
    - RDS PostgreSQL (myproject-prod-db) — DATA WILL BE LOST unless a snapshot exists
    - ElastiCache Redis (myproject-prod-cache) — DATA WILL BE LOST
    - NAT Gateway
    - S3 Bucket (myproject-prod-uploads) — BUCKET AND ALL OBJECTS WILL BE DELETED
    - Subnets, Route Tables, Security Groups, VPC
    - Elastic IP

  This action cannot be undone.
  Type the stack name to confirm: my-project-prod
```

On confirmation:

1. The API removes the composition files from the git repo and commits:
   ```
   feat(stack-7f3a2b1c): destroy stack my-project-prod
   ```
2. ArgoCD syncs the deletion.
3. Crossplane deletes all managed resources in reverse dependency order.
4. The stack status in PostgreSQL is updated to `destroyed`.
5. The stack record is retained (not deleted) for audit purposes.

---

## Data Model

### Stacks Table

| Column | Type | Description |
|---|---|---|
| `id` | `UUID` | Primary key, globally unique stack identifier |
| `name` | `VARCHAR(255)` | User-provided name, unique per owner |
| `owner_id` | `UUID` | Foreign key to the user (from Keycloak `sub` claim) |
| `analysis_id` | `VARCHAR(255)` | Links to the code analysis that seeded this stack |
| `spec` | `JSONB` | The current validated infrastructure spec |
| `status` | `ENUM` | `pending_generation`, `syncing`, `running`, `updating`, `destroying`, `destroyed`, `error` |
| `aws_credential_id` | `VARCHAR(255)` | Reference to the encrypted AWS credential |
| `git_credential_id` | `VARCHAR(255)` | Reference to the encrypted Git credential |
| `git_repo_url` | `TEXT` | The git repository URL where compositions are stored |
| `target_region` | `VARCHAR(20)` | AWS region |
| `created_at` | `TIMESTAMP` | Creation time |
| `updated_at` | `TIMESTAMP` | Last modification time |

### Stack Versions Table

| Column | Type | Description |
|---|---|---|
| `id` | `UUID` | Primary key |
| `stack_id` | `UUID` | Foreign key to `stacks` |
| `version_number` | `INTEGER` | Monotonically increasing, starts at 1 |
| `spec_snapshot` | `JSONB` | Full spec at this version (immutable snapshot) |
| `changed_by` | `UUID` | User who made the change, or NULL for system |
| `change_reason` | `TEXT` | Human-readable description of what changed |
| `git_commit_sha` | `VARCHAR(40)` | The git commit that corresponds to this version |
| `created_at` | `TIMESTAMP` | When this version was created |

### Credentials Table

| Column | Type | Description |
|---|---|---|
| `id` | `VARCHAR(255)` | User-facing credential ID (e.g., `aws-cred-a1b2`) |
| `owner_id` | `UUID` | Foreign key to the user |
| `type` | `ENUM` | `aws`, `git`, `gcp`, `azure` |
| `encrypted_payload` | `BYTEA` | The credential, encrypted with envelope encryption |
| `kms_key_id` | `VARCHAR(255)` | The KMS key used to encrypt the data key |
| `encrypted_data_key` | `BYTEA` | The data key, encrypted by the KMS key |
| `created_at` | `TIMESTAMP` | |
| `rotated_at` | `TIMESTAMP` | Last rotation time, NULL if never rotated |
| `revoked_at` | `TIMESTAMP` | Revocation time, NULL if active |

---

## Git Repository Structure

Each stack gets its own directory within a shared or dedicated git repository:

```
accio-infra-my-project-prod/
├── README.md                           # Auto-generated, describes the stack
├── compositions/
│   └── stack-7f3a2b1c/
│       ├── meta/
│       │   └── stack-composite.yaml   # Top-level CompositeResource
│       ├── networking/
│       │   ├── vpc.yaml
│       │   ├── subnets.yaml
│       │   ├── nat-gateway.yaml
│       │   ├── route-tables.yaml
│       │   └── security-groups.yaml
│       ├── compute/
│       │   ├── ecs-cluster.yaml
│       │   ├── ecs-task-definition.yaml
│       │   ├── ecs-service.yaml
│       │   └── alb.yaml
│       ├── database/
│       │   ├── rds-instance.yaml
│       │   ├── rds-subnet-group.yaml
│       │   └── rds-parameter-group.yaml
│       ├── cache/
│       │   ├── elasticache-redis.yaml
│       │   └── elasticache-subnet-group.yaml
│       └── storage/
│           └── s3-bucket.yaml
└── .argocd-app.yaml                    # ArgoCD Application manifest (optional)
```

If multiple stacks share an infra repo (e.g., `prod` and `staging` for the same project), they each occupy their own directory under `compositions/`.

---

## Crossplane Composition Design

Each Crossplane Composition follows a consistent structure. Here is an illustrative example for the RDS resource:

```yaml
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: stack-7f3a2b1c-rds
  labels:
    accio.dev/stack-id: stack-7f3a2b1c
    accio.dev/resource-type: database
    accio.dev/version: "1"
spec:
  compositeTypeRef:
    apiVersion: database.accio.dev/v1alpha1
    kind: XDatabase
  mode: Pipeline
  pipeline:
    - step: rds-subnet-group
      fn.crossplane.io/auto-ready: true
      patches:
        - type: FromCompositeFieldPath
          fromFieldPath: spec.forProvider.networking.privateSubnetIds
          toFieldPath: spec.forProvider.subnetIds
    - step: rds-instance
      fn.crossplane.io/auto-ready: true
      patches:
        - type: FromCompositeFieldPath
          fromFieldPath: spec.forProvider.database.instanceClass
          toFieldPath: spec.forProvider.instanceClass
        - type: FromCompositeFieldPath
          fromFieldPath: spec.forProvider.database.multiAZ
          toFieldPath: spec.forProvider.multiAZ
        - type: ToCompositeFieldPath
          fromFieldPath: status.atProvider.endpoint
          toFieldPath: status.atProvider.databaseEndpoint
  # All resources created by this composition are namespaced
  # under the stack's namespace via the CompositeResource's
  # spec.claimRef or metadata.namespace
```

Key design constraints enforced across all generated compositions:

- **Namespacing:** Every `ManagedResource` is created in the namespace `ns-stack-<stack-id>`. This is enforced by the Composition's patch logic, not left to the user.
- **Labeling:** All resources carry labels for stack ID, resource type, and generation version. This enables `accio stack status` to query resources by label.
- **Secret handling:** Passwords and connection strings are generated by Crossplane as `ConnectionDetails` and stored as Kubernetes Secrets. They are never written to the Composition YAML.
- **Drift detection:** Crossplane's default behavior is to reconcile drift — if someone manually changes a resource in the AWS console, Crossplane will revert it to the desired state on the next sync cycle. This is intentional and aligns with Accio's GitOps-first philosophy.

---

## Error Handling & Recovery

Accio is designed to fail explicitly and recover cleanly. Here is how errors are handled at each phase:

| Phase | Failure Scenario | Behavior |
|---|---|---|
| **Auth** | Token expired mid-session | CLI transparently refreshes using the refresh token. If refresh fails, prompts re-login. |
| **Auth** | Keycloak unreachable | CLI fails with a clear error and the Keycloak URL for troubleshooting. |
| **Analysis** | Unrecognized project structure | Accio reports what it found and what it could not identify. User can proceed with a manual/custom spec. |
| **Recommendation** | AI engine times out | CLI surfaces the timeout and offers to retry or skip to manual spec creation. |
| **Validation** | Subnet CIDR overlap | CLI blocks submission and explains the conflict with specific CIDR ranges. |
| **Validation** | AWS quota exceeded | API surfaces the specific quota, current usage, and the requested amount. |
| **Git** | PAT revoked or invalid | API returns an auth error. CLI prompts to re-register git credentials. |
| **Git** | Push conflict (branch diverged) | API pulls latest, rebases the generated files (non-conflicting by design), and retries. |
| **Crossplane** | Resource creation fails | Status surfaces the AWS error. User can run `accio stack diagnose` for guided resolution. |
| **Crossplane** | Partial failure (some resources created, others not) | Crossplane retries failed resources on subsequent reconciliation loops. No manual cleanup needed. |
| **Destroy** | Resource has deletion protection | Crossplane surfaces the error. User must either disable protection via `accio stack edit` or acknowledge the resource will be left orphaned. |

### `accio stack diagnose`

This command runs a set of automated checks against a stack that is in an `error` state:

```bash
accio stack diagnose stack-7f3a2b1c
```

It queries Crossplane for the most recent error events, cross-references them against known failure patterns (maintained in a knowledge base), and suggests corrective actions:

```
$ accio stack diagnose stack-7f3a2b1c

Diagnosing stack-7f3a2b1c...

Issue found: RDS instance creation failed
  AWS Error: The DB instance identifier 'myproject-prod-db' is already in use.
  
  Likely cause: A previous stack creation attempt partially succeeded and left
  this RDS instance behind, or another stack in the same AWS account uses
  the same instance identifier.

  Suggested actions:
    1. Check for orphaned RDS instances:
       accio stack orphans --account aws-cred-a1b2c3
    2. Rename this stack's RDS instance in the spec and retry:
       accio stack edit stack-7f3a2b1c  (change database.instance_name)
    3. If the existing instance belongs to a destroyed stack, it can be
       safely deleted via the AWS console or AWS CLI.
```

---

## Security Model

### Credential Scoping (Principle of Least Privilege)

The AWS IAM credentials you provide to Accio should be scoped to only the permissions Accio needs. The recommended IAM policy covers:

- `ec2:*` — VPC, subnets, NAT Gateways, security groups, Elastic IPs
- `ecs:*` — Clusters, task definitions, services
- `elasticloadbalancing:*` — ALBs, listeners, target groups
- `rds:*` — DB instances, subnet groups, parameter groups
- `elasticache:*` — Redis clusters, subnet groups
- `s3:*` — Bucket creation, object operations
- `cloudwatch:GetMetricData` — For cost and usage monitoring
- `pricing:GetProducts` — For cost estimation

This policy should be attached to a dedicated IAM user or role that is used exclusively for Accio. It should not be a root account or an admin policy.

### Credential Storage

Credentials are never stored on the developer's machine after submission (beyond the initial input). They are encrypted at rest in PostgreSQL using envelope encryption: a unique data key per credential, encrypted by an AWS KMS key. The KMS key is the only long-lived secret, and access to it is controlled via KMS key policies.

### Network Security

- All CLI-to-API communication is over TLS.
- The Accio API should be deployed behind a load balancer with TLS termination.
- The Crossplane provider credentials (used by Crossplane itself to create AWS resources) are stored as Kubernetes Secrets and are separate from the credentials Accio uses for validation and cost queries.

### Namespace Isolation

Every stack's resources are namespaced in Kubernetes. This means:
- A user cannot accidentally query or modify another stack's resources through Crossplane.
- RBAC policies on the Kubernetes cluster can enforce per-team or per-stack access boundaries.
- Resource cleanup on stack destruction is scoped — only resources in the stack's namespace are affected.

---

## CLI Command Reference

| Command | Description |
|---|---|
| `accio login` | Authenticate via Keycloak OIDC flow |
| `accio logout` | Revoke local token and clear session |
| `accio credentials add <type>` | Register a new cloud or git credential |
| `accio credentials list` | List all registered credentials |
| `accio credentials rotate <id>` | Rotate a credential (re-encrypt with new key) |
| `accio credentials revoke <id>` | Revoke a credential permanently |
| `accio analyze` | Analyze the current project directory |
| `accio recommend` | Generate infrastructure recommendations from analysis |
| `accio spec validate` | Validate a spec file and prompt for missing fields |
| `accio spec ask "<question>"` | Ask the AI engine a question about your spec |
| `accio stack create --spec <file>` | Submit a spec and create a stack |
| `accio stack list` | List all stacks owned by the current user |
| `accio stack status <id>` | Show current status and resource states |
| `accio stack status <id> --watch` | Continuously poll and display status updates |
| `accio stack edit <id>` | Open the current spec for editing |
| `accio stack cost <id>` | Show cost report for a stack |
| `accio stack cost <id> --breakdown` | Show per-resource cost breakdown |
| `accio stack ask <id> "<question>"` | Ask the AI engine about a running stack |
| `accio stack diagnose <id>` | Run automated diagnostics on a failed stack |
| `accio stack orphans --account <id>` | Find cloud resources not tracked by any stack |
| `accio stack destroy <id>` | Destroy all resources in a stack (requires confirmation) |