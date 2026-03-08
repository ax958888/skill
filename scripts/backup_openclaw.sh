#!/bin/bash
# OpenClaw Backup Script
# Automatically backs up OpenClaw configuration and data from Kubernetes deployment

set -e

# Configuration
NAMESPACE="${NAMESPACE:-openclaw}"
BACKUP_DIR="${BACKUP_DIR:-$HOME/openclaw-backups}"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_NAME="openclaw-backup-$TIMESTAMP"
FULL_BACKUP_PATH="$BACKUP_DIR/$BACKUP_NAME"

echo "🔧 OpenClaw Backup Utility"
echo "=========================="

# Check prerequisites
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found. Please install kubectl first."
    exit 1
fi

# Create backup directory
mkdir -p "$FULL_BACKUP_PATH"

echo "📊 Backup Information:"
echo "  Namespace: $NAMESPACE"
echo "  Backup directory: $FULL_BACKUP_PATH"
echo "  Timestamp: $TIMESTAMP"
echo ""

# Get pod name
POD_NAME=$(kubectl get pod -n $NAMESPACE -l app.kubernetes.io/name=openclaw-helm -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || true)

if [ -z "$POD_NAME" ]; then
    echo "⚠️  No OpenClaw pod found in namespace '$NAMESPACE'"
    echo "   Trying alternate label selector..."
    POD_NAME=$(kubectl get pod -n $NAMESPACE -l app=openclaw -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || true)
fi

if [ -z "$POD_NAME" ]; then
    echo "❌ Could not find OpenClaw pod. Please check:"
    echo "   1. Namespace '$NAMESPACE' exists"
    echo "   2. OpenClaw is running in this namespace"
    echo "   3. You have kubectl access to the cluster"
    exit 1
fi

echo "✅ Found pod: $POD_NAME"

# Create backup manifest
echo "📝 Creating backup manifest..."
cat > "$FULL_BACKUP_PATH/manifest.json" << EOF
{
  "backup": {
    "name": "$BACKUP_NAME",
    "timestamp": "$(date -Iseconds)",
    "namespace": "$NAMESPACE",
    "pod": "$POD_NAME",
    "cluster": "$(kubectl config current-context 2>/dev/null || echo "unknown")"
  },
  "openclaw": {
    "version": "$(kubectl exec -n $NAMESPACE $POD_NAME -- openclaw version 2>/dev/null | head -1 || echo "unknown")"
  }
}
EOF

# Backup configuration
echo "💾 Backing up OpenClaw configuration..."
kubectl cp "$NAMESPACE/$POD_NAME:/home/node/.openclaw" "$FULL_BACKUP_PATH/config" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    CONFIG_SIZE=$(du -sh "$FULL_BACKUP_PATH/config" | cut -f1)
    echo "✅ Configuration backed up ($CONFIG_SIZE)"
else
    echo "⚠️  Could not copy configuration. The pod may not be ready."
    echo "   Trying alternative approach..."
    
    # Try getting specific files
    mkdir -p "$FULL_BACKUP_PATH/config"
    kubectl exec -n $NAMESPACE $POD_NAME -- tar czf - -C /home/node/.openclaw . 2>/dev/null | \
        tar xzf - -C "$FULL_BACKUP_PATH/config" 2>/dev/null || true
    
    if [ -f "$FULL_BACKUP_PATH/config/openclaw.json" ]; then
        CONFIG_SIZE=$(du -sh "$FULL_BACKUP_PATH/config" | cut -f1)
        echo "✅ Configuration backed up via tar ($CONFIG_SIZE)"
    else
        echo "❌ Failed to backup configuration"
    fi
fi

# Backup Kubernetes resources
echo "📦 Backing up Kubernetes resources..."
kubectl get deployment -n $NAMESPACE -l app.kubernetes.io/name=openclaw-helm -o yaml > "$FULL_BACKUP_PATH/deployment.yaml" 2>/dev/null || true
kubectl get service -n $NAMESPACE -l app.kubernetes.io/name=openclaw-helm -o yaml > "$FULL_BACKUP_PATH/service.yaml" 2>/dev/null || true
kubectl get pvc -n $NAMESPACE -l app.kubernetes.io/name=openclaw-helm -o yaml > "$FULL_BACKUP_PATH/pvc.yaml" 2>/dev/null || true
kubectl get secret -n $NAMESPACE -l app.kubernetes.io/name=openclaw-helm -o yaml > "$FULL_BACKUP_PATH/secret.yaml" 2>/dev/null || true

echo "✅ Kubernetes resources backed up"

# Create compressed archive
echo "📦 Creating compressed archive..."
cd "$BACKUP_DIR"
tar -czf "$BACKUP_NAME.tar.gz" "$BACKUP_NAME"

