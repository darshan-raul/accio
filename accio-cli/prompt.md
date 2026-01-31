Nice, this is a solid direction 👍
Below is a **single, end-to-end prompt** you can directly give to **Codex / Antigravity / any code-gen IDE agent** to build **accio**, your Bubble Tea–based TUI CLI.

I’ve written it as a **“build this repo from scratch” system prompt**, with **clear architecture, constraints, UX expectations, and deliverables**—so the agent doesn’t hallucinate scope or drift away from your platform vision.

---

## 🔮 MASTER PROMPT: Build *accio* CLI

> **Role**: You are a senior Go platform engineer and TUI designer.
> **Goal**: Build a production-grade TUI CLI called **accio**, inspired by CloudCode, using **Golang + Bubble Tea**, which acts as the primary developer entry point to a multi-cloud infrastructure platform built on **Crossplane, Argo CD, and Kubernetes**.

---

### 🧠 High-Level Context

accio is **NOT** a thin wrapper over kubectl.
It is a **platform CLI** that allows developers and DevOps engineers to:

* Define cloud infrastructure **intentionally**, not imperatively
* Work across **AWS, GCP, Azure**
* Interact with a **platform API** (not cloud APIs directly)
* Generate, validate, preview, and submit **infrastructure intents**
* Observe status of Crossplane-managed resources

accio is the **starting point** of the platform UX.

---

## 🧱 Core Principles

* **TUI-first** (Bubble Tea)
* **Opinionated UX** (guided flows > free-form commands)
* **Declarative mental model**
* **Cloud-agnostic vocabulary**
* **API-driven** (no direct SDK calls to AWS/GCP/Azure)
* **Extensible** (plugins, future MCP/AI hooks)

---

## 🧰 Tech Stack (Mandatory)

* Language: **Go**
* TUI Framework:

  * `bubbletea`
  * `bubbles` (list, table, spinner, viewport, textinput)
  * `lipgloss`
* HTTP client for API interaction
* YAML/JSON handling
* Cobra **only for entrypoint**, not UX
* Works on **Linux / macOS**
* No GUI frameworks

---

## 🏗️ Architecture Requirements

### 1. Binary

```
accio
```

### 2. Layered Design

```
cmd/
internal/
  tui/
    screens/
    components/
    styles/
  api/
  config/
  models/
  workflows/
  state/
pkg/
```

* `tui/screens` → full-screen views (dashboard, create infra, status)
* `tui/components` → reusable widgets
* `api` → typed client for platform API
* `models` → InfraIntent, Stack, Environment, ResourceStatus
* `workflows` → multi-step flows (wizard logic)
* `state` → Bubble Tea model + update logic

---

## 🖥️ TUI UX REQUIREMENTS

### Initial Screen (Landing Dashboard)

* Platform status
* Current context:

  * Org
  * Project
  * Environment (dev/stage/prod)
* Quick actions:

  * `Create Infrastructure`
  * `View Stacks`
  * `Preview Changes`
  * `Deploy`
  * `Observe Status`
  * `Settings`

Keyboard-driven:

* `↑ ↓` navigate
* `Enter` select
* `Esc` back
* `q` quit

---

### Infrastructure Creation Flow (Wizard)

**Step-based TUI wizard**

1. Choose Cloud

   * AWS / GCP / Azure
2. Choose Environment
3. Choose Capability (abstracted)

   * Compute
   * Kubernetes Cluster
   * Database
   * Object Storage
   * Networking
4. Collect Inputs

   * region
   * size / tier
   * scaling intent
5. Generate **Infra Intent**

   * Show YAML preview
6. Validate
7. Submit to API

---

### Example Intent (Internal Model)

```yaml
apiVersion: accio.io/v1alpha1
kind: InfraIntent
metadata:
  name: payments-backend
spec:
  cloud: aws
  environment: prod
  capabilities:
    - type: kubernetes
      size: medium
      region: ap-south-1
```

---

## 🌐 API Interaction (Abstract)

accio talks to a **Platform API**, not Kubernetes directly.

Endpoints (assume):

* `POST /intents`
* `GET /intents`
* `GET /intents/{id}/status`
* `POST /preview`
* `POST /deploy`

Implement:

* Typed client
* Error handling
* Loading spinners
* Retry/backoff

Mock server acceptable for now.

---

## 🔍 Observability Screen

* List of stacks
* Sync status (Crossplane / Argo CD)
* Health:

  * Ready
  * Reconciling
  * Failed
* Resource drill-down

---

## ⚙️ Config & Context

* Config file:

  ```
  ~/.accio/config.yaml
  ```
* Stores:

  * API endpoint
  * Auth token
  * Default org/project/env

Support:

* Context switching inside TUI
* Future: `accio context use prod`

---

## 🧪 Quality Bar

* Idiomatic Go
* Clean Bubble Tea update/view separation
* No god-models
* Readable, testable code
* Comments where UX logic is non-obvious

---

## 📦 Deliverables

1. Fully working Go module
2. TUI application with:

   * Navigation
   * Wizard flow
   * API interaction
3. README covering:

   * Architecture
   * UX philosophy
   * How it maps to Crossplane / Argo CD
4. Mock API server for local testing

---

## 🚀 Stretch (If Time Permits)

* Diff view for previews
* Plugin system
* AI/MCP hook (intent suggestion)
* Keyboard shortcut help modal

---

### Final Instruction

> Build this as if it will be the **primary interface** for a multi-tenant cloud platform.
> Favor **clarity, flow, and correctness** over cleverness.

---

If you want, next we can:

* Break this into **milestones**
* Design **exact Bubble Tea models per screen**
* Or map **accio intents → Crossplane XRDs** directly 👀
