---
name: openclaw-helm
description: Helm-based deployment and management of OpenClaw in Kubernetes environments. Use when: (1) Deploying OpenClaw to Kubernetes clusters, (2) Managing OpenClaw Helm chart installations, (3) Configuring OpenClaw for production use, (4) Setting up backup and restore procedures, (5) Upgrading OpenClaw deployments, or (6) Troubleshooting OpenClaw in Kubernetes.
---

# OpenClaw Helm Deployment

## Overview

This skill provides comprehensive guidance for deploying and managing OpenClaw AI Agent in Kubernetes using Helm charts. It covers installation, configuration, backup, upgrade, and troubleshooting procedures for production-grade OpenClaw deployments.

## Quick Start

### Prerequisites
- Kubernetes cluster (k3s, minikube, EKS, etc.)
- Helm v3.8+ installed
- kubectl configured
- Sudo access (for k3s)

### Basic Installation

```bash
# Add repository and install
helm repo add openclaw https://thepagent.github.io/openclaw-helm
helm install openclaw openclaw/openclaw-helm -n openclaw --create-namespace
```

### Minimal Values Configuration

Create `values.yaml` for resource limits and basic configuration:

```yaml
replicaCount: 1
resources:
  requests:
    memory: "512Mi"
    cpu: "250m"
  limits:
    memory: "1Gi"
    cpu: "500m"
envFrom:
  - secretRef:
      name: openclaw-api-keys
skills:
  - weather
  - healthcheck
```

Install with custom values:

```bash
helm install openclaw openclaw/openclaw-helm -n openclaw -f values.yaml
```

## Deployment Workflows

### 1. First-Time Setup

#### Step 1: Install the Chart

```bash
# For k3s users
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
sudo -E helm install openclaw oci://ghcr.io/thepagent/openclaw-helm -n openclaw --create-namespace
```

#### Step 2: Run Onboarding Wizard

```bash
POD=$(kubectl get pod -n openclaw -l app.kubernetes.io/name=openclaw-helm -o jsonpath='{.items[0].metadata.name}')
kubectl exec -it -n openclaw $POD -c main -- openclaw onboard
```

The wizard guides you through:
- Security acknowledgment
- AI provider selection (OpenAI, Anthropic, etc.)
- Authentication setup
- Channel configuration (Telegram, WhatsApp, etc.)
- Skills installation

#### Step 3: Verify Deployment

```bash
# Run connectivity test
helm test openclaw -n openclaw

# Check gateway status
kubectl exec -n openclaw $POD -- openclaw models status
```

### 2. Channel Configuration

#### Telegram Setup

```bash
# Enter configuration mode
kubectl exec -it -n openclaw $POD -c main -- openclaw configure --section channels

# Pair Telegram (get pairing code from bot)
kubectl exec -it -n openclaw $POD -c main -- openclaw pairing approve telegram <pairing-code>
```

#### API Key via Secret (Alternative)

```bash
# Create secret with API keys
kubectl create secret generic openclaw-api-keys \
  --from-literal=OPENAI_API_KEY=sk-... \
  --from-literal=ANTHROPIC_API_KEY=sk-... \
  --namespace openclaw
```

Update values.yaml:

```yaml
envFrom:
  - secretRef:
      name: openclaw-api-keys
```

Upgrade:

```bash
helm upgrade openclaw oci://ghcr.io/thepagent/openclaw-helm -n openclaw -f values.yaml
```

### 3. Backup and Restore

#### Create Backup

```bash
# Create backup directory
BACKUP_DIR=~/openclaw-backup-$(date +%Y%m%d-%H%M%S)
mkdir -p $BACKUP_DIR

# Get pod name and copy data
POD=$(kubectl get pod -n openclaw -l app.kubernetes.io/name=openclaw-helm -o jsonpath='{.items[0].metadata.name}')
kubectl cp openclaw/$POD:/home/node/.openclaw $BACKUP_DIR

# Compress backup
tar -czf $BACKUP_DIR.tar.gz $BACKUP_DIR
```

#### Restore from Backup

```bash
# Decompress backup
tar -xzf backup-file.tar.gz

# Get current pod
POD=$(kubectl get pod -n openclaw -l app.kubernetes.io/name=openclaw-helm -o jsonpath='{.items[0].metadata.name}')

# Restore configuration
kubectl cp ./backup-dir/. openclaw/$POD:/home/node/.openclaw/

# Restart deployment
kubectl rollout restart deployment/openclaw-openclaw-helm -n openclaw
```

