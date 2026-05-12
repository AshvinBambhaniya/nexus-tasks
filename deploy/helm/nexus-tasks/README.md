# Nexus Tasks Helm Chart

[Nexus Tasks](https://github.com/AshvinBambhaniya/nexus-tasks) is a developer-focused project management platform designed for speed, simplicity, and Markdown support. This Helm chart bootstraps a Nexus Tasks deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.22+
- Helm 3.8.0+
- PV provisioner support in the underlying infrastructure

## Installation

### Add Helm Repository

```bash
helm repo add nexus-tasks https://AshvinBambhaniya.github.io/nexus-tasks/
helm repo update
```

### Install Chart

```bash
helm install my-nexus-tasks nexus-tasks/nexus-tasks --version 2.0.0
```

## Configuration

The following table lists the configurable parameters of the Nexus Tasks chart and their default values.

### Global & App Configuration
| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas for deployments | `1` |
| `serviceAccount.create` | Create a service account | `true` |
| `config.appEnv` | Application environment | `local` |
| `config.debug` | Enable debug mode | `true` |
| `config.allowedOrigins` | CORS allowed origins | `http://localhost:3000` |
| `resources` | Resource limits and requests | (Sane defaults provided) |

### Components
| Parameter | Description | Default |
|-----------|-------------|---------|
| `backend.image.repository` | Backend image repository | `ghcr.io/ashvinbambhaniya/nexus-tasks-backend` |
| `backend.image.tag` | Backend image tag (defaults to `.Chart.AppVersion`) | `""` |
| `frontend.image.repository` | Frontend image repository | `ghcr.io/ashvinbambhaniya/nexus-tasks-frontend` |
| `frontend.image.tag` | Frontend image tag (defaults to `.Chart.AppVersion`) | `""` |
| `worker.enabled` | Enable background task processor | `false` |
| `dlqWorker.enabled` | Enable Dead Letter Queue worker | `false` |
| `postgresql.enabled` | Deploy PostgreSQL subchart | `true` |
| `redis.enabled` | Deploy Redis subchart | `true` |

### SMTP Configuration
| Parameter | Description | Default |
|-----------|-------------|---------|
| `smtp.enabled` | Enable SMTP environment variables | `false` |
| `smtp.host` | SMTP server host | `""` |
| `smtp.port` | SMTP server port | `587` |
| `smtp.fromEmail` | Sender email address | `noreply@example.com` |

Refer to [values.yaml](values.yaml) for the full list of configuration options.

## Modular Architecture

This chart is designed for a **Lean-by-Default** deployment:
- **Core Only**: Installs API, Frontend, and Databases.
- **Collaborative**: Set `worker.enabled: true` and `smtp.enabled: true` to process email invitations.
- **High-Availability**: Enable `dlqWorker` for enterprise-grade error handling.

## Persistence

The chart uses Persistent Volume Claims (PVCs) for PostgreSQL to ensure data is preserved across pod restarts. Redis persistence is disabled by default to save resources but can be enabled in `values.yaml`.

## Ingress

By default, an Ingress is created. The chart automatically configures the frontend and backend routes. If TLS is enabled, the backend API URLs will automatically switch to `https`.
