<p align="center">
  <img src="https://raw.githubusercontent.com/AshvinBambhaniya/nexus-tasks/main/frontend/public/logo.svg" alt="Nexus Tasks Logo" width="128" />
</p>

<h1 align="center">Nexus Tasks Helm Charts</h1>

<p align="center">
  <strong>Official Helm Chart Repository for the Developer-Focused Project Management Platform</strong>
</p>

<p align="center">
  <a href="https://artifacthub.io/packages/helm/nexus-tasks/nexus-tasks"><img src="https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/nexus-tasks" alt="Artifact Hub" /></a>
  <a href="https://github.com/AshvinBambhaniya/nexus-tasks/releases"><img src="https://img.shields.io/badge/Helm%20Chart-v2.0.0-success.svg?logo=helm" alt="Helm Chart" /></a>
  <a href="https://github.com/AshvinBambhaniya/nexus-tasks/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-AGPL_v3-blue.svg" alt="License" /></a>
</p>

---

Welcome to the **Nexus Tasks** Helm repository. This site hosts the automated distribution of our Kubernetes deployment manifests. Nexus Tasks is designed to bridge the gap between personal productivity and team collaboration with a high-performance, developer-centric UI.

## Usage

### 1. Add this repository
To access the Nexus Tasks charts, add this repository to your local Helm configuration:

```bash
helm repo add nexus-tasks https://AshvinBambhaniya.github.io/nexus-tasks/
helm repo update
```

### 2. Install the Chart
Deploy the full stack (Frontend, Backend, PostgreSQL, and Redis) using the default configuration:

```bash
# Install the latest stable release
helm install my-nexus nexus-tasks/nexus-tasks
```

### 3. Verify Deployment
```bash
kubectl get pods -l app.kubernetes.io/instance=my-nexus
```

---

## Configuration

We provide a comprehensive **JSON Schema** for robust value validation and a rich configuration experience. 

### Inspect Default Values
You can explore all available parameters directly from your CLI:

```bash
helm show values nexus-tasks/nexus-tasks
```

### Advanced Setup
For detailed information on configuring SMTP, external databases, Ingress TLS, and persistent storage, please refer to the [official Chart Documentation](https://github.com/AshvinBambhaniya/nexus-tasks/tree/main/deploy/helm/nexus-tasks).

---

## Project Links
- **[Source Code](https://github.com/AshvinBambhaniya/nexus-tasks)**: Issues, Feature Requests, and Contributions.
- **[Artifact Hub](https://artifacthub.io/packages/helm/nexus-tasks/nexus-tasks)**: Interactive UI for parameter configuration.
- **[Security Policy](https://github.com/AshvinBambhaniya/nexus-tasks/security/policy)**: Report vulnerabilities.

---

## Repository Maintenance
This branch (`gh-pages`) is a **production distribution endpoint** managed automatically. 

- **Maintenance**: Changes to chart templates must be submitted via PRs to the `main` branch.
- **Automation**: Managed by [helm/chart-releaser-action](https://github.com/helm/chart-releaser-action).

> **Warning**: Do not manually edit the `index.yaml` file in this branch as it will be overwritten by the next automated release.

---

<p align="center">
  Distributed under the <strong>GNU Affero General Public License v3.0</strong>
</p>
