# OpenClaw Helm Troubleshooting Guide

## Quick Diagnosis Flowchart

```
Start
  ↓
Can you connect to Kubernetes? → No → Check kubectl config
  ↓ Yes
Does namespace exist? → No → Create namespace
  ↓ Yes
Are there any pods? → No → Check deployment
  ↓ Yes
Is pod running? → No → Check pod events
  ↓ Yes
Is pod ready? → No → Check container logs
  ↓ Yes
Is gateway responding? → No → Check gateway logs
  ↓ Yes
✅ OpenClaw is healthy
```

## Common Issues and Solutions

### 1. Installation Issues

#### Issue: "Error: INSTALLATION FAILED"

**Symptoms:**
```
Error: INSTALLATION FAILED: cannot re-use a name that is still in use
```

**Solution:**
```bash
# Check if release already exists
helm list -n openclaw

# If it exists, uninstall first
helm uninstall openclaw -n openclaw

# Or install with different name
helm install openclaw2 oci://ghcr.io/thepagent/openclaw-helm -n openclaw
```

#### Issue: "failed to download" (OCI registry)

**Symptoms:**
```
Error: failed to download "oci://ghcr.io/thepagent/openclaw-helm"
```

**Solution:**
```bash
# Add repository and install from index
helm repo add openclaw https://thepagent.github.io/openclaw-helm
helm repo update
helm install openclaw openclaw/openclaw-helm -n openclaw
```

### 2. Pod Issues

#### Issue: Pod in "Pending" state

**Diagnosis:**
```bash
# Describe pod to see events
kubectl describe pod -n openclaw <pod-name>

# Common causes:
# - Insufficient resources
# - No nodes available
# - PVC pending
```

**Solutions:**

**Resource issues:**
```bash
# Check resource requests vs available
kubectl describe nodes

# Reduce resource requests in values.yaml
resources:
  requests:
    memory: "128Mi"  # Reduced from 256Mi
    cpu: "50m"       # Reduced from 100m
```

**PVC issues:**
```bash
# Check PVC status
kubectl get pvc -n openclaw

# If using k3s, ensure local-path provisioner is running
kubectl get pods -n local-path-storage
```

#### Issue: Pod in "CrashLoopBackOff" state

**Diagnosis:**
```bash
# Check pod logs
kubectl logs -n openclaw <pod-name> -c main --previous

# Common error patterns:
# - "Cannot find module" → Node.js dependency issue
# - "Permission denied" → Filesystem permissions
# - "Connection refused" → Database/API connectivity
```

**Solutions:**

**Module dependency issue:**
```bash
# Delete pod to force fresh start (will re-pull image)
kubectl delete pod -n openclaw <pod-name>

# Or update to newer image tag
helm upgrade openclaw oci://ghcr.io/thepagent/openclaw-helm --set image.tag=latest
```

**Permission issues:**
```yaml
# Update security context in values.yaml
securityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
```

#### Issue: Pod in "Error" state

**Diagnosis:**
```bash
# Get pod status details
kubectl get pod -n openclaw <pod-name> -o yaml | grep -A 10 "lastState:"

# Common causes:
# - OOMKilled (Out of Memory)
# - ImagePullBackOff
# - ContainerCannotRun
```

**Solutions:**

**OOMKilled:**
```yaml
# Increase memory limits
resources:
  limits:
    memory: "1Gi"  # Increased from 512Mi
```

**ImagePullBackOff:**
```bash
# Check image name and tag
kubectl describe pod -n openclaw <pod-name> | grep -i image

# Try pulling manually
docker pull ghcr.io/thepagent/openclaw:latest

# Use different tag
helm upgrade openclaw oci://ghcr.io/thepagent/openclaw-helm --set image.tag=stable
```

### 3. Gateway Issues

#### Issue: Gateway not responding on port 8000

**Diagnosis:**
```bash
# Check if gateway process is running
kubectl exec -n openclaw <pod-name> -- ps aux | grep gateway

# Test connectivity from within pod
kubectl exec -n openclaw <pod-name> -- curl -s http://localhost:8000

# Check gateway logs
kubectl logs -n openclaw <pod-name> -c main --tail=50
```

**Solutions:**

**Port conflict:**
```yaml
# Change gateway port
env:
  - name: OPENCLAW_GATEWAY_PORT
    value: "8080"
```

**Startup failure:**
```bash
# Increase initial delay for probes
livenessProbe:
  initialDelaySeconds: 60  # Increased from 30

readinessProbe:
  initialDelaySeconds: 30  # Increased from 5
```

#### Issue: "Cannot connect to AI provider"

**Diagnosis:**
```bash
# Check API key configuration
kubectl exec -n openclaw <pod-name> -- env | grep API_KEY

# Test API connectivity
kubectl exec -n openclaw <pod-name> -- curl -s https://api.openai.com/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY"
```

**Solutions:**

**Missing API keys:**
```bash
# Create secret with API keys
kubectl create secret generic openclaw-secrets \
  --from-literal=OPENAI_API_KEY=sk-... \
  --namespace openclaw

# Update deployment to use secret
envFrom:
  - secretRef:
      name: openclaw-secrets
```