### 4. Upgrade Procedures

#### Standard Upgrade

```bash
# For k3s users
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

# Backup first (recommended)
# ... backup commands ...

# Upgrade chart
sudo -E helm upgrade openclaw oci://ghcr.io/thepagent/openclaw-helm -n openclaw
```

#### Rollback Procedure

```bash
# List revisions
helm history openclaw -n openclaw

# Rollback to specific revision
helm rollback openclaw <revision-number> -n openclaw
```

## Troubleshooting Guide

### Common Issues

#### 1. Pod Fails to Start

```bash
# Check pod status
kubectl get pods -n openclaw

# View logs
kubectl logs -n openclaw <pod-name> -c main

# Describe pod for events
kubectl describe pod -n openclaw <pod-name>
```

**Common causes:**
- Insufficient resources (increase memory/cpu in values.yaml)
- Missing API keys (create secret or run onboarding)
- PVC issues (check storage class)

#### 2. Gateway Not Responding

```bash
# Check service
kubectl get svc -n openclaw

# Test connectivity
kubectl exec -n openclaw <pod-name> -- curl -s http://localhost:8000

# Check gateway logs
kubectl logs -n openclaw <pod-name> -c main --tail=50
```

#### 3. Skills Not Loading

```bash
# Check init container logs
kubectl logs -n openclaw <pod-name> -c init-skills

# Verify skills directory
kubectl exec -n openclaw <pod-name> -- ls -la /home/node/.openclaw/workspace/skills
```

### k3s-Specific Issues

#### Permission Issues

```bash
# Set KUBECONFIG for sudo operations
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
sudo -E helm install openclaw oci://ghcr.io/thepagent/openclaw-helm -n openclaw
```

#### Storage Class Issues

Create custom values.yaml for k3s:

```yaml
persistence:
  storageClass: "local-path"
  accessModes:
    - ReadWriteOnce
  size: 5Gi
```

## Advanced Configuration

### Custom Values.yaml Examples

#### Production Configuration

```yaml
replicaCount: 2
resources:
  requests:
    memory: "1Gi"
    cpu: "500m"
  limits:
    memory: "2Gi"
    cpu: "1000m"

persistence:
  storageClass: "fast-ssd"
  size: "20Gi"

skills:
  - weather
  - healthcheck
  - coding-agent
  - openai-whisper-api

env:
  - name: OPENCLAW_LOG_LEVEL
    value: "info"
  - name: OPENCLAW_GATEWAY_PORT
    value: "8000"
```

#### Development Configuration

```yaml
replicaCount: 1
resources:
  requests:
    memory: "256Mi"
    cpu: "100m"
  limits:
    memory: "512Mi"
    cpu: "250m"

skills: []

env:
  - name: OPENCLAW_LOG_LEVEL
    value: "debug"
```

### Monitoring and Logging

```bash
# View real-time logs
kubectl logs -n openclaw -f <pod-name> -c main

# Check resource usage
kubectl top pod -n openclaw

# Port forward for local testing
kubectl port-forward -n openclaw svc/openclaw-openclaw-helm 8080:8000
```

## Design Philosophy

This Helm chart follows these principles:

1. **Minimal Footprint**: Optimized for small cloud instances (2 vCPU, 2GB RAM)
2. **Decoupled Browser**: No Chromium integration; use Vercel agent-browser for browser automation
3. **Security by Design**: Gateway binds to loopback (127.0.0.1) within Kubernetes
4. **Persistent Configuration**: All settings stored in PVC; upgrades preserve configuration

## Resources

### Scripts

This skill includes utility scripts in the `scripts/` directory:

- `backup_openclaw.sh` - Automated backup script
- `restore_openclaw.sh` - Restoration script
- `health_check.sh` - Health check and validation

### References

- `CHART_REFERENCE.md` - Detailed chart configuration options
- `TROUBLESHOOTING.md` - Extended troubleshooting guide
- `SECURITY.md` - Security considerations and hardening

### Templates

- `values-prod.yaml` - Production-ready values template
- `values-dev.yaml` - Development values template
- `backup-script-template.sh` - Customizable backup script template

---

**Note**: Always backup your OpenClaw configuration before upgrades or major changes. The Helm chart preserves existing `openclaw.json` configuration but backups provide safety against unexpected issues.