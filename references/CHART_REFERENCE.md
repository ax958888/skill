# OpenClaw Helm Chart Reference

## Chart Values Reference

### Global Configuration

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `replicaCount` | int | `1` | Number of replicas |
| `image.repository` | string | `ghcr.io/thepagent/openclaw` | Container image repository |
| `image.tag` | string | `latest` | Container image tag |
| `image.pullPolicy` | string | `IfNotPresent` | Image pull policy |
| `imagePullSecrets` | list | `[]` | Image pull secrets |

### Resources

```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "100m"
  limits:
    memory: "512Mi"
    cpu: "250m"
```

### Persistence

```yaml
persistence:
  enabled: true
  storageClass: ""  # Use default storage class
  accessModes:
    - ReadWriteOnce
  size: 5Gi
  annotations: {}
```

### Environment Variables

#### Basic Configuration

```yaml
env:
  - name: OPENCLAW_LOG_LEVEL
    value: "info"
  - name: OPENCLAW_GATEWAY_PORT
    value: "8000"
  - name: NODE_ENV
    value: "production"
```

#### AI Provider Configuration

```yaml
env:
  - name: OPENAI_API_KEY
    valueFrom:
      secretKeyRef:
        name: openclaw-secrets
        key: openai-api-key
  - name: ANTHROPIC_API_KEY
    valueFrom:
      secretKeyRef:
        name: openclaw-secrets
        key: anthropic-api-key
  - name: GROQ_API_KEY
    valueFrom:
      secretKeyRef:
        name: openclaw-secrets
        key: groq-api-key
```

### Secrets Configuration

Create a secret first:

```bash
kubectl create secret generic openclaw-secrets \
  --from-literal=openai-api-key=sk-... \
  --from-literal=anthropic-api-key=sk-... \
  --namespace openclaw
```

Then reference in values.yaml:

```yaml
envFrom:
  - secretRef:
      name: openclaw-secrets
```

### Skills Configuration

```yaml
skills:
  - weather
  - healthcheck
  - coding-agent
  - openai-whisper-api
  - openai-image-gen
```

### Service Configuration

```yaml
service:
  type: ClusterIP
  port: 8000
  annotations: {}
```

### Ingress Configuration

```yaml
ingress:
  enabled: false
  className: ""
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: openclaw.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: openclaw-tls
      hosts:
        - openclaw.example.com
```

### Pod Configuration

```yaml
podAnnotations: {}
podLabels: {}
nodeSelector: {}
tolerations: []
affinity: {}
securityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
```

### Liveness and Readiness Probes

```yaml
livenessProbe:
  httpGet:
    path: /
    port: 8000
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /
    port: 8000
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 1
```

## Complete Values Example

### Production Configuration

```yaml
# production-values.yaml
replicaCount: 2

image:
  repository: ghcr.io/thepagent/openclaw
  tag: stable
  pullPolicy: IfNotPresent

resources:
  requests:
    memory: "1Gi"
    cpu: "500m"
  limits:
    memory: "2Gi"
    cpu: "1000m"

persistence:
  storageClass: "fast-ssd"
  accessModes:
    - ReadWriteOnce
  size: 20Gi

skills:
  - weather
  - healthcheck
  - coding-agent
  - openai-whisper-api
  - openai-image-gen
  - gemini

env:
  - name: OPENCLAW_LOG_LEVEL
    value: "info"
  - name: NODE_ENV
    value: "production"

envFrom:
  - secretRef:
      name: openclaw-secrets

service:
  type: ClusterIP
  port: 8000

ingress:
  enabled: true
  className: "nginx"
  annotations:
    kubernetes.io/ingress.class: nginx
  hosts:
    - host: openclaw.company.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: openclaw-tls
      hosts:
        - openclaw.company.com
```

### Development Configuration

```yaml
# development-values.yaml
replicaCount: 1

resources:
  requests:
    memory: "256Mi"
    cpu: "100m"
  limits:
    memory: "512Mi"
    cpu: "250m"

persistence:
  size: 2Gi

skills: []

env:
  - name: OPENCLAW_LOG_LEVEL
    value: "debug"
  - name: NODE_ENV
    value: "development"

service:
  type: ClusterIP
  port: 8000
```

### High-Availability Configuration

```yaml
# ha-values.yaml
replicaCount: 3

resources:
  requests:
    memory: "512Mi"
    cpu: "250m"
  limits:
    memory: "1Gi"
    cpu: "500m"

persistence:
  storageClass: "csi-rbd"  # Shared storage for RWX
  accessModes:
    - ReadWriteMany
  size: 50Gi

affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchExpressions:
          - key: app.kubernetes.io/name
            operator: In
            values:
            - openclaw-helm
        topologyKey: kubernetes.io/hostname

skills:
  - weather
  - healthcheck
  - coding-agent

service:
  type: LoadBalancer
  port: 8000
```