ARCHIVE_SIZE=$(du -h "$BACKUP_NAME.tar.gz" | cut -f1)
echo "✅ Archive created: $BACKUP_NAME.tar.gz ($ARCHIVE_SIZE)"

# Cleanup uncompressed directory
rm -rf "$FULL_BACKUP_PATH"

# Create restore script
cat > "$BACKUP_DIR/restore-$BACKUP_NAME.sh" << 'EOF'
#!/bin/bash
# OpenClaw Restore Script
# Usage: ./restore-<backup-name>.sh [namespace]

set -e

BACKUP_FILE="$1"
NAMESPACE="${2:-openclaw}"

if [ -z "$BACKUP_FILE" ] || [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Usage: $0 <backup-file.tar.gz> [namespace]"
    exit 1
fi

echo "🔧 OpenClaw Restore Utility"
echo "=========================="

# Extract backup
EXTRACT_DIR=$(mktemp -d)
echo "📂 Extracting backup to $EXTRACT_DIR..."
tar -xzf "$BACKUP_FILE" -C "$EXTRACT_DIR"

BACKUP_DIR=$(find "$EXTRACT_DIR" -type d -name "openclaw-backup-*" | head -1)
if [ -z "$BACKUP_DIR" ]; then
    echo "❌ Invalid backup format"
    rm -rf "$EXTRACT_DIR"
    exit 1
fi

# Check namespace exists
if ! kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
    echo "📦 Creating namespace $NAMESPACE..."
    kubectl create namespace "$NAMESPACE"
fi

# Get current pod name
POD_NAME=$(kubectl get pod -n $NAMESPACE -l app.kubernetes.io/name=openclaw-helm -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || true)

if [ -z "$POD_NAME" ]; then
    echo "⚠️  No running OpenClaw pod found in namespace '$NAMESPACE'"
    echo "   You need to install OpenClaw first before restoring backup."
    echo "   Install with: helm install openclaw oci://ghcr.io/thepagent/openclaw-helm -n $NAMESPACE"
    rm -rf "$EXTRACT_DIR"
    exit 1
fi

echo "✅ Found pod: $POD_NAME"

# Restore configuration
if [ -d "$BACKUP_DIR/config" ]; then
    echo "💾 Restoring OpenClaw configuration..."
    
    # Stop the pod to prevent writes during restore
    echo "⏸️  Scaling deployment to 0 replicas..."
    kubectl scale deployment -n $NAMESPACE openclaw-openclaw-helm --replicas=0
    
    # Wait for pod termination
    sleep 10
    
    # Copy backup
    kubectl cp "$BACKUP_DIR/config" "$NAMESPACE/$POD_NAME:/home/node/.openclaw"
    
    # Restart
    echo "▶️  Scaling deployment back to 1 replica..."
    kubectl scale deployment -n $NAMESPACE openclaw-openclaw-helm --replicas=1
    
    echo "✅ Configuration restored"
else
    echo "⚠️  No configuration found in backup"
fi

# Cleanup
rm -rf "$EXTRACT_DIR"

echo ""
echo "🎉 Restore complete!"
echo ""
echo "Next steps:"
echo "1. Verify the pod is running: kubectl get pods -n $NAMESPACE"
echo "2. Check logs: kubectl logs -n $NAMESPACE $POD_NAME -c main --tail=50"
echo "3. Test functionality: kubectl exec -n $NAMESPACE $POD_NAME -- openclaw models status"
EOF

chmod +x "$BACKUP_DIR/restore-$BACKUP_NAME.sh"

echo ""
echo "🎉 Backup completed successfully!"
echo ""
echo "📋 Summary:"
echo "  Backup file: $BACKUP_DIR/$BACKUP_NAME.tar.gz"
echo "  Restore script: $BACKUP_DIR/restore-$BACKUP_NAME.sh"
echo ""
echo "🔧 Restore command:"
echo "  $BACKUP_DIR/restore-$BACKUP_NAME.sh $BACKUP_DIR/$BACKUP_NAME.tar.gz"
echo ""
echo "💡 Tip: Schedule regular backups with cron:"
echo "  0 2 * * * $BACKUP_DIR/backup_openclaw.sh"
echo ""

# Retention policy (keep last 7 days)
echo "🧹 Cleaning up old backups (keeping last 7 days)..."
find "$BACKUP_DIR" -name "openclaw-backup-*.tar.gz" -mtime +7 -delete 2>/dev/null || true
find "$BACKUP_DIR" -name "restore-openclaw-backup-*.sh" -mtime +7 -delete 2>/dev/null || true

echo "✅ Cleanup complete"