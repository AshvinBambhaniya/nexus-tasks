# Deployment & Publishing Guide

This guide covers the deployment of Nexus Tasks v2.0.0 using Docker and Kubernetes (Helm).

## 1. Local Development (Docker Compose)

The fastest way to get started is with Docker Compose.

```bash
docker-compose up --build
```
- **Frontend**: http://localhost:3000
- **API**: http://localhost:8000

---

## 2. Kubernetes Deployment (Helm)

### Prerequisites
- Kubernetes Cluster (Docker Desktop, Minikube, or Cloud)
- Helm 3.8.0+

### Step 1: Images
By default, this chart uses official images from **GHCR** (`ghcr.io/ashvinbambhaniya/nexus-tasks-*`). If you wish to use your own local images, you must build them and update the `values.yaml`:

```bash
# Optional: Build local images for development
docker build -t nexus-tasks-backend:2.0.0 ./backend
docker build -t nexus-tasks-frontend:2.0.0 ./frontend
```

### Step 2: Initialize Chart
```bash
cd deploy/helm/nexus-tasks
helm dependency update
```

### Step 3: Install
For a **minimal** installation:
```bash
helm install nexus .
```

For a **full collaborative** installation (with Workers and SMTP):
```bash
helm install nexus . \
  --set worker.enabled=true \
  --set smtp.enabled=true \
  --set smtp.host="smtp.your-provider.com" \
  --set smtp.username="your-user" \
  --set smtp.password="your-password"
```

---

## 3. Publishing to Artifact Hub (Automated)

We use **GitHub Actions** and **GitHub Pages** to automatically package and publish our Helm charts.

### How it works:
1. **Trigger**: Any push to the `main` branch that modifies files in `deploy/helm/nexus-tasks` triggers the `Release Charts` workflow.
2. **Action**: The `helm/chart-releaser-action` automatically:
   - Packages the chart.
   - Creates a GitHub Release for the new version.
   - Updates the `index.yaml` on the `gh-pages` branch.
3. **Distribution**: The chart is served via GitHub Pages at `https://AshvinBambhaniya.github.io/nexus-tasks/`.

### Manual Registration (One-time):
To list the chart on Artifact Hub:
1. Ensure the repository has a `deploy/helm/artifacthub-repo.yml` file.
2. Log in to [Artifact Hub](https://artifacthub.io/).
3. Add a new repository:
   - **Type**: Helm charts
   - **Name**: `nexus-tasks`
   - **URL**: `https://AshvinBambhaniya.github.io/nexus-tasks/`

Artifact Hub will automatically sync whenever the GitHub Action publishes a new version.

---

## 4. Production Best Practices
- **Persistence**: Ensure `postgresql.primary.persistence.enabled` is `true`.
- **Security**: Always change `config.jwtSecret` in production.
- **TLS**: Enable Ingress TLS in `values.yaml` for production-grade security.
