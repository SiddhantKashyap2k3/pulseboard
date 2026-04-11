## PulseBoard — Real-time SaaS Analytics Dashboard

The full breakdown of what was build.

---

### What it does

A multi-tenant SaaS platform where users ingest events (page views, API calls, custom events) via an SDK or HTTP endpoint, and see them visualized in real-time on a dashboard — with charts, funnels, retention tables, and live counters. Think a lightweight Mixpanel/Datadog hybrid.

---

### Tech stack breakdown

**Frontend** — React + TypeScript with Vite. Recharts for graphs, WebSocket for live data, TailwindCSS for styling. Served via Nginx inside a Docker container.

**Backend (Go)** — Two services: a REST API (Gin framework) handling auth, CRUD, and event ingestion; and a WebSocket server that pushes live metrics to connected clients. Go is a standout choice here — it signals maturity and real systems knowledge.

**Data layer** — PostgreSQL as the primary store (tenants, dashboards, aggregated stats), Redis for caching + pub/sub between the REST and WebSocket layers, and Kafka for durable event streaming from the ingest endpoint to the worker pods.

**Kubernetes on EKS** — Each service runs as a separate Deployment with HPA (Horizontal Pod Autoscaler). Prometheus + Grafana for internal metrics.

**GitHub Actions CI/CD** — Four stages: lint/test → Docker build + push to ECR → Helm deploy to staging → manual-approval prod deploy. Every PR triggers the first three automatically.

**Terraform** — All AWS infrastructure (VPC, EKS cluster, RDS, ElastiCache, ECR, S3, IAM with IRSA) written as code in a `terraform/` directory. Remote state in S3 with DynamoDB locking.

---

### Repository structure

```
pulseboard/
├── backend/
│   ├── api/           # Go REST API (Gin)
│   ├── ws/            # Go WebSocket server
│   ├── worker/        # Kafka consumer
│   └── Dockerfile
├── frontend/
│   ├── src/
│   └── Dockerfile
├── terraform/
│   ├── modules/       # vpc, eks, rds, cache
│   └── environments/  # staging/, prod/
├── k8s/
│   └── helm/          # Helm chart for all services
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── deploy.yml
└── docker-compose.yml # Local dev stack
```

---

### Tech Dominance

| Signal | What it shows |
|---|---|
| Go backend | Production-grade language, not just JS everywhere |
| Kafka + WebSockets | Real-time distributed systems experience |
| Multi-tenant auth | Real SaaS thinking — not just a CRUD app |
| Terraform modules | IaC in reusable, production patterns |
| Helm + HPA | Kubernetes beyond "kubectl apply" |
| CI/CD with staging gate | Full delivery lifecycle ownership |
| Prometheus + Grafana | Observability as a first-class concern |

---

### Phases

**Phase 1 (Week 1–2):** Local dev — Docker Compose with Go API + PostgreSQL + Redis. Basic auth + dashboard CRUD.

**Phase 2 (Week 2–3):** Add Kafka ingest pipeline, WebSocket live push, React frontend with charts.

**Phase 3 (Week 3–4):** Write Terraform modules, provision EKS, write Helm chart, deploy manually.

**Phase 4 (Week 4–5):** Wire up GitHub Actions CI/CD. Add Prometheus metrics to Go services. Write a README with architecture diagram.

---