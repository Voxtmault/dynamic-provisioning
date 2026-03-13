# Dynamic Multi-Tenant Messaging Board (PoC)

A Proof-of-Concept (PoC) demonstrating a highly scalable, dynamic white-label SaaS architecture. This project automates the provisioning of isolated tenant environments, spinning up standalone backend and frontend containers on-the-fly upon user registration, complete with automated subdomain routing and SSL generation.

This repository serves as the foundational architecture for a project i have in a company i currently work in (and mostly shower thoughts).

## Key Features

* **Dynamic Container Provisioning:** Uses the Docker Engine API to programmatically spin up standalone backend and frontend containers for each registered tenant without manual intervention.
* **White-Label Frontend:** A single, generic compiled Flutter Web SPA (Single Page Application) that dynamically fetches its tenant-specific profile (branding, logos, names) at runtime.
* **Automated Subdomain Routing:** Seamless integration with Traefik to dynamically route traffic to newly spun-up containers (e.g., `tenant1.domain.com`) without proxy restarts.
* **Automated Wildcard SSL:** Utilizes Let's Encrypt with the DNS-01 challenge to automatically secure all tenant subdomains.
* **Hybrid Data Isolation:** * **PostgreSQL:** Dedicated logical databases per tenant managed via dynamic GORM migrations.
    * **Garage (S3):** Dedicated object storage buckets per tenant for picture and media uploads.
    * **Redis:** Centralized caching with strict tenant-ID key prefixing.
    * **OpenBao:** Centralized secrets management with isolated, path-based access control per tenant.
* **Zero-Downtime Schema Updates:** Orchestrator-driven batch rollout mechanism for applying database migrations securely across all tenant containers.

## Technology Stack

* **Frontend:** Flutter Web, Nginx (Multi-stage Docker build)
* **Backend:** Go, GORM
* **Databases & Caching:** PostgreSQL, Redis
* **Object Storage:** Garage (S3-compatible)
* **Secrets Management:** OpenBao
* **Infrastructure & Routing:** Docker, Traefik Reverse Proxy

## Repository Structure

This project utilizes a monorepo structure to ensure seamless collaboration across UI/UX, backend logic, and DevOps orchestration.

```text
.
├── /tenant-backend             # Go application, REST APIs, and GORM database models for tenants
├── /tenant-frontend            # Flutter Web source code and Nginx multi-stage Dockerfile for tenants
├── /admin-backend              # Go application, REST APIs, and GORM database models for admin
├── /admin-frontend             # Flutter Web source code and Nginx multi-stage Dockerfile for admin
└── /infrastructure             # Core static stack (docker-compose.yml for Traefik, Postgres, Redis, etc.)
```

## Architecture Overview

When a new client registers via the admin interface:

1. The Orchestrator creates a new logical PostgreSQL database and a dedicated Garage S3 bucket.

2. Unique credentials are generated and securely stored in OpenBao under a tenant-specific path.

3. The Orchestrator calls the Docker Engine API to spin up two new standalone containers (Frontend & Backend).

4. Traefik detects the new containers via Docker labels and dynamically creates routing rules for the new subdomain.

5. The compiled Flutter Frontend boots in the user's browser, requests the tenant profile from the backend, and renders the customized messaging board.

## Getting Started
1. Navigate to the /infrastructure directory and start the core services:

```bash
docker compose up -d
```

2. Build the generic frontend and backend Docker images.

3. Run the Orchestrator service to initialize the environment and listen for tenant registrations.