## Network Policies

### Restrict Access

```yaml
# network-policy.yaml
networkPolicy:
  enabled: true
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: monitoring
      ports:
        - port: 8000
          protocol: TCP
    - from:
        - ipBlock:
            cidr: 10.0.0.0/8
      ports:
        - port: 8000
          protocol: TCP
```

## Custom Resource Definitions

### Horizontal Pod Autoscaler

```yaml
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
  targetMemoryUtilizationPercentage: 80
```

### Pod Disruption Budget

```yaml
pdb:
  enabled: true
  minAvailable: 1
```

## Storage Classes by Provider

### AWS EKS

```yaml
persistence:
  storageClass: "gp2"
  size: "10Gi"
```

### GCP GKE

```yaml
persistence:
  storageClass: "standard"
  size: "10Gi"
```

### Azure AKS

```yaml
persistence:
  storageClass: "managed-premium"
  size: "10Gi"
```

### k3s

```yaml
persistence:
  storageClass: "local-path"
  size: "5Gi"
```

## Monitoring Configuration

### Prometheus ServiceMonitor

```yaml
serviceMonitor:
  enabled: true
  interval: 30s
  scrapeTimeout: 10s
  labels:
    release: prometheus
```

### Grafana Dashboard

```yaml
grafana:
  dashboard:
    enabled: true
    label: openclaw
```

## Troubleshooting Values

### Debug Mode

```yaml
# debug-values.yaml
replicaCount: 1

resources:
  requests:
    memory: "1Gi"
    cpu: "500m"
  limits:
    memory: "2Gi"
    cpu: "1000m"

env:
  - name: OPENCLAW_LOG_LEVEL
    value: "debug"
  - name: NODE_OPTIONS
    value: "--inspect=0.0.0.0:9229"
  - name: DEBUG
    value: "*"

livenessProbe:
  initialDelaySeconds: 120
  periodSeconds: 30

readinessProbe:
  initialDelaySeconds: 60
  periodSeconds: 30

persistence:
  size: "10Gi"
```

### Resource Investigation

```yaml
# resource-investigation.yaml
resources:
  requests:
    memory: "100Mi"
    cpu: "50m"
  limits:
    memory: "200Mi"
    cpu: "100m"

env:
  - name: OPENCLAW_LOG_LEVEL
    value: "trace"

skills: []
```

## Upgrade Configuration

### Preserve PVC During Upgrade

```yaml
# upgrade-values.yaml
persistence:
  existingClaim: "openclaw-openclaw-helm-pvc"  # Use existing PVC

strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1
    maxUnavailable: 0
```

### Canary Deployment

```yaml
# canary-values.yaml
replicaCount: 2

strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1
    maxUnavailable: 50%

podAnnotations:
  sidecar.istio.io/inject: "true"

skills:
  - weather  # Only essential skills for canary
```

## Security Hardening

### Security Context

```yaml
# security-values.yaml
securityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
  runAsNonRoot: true
  allowPrivilegeEscalation: false
  capabilities:
    drop:
      - ALL
  readOnlyRootFilesystem: true
  seccompProfile:
    type: RuntimeDefault

podSecurityContext:
  seccompProfile:
    type: RuntimeDefault

serviceAccount:
  create: true
  name: openclaw-service-account
  annotations: {}

networkPolicy:
  enabled: true
```

### Pod Security Standards

```yaml
# pss-values.yaml
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 1000

securityContext:
  allowPrivilegeEscalation: false
  capabilities:
    drop:
      - ALL
  readOnlyRootFilesystem: true
  seccompProfile:
    type: RuntimeDefault
```

## Custom Templates

### Custom Init Scripts

```yaml
# custom-init-values.yaml
extraInitContainers:
  - name: custom-setup
    image: busybox:latest
    command: ['sh', '-c', 'echo "Custom setup complete"']
    volumeMounts:
      - name: data
        mountPath: /home/node/.openclaw

extraVolumes:
  - name: config-map
    configMap:
      name: openclaw-custom-config

extraVolumeMounts:
  - name: config-map
    mountPath: /etc/openclaw/custom
```

### Sidecar Containers

```yaml
# sidecar-values.yaml
sidecars:
  - name: log-forwarder
    image: fluent/fluent-bit:latest
    env:
      - name: FLUENT_BIT_CONFIG
        value: "/fluent-bit/etc/fluent-bit.conf"
    volumeMounts:
      - name: data
        mountPath: /var/log/openclaw
      - name: config
        mountPath: /fluent-bit/etc

extraVolumes:
  - name: config
    configMap:
      name: fluent-bit-config
```