**Invalid API keys:**
```bash
# Verify API key works
curl -s https://api.openai.com/v1/models \
  -H "Authorization: Bearer sk-..."

# If invalid, regenerate key and update secret
kubectl create secret generic openclaw-secrets \
  --from-literal=OPENAI_API_KEY=new-key-here \
  --namespace openclaw --dry-run=client -o yaml | \
  kubectl apply -f -
```

### 4. Storage Issues

#### Issue: "PVC pending" or "unbound"

**Diagnosis:**
```bash
# Check PVC details
kubectl describe pvc -n openclaw <pvc-name>

# Check storage class availability
kubectl get storageclass
```

**Solutions:**

**Missing storage class (k3s):**
```yaml
# Use local-path provisioner
persistence:
  storageClass: "local-path"
```

**Insufficient storage:**
```yaml
# Reduce requested size
persistence:
  size: "2Gi"  # Reduced from 5Gi
```

#### Issue: "Permission denied" on PVC

**Diagnosis:**
```bash
# Check pod logs for permission errors
kubectl logs -n openclaw <pod-name> | grep -i "permission\|EACCES"
```

**Solutions:**

**Fix permissions:**
```yaml
# Add proper security context
securityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
```

**Alternative: Use emptyDir (non-persistent):**
```yaml
persistence:
  enabled: false
  
extraVolumes:
  - name: data
    emptyDir: {}

extraVolumeMounts:
  - name: data
    mountPath: /home/node/.openclaw
```

### 5. Network Issues

#### Issue: Cannot access service from outside cluster

**Diagnosis:**
```bash
# Check service type and external IP
kubectl get svc -n openclaw

# Test from within cluster
kubectl run test --rm -i --tty --image=busybox -- \
  wget -qO- http://openclaw-openclaw-helm.openclaw.svc.cluster.local:8000
```

**Solutions:**

**Change service type:**
```yaml
# Use NodePort or LoadBalancer
service:
  type: NodePort
  # or
  type: LoadBalancer
```

**Port forward for testing:**
```bash
# Forward local port to service
kubectl port-forward -n openclaw svc/openclaw-openclaw-helm 8080:8000

# Now accessible at http://localhost:8080
```

#### Issue: DNS resolution issues

**Diagnosis:**
```bash
# Test DNS from within pod
kubectl exec -n openclaw <pod-name> -- nslookup api.openai.com

# Check CoreDNS/kube-dns
kubectl get pods -n kube-system | grep -E "coredns|kube-dns"
```

**Solutions:**

**Add DNS configuration:**
```yaml
# Custom DNS settings
dnsConfig:
  nameservers:
    - 8.8.8.8
    - 8.8.4.4
  searches:
    - openclaw.svc.cluster.local
    - svc.cluster.local
    - cluster.local
```

### 6. Skills Issues

#### Issue: Skills not installing

**Diagnosis:**
```bash
# Check init-skills container logs
kubectl logs -n openclaw <pod-name> -c init-skills

# Check skills directory
kubectl exec -n openclaw <pod-name> -- ls -la /home/node/.openclaw/workspace/skills/
```

**Solutions:**

**Network issues (skills download fails):**
```bash
# Install skills manually
kubectl exec -n openclaw <pod-name> -- \
  npx -y clawhub install weather --no-input

# Or specify skills in values.yaml
skills:
  - weather
  - healthcheck
```

**Insufficient storage for skills:**
```yaml
# Increase PVC size
persistence:
  size: "10Gi"  # Increased from 5Gi
```

### 7. Upgrade Issues

#### Issue: Upgrade fails or causes downtime

**Diagnosis:**
```bash
# Check upgrade history
helm history openclaw -n openclaw

# Check current pod status after upgrade
kubectl get pods -n openclaw -w
```

**Solutions:**

**Use strategic upgrade approach:**
```bash
# 1. Backup first
./backup_openclaw.sh

# 2. Upgrade with longer timeout
helm upgrade openclaw oci://ghcr.io/thepagent/openclaw-helm \
  -n openclaw --timeout 10m

# 3. If fails, rollback
helm rollback openclaw <previous-revision> -n openclaw
```

**Configure upgrade strategy:**
```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1
    maxUnavailable: 0
```

### 8. Performance Issues

#### Issue: High memory or CPU usage

**Diagnosis:**
```bash
# Check resource usage
kubectl top pods -n openclaw

# Check for memory leaks
kubectl logs -n openclaw <pod-name> | grep -i "memory\|heap"
```

**Solutions:**

**Adjust resource limits:**
```yaml
resources:
  requests:
    memory: "512Mi"
    cpu: "250m"
  limits:
    memory: "1Gi"
    cpu: "500m"
```

**Configure garbage collection:**
```yaml
env:
  - name: NODE_OPTIONS
    value: "--max-old-space-size=512"
```

#### Issue: Slow response times

**Diagnosis:**
```bash
# Check gateway response time
time kubectl exec -n openclaw <pod-name> -- \
  curl -s -o /dev/null -w "%{time_total}" http://localhost:8000
```

**Solutions:**

**Increase resources:**
```yaml
resources:
  requests:
    memory: "1Gi"
    cpu: "500m"
```

**Configure caching:**
```yaml
env:
  - name: OPENCLAW_CACHE_ENABLED
    value: "true"
  - name: OPENCLAW_CACHE_TTL
    value: "300"  # 5 minutes
```

### 9. Logging and Debugging

#### Collect Debug Information

```bash
#!/bin/bash
# debug-collector.sh

NAMESPACE=${1:-openclaw}
DEBUG_DIR="/tmp/openclaw-debug-$(date +%Y%m%d-%H%M%S)"

mkdir -p "$DEBUG_DIR"

echo "Collecting OpenClaw debug information..."

# Cluster info
kubectl cluster-info > "$DEBUG_DIR/cluster-info.txt"

# Namespace resources
kubectl get all -n $NAMESPACE > "$DEBUG_DIR/resources.txt"

# Pod details
POD=$(kubectl get pod -n $NAMESPACE -l app.kubernetes.io/name=openclaw-helm -o name | head -1)
if [ -n "$POD" ]; then
    kubectl describe $POD -n $NAMESPACE > "$DEBUG_DIR/pod-describe.txt"
    kubectl logs $POD -n $NAMESPACE -c main > "$DEBUG_DIR/pod-logs.txt"
    kubectl logs $POD -n $NAMESPACE -c init-skills > "$DEBUG_DIR/init-logs.txt"
fi

# Events
kubectl get events -n $NAMESPACE --sort-by='.lastTimestamp' > "$DEBUG_DIR/events.txt"

# PVC status
kubectl get pvc -n $NAMESPACE > "$DEBUG_DIR/pvc.txt"

# ConfigMaps and Secrets
kubectl get configmap,secret -n $NAMESPACE > "$DEBUG_DIR/config.txt"

echo "Debug information saved to: $DEBUG_DIR"
```

#### Enable Verbose Logging

```yaml
# debug-values.yaml
env:
  - name: OPENCLAW_LOG_LEVEL
    value: "debug"
  - name: DEBUG
    value: "*"
  - name: NODE_ENV
    value: "development"

resources:
  requests:
    memory: "1Gi"
    cpu: "500m"
```

### 10. k3s-Specific Issues

#### Issue: "permission denied" with sudo

**Solution:**
```bash
# Always use -E flag to preserve environment
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
sudo -E helm install openclaw oci://ghcr.io/thepagent/openclaw-helm -n openclaw
```

#### Issue: Local storage not working

**Solution:**
```yaml
# Ensure local-path storage class is used
persistence:
  storageClass: "local-path"
  size: "5Gi"
```

### Emergency Recovery Procedures

#### Complete Reinstallation

```bash
#!/bin/bash
# emergency-reinstall.sh

NAMESPACE=openclaw

echo "⚠️  Emergency reinstallation of OpenClaw..."

# 1. Backup existing data
./backup_openclaw.sh

# 2. Uninstall everything
helm uninstall openclaw -n $NAMESPACE --wait
kubectl delete pvc -n $NAMESPACE --all
kubectl delete secret -n $NAMESPACE --all

# 3. Wait for cleanup
sleep 30

# 4. Fresh install
helm install openclaw oci://ghcr.io/thepagent/openclaw-helm \
  -n $NAMESPACE --create-namespace \
  --set persistence.size=5Gi

# 5. Restore from backup if needed
# ./restore-<backup-name>.sh <backup-file>.tar.gz
```

#### Data Recovery from PVC

```bash
# Create debug pod with PVC attached
kubectl run recovery-pod -n openclaw \
  --image=busybox \
  --restart=Never \
  --overrides='
{
  "spec": {
    "containers": [{
      "name": "recovery",
      "image": "busybox",
      "command": ["sleep", "3600"],
      "volumeMounts": [{
        "mountPath": "/data",
        "name": "data"
      }]
    }],
    "volumes": [{
      "name": "data",
      "persistentVolumeClaim": {
        "claimName": "openclaw-openclaw-helm-pvc"
      }
    }]
  }
}'

# Copy data out
kubectl cp recovery-pod:/data ~/recovered-data -n openclaw

# Cleanup
kubectl delete pod recovery-pod -n openclaw
```

## Support Matrix

### Kubernetes Versions

| Version | Status | Notes |
|---------|--------|-------|
| 1.24+ | ✅ Supported | Fully compatible |
| 1.20-1.23 | ⚠️ Limited | May need adjustments |
| <1.20 | ❌ Not supported | Use newer version |

### Storage Classes

| Provider | Storage Class | Status |
|----------|---------------|--------|
| k3s | local-path | ✅ Recommended |
| AWS EKS | gp2/gp3 | ✅ Supported |
| GCP GKE | standard | ✅ Supported |
| Azure AKS | managed-premium | ✅ Supported |
| MicroK8s | microk8s-hostpath | ✅ Supported |

### Network Plugins

| Plugin | Status | Notes |
|--------|--------|-------|
| Calico | ✅ Supported | Default for many distros |
| Flannel | ✅ Supported | Common in k3s |
| Cilium | ✅ Supported | Advanced features available |
| Weave | ✅ Supported | No issues reported |