#!/bin/bash
# OpenClaw Health Check Script
# Comprehensive health check for OpenClaw Kubernetes deployments

set -e

# Configuration
NAMESPACE="${NAMESPACE:-openclaw}"
COLOR_GREEN='\033[0;32m'
COLOR_RED='\033[0;31m'
COLOR_YELLOW='\033[1;33m'
COLOR_BLUE='\033[0;34m'
COLOR_RESET='\033[0m'

echo -e "${COLOR_BLUE}🔧 OpenClaw Health Check Utility${COLOR_RESET}"
echo "================================"
echo "Namespace: $NAMESPACE"
echo "Timestamp: $(date)"
echo ""

# Function to print status
print_status() {
    local status=$1
    local message=$2
    
    if [ "$status" = "OK" ]; then
        echo -e "${COLOR_GREEN}✓${COLOR_RESET} $message"
    elif [ "$status" = "WARNING" ]; then
        echo -e "${COLOR_YELLOW}⚠${COLOR_RESET} $message"
    else
        echo -e "${COLOR_RED}✗${COLOR_RESET} $message"
    fi
}

# Check 1: Namespace exists
echo -e "${COLOR_BLUE}1. Checking namespace...${COLOR_RESET}"
if kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
    print_status "OK" "Namespace '$NAMESPACE' exists"
else
    print_status "ERROR" "Namespace '$NAMESPACE' does not exist"
    exit 1
fi

# Check 2: Pod status
echo -e "\n${COLOR_BLUE}2. Checking pods...${COLOR_RESET}"
PODS=$(kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=openclaw-helm -o jsonpath='{.items[*].metadata.name}' 2>/dev/null || true)

if [ -z "$PODS" ]; then
    print_status "ERROR" "No OpenClaw pods found"
    echo "   Trying alternate label selector..."
    PODS=$(kubectl get pods -n $NAMESPACE -l app=openclaw -o jsonpath='{.items[*].metadata.name}' 2>/dev/null || true)
fi

if [ -z "$PODS" ]; then
    print_status "ERROR" "No OpenClaw pods found with any label"
    exit 1
fi

for POD in $PODS; do
    POD_STATUS=$(kubectl get pod -n $NAMESPACE $POD -o jsonpath='{.status.phase}')
    POD_READY=$(kubectl get pod -n $NAMESPACE $POD -o jsonpath='{.status.containerStatuses[0].ready}')
    RESTARTS=$(kubectl get pod -n $NAMESPACE $POD -o jsonpath='{.status.containerStatuses[0].restartCount}')
    
    if [ "$POD_STATUS" = "Running" ] && [ "$POD_READY" = "true" ]; then
        print_status "OK" "Pod $POD: Running, Ready"
    else
        print_status "WARNING" "Pod $POD: Status=$POD_STATUS, Ready=$POD_READY"
    fi
    
    if [ "$RESTARTS" -gt 0 ]; then
        print_status "WARNING" "Pod $POD: Restarted $RESTARTS times"
    fi
done

POD_NAME=$(echo $PODS | awk '{print $1}')

# Check 3: Service status
echo -e "\n${COLOR_BLUE}3. Checking services...${COLOR_RESET}"
SERVICES=$(kubectl get svc -n $NAMESPACE -l app.kubernetes.io/name=openclaw-helm -o jsonpath='{.items[*].metadata.name}' 2>/dev/null || true)

if [ -n "$SERVICES" ]; then
    for SVC in $SERVICES; do
        SVC_TYPE=$(kubectl get svc -n $NAMESPACE $SVC -o jsonpath='{.spec.type}')
        CLUSTER_IP=$(kubectl get svc -n $NAMESPACE $SVC -o jsonpath='{.spec.clusterIP}')
        print_status "OK" "Service $SVC: Type=$SVC_TYPE, ClusterIP=$CLUSTER_IP"
    done
else
    print_status "WARNING" "No services found (this may be normal for headless deployments)"
fi

# Check 4: PVC status
echo -e "\n${COLOR_BLUE}4. Checking storage...${COLOR_RESET}"
PVC=$(kubectl get pvc -n $NAMESPACE -l app.kubernetes.io/name=openclaw-helm -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || true)

if [ -n "$PVC" ]; then
    PVC_STATUS=$(kubectl get pvc -n $NAMESPACE $PVC -o jsonpath='{.status.phase}')
    CAPACITY=$(kubectl get pvc -n $NAMESPACE $PVC -o jsonpath='{.status.capacity.storage}' 2>/dev/null || echo "unknown")
    
    if [ "$PVC_STATUS" = "Bound" ]; then
        print_status "OK" "PVC $PVC: Bound ($CAPACITY)"
    else
        print_status "WARNING" "PVC $PVC: Status=$PVC_STATUS"
    fi
else
    print_status "WARNING" "No PVC found (stateless deployment)"
fi

# Check 5: Gateway connectivity
echo -e "\n${COLOR_BLUE}5. Checking gateway connectivity...${COLOR_RESET}"
if [ -n "$POD_NAME" ]; then
    # Check if gateway is responding
    if kubectl exec -n $NAMESPACE $POD_NAME -c main -- curl -s -f http://localhost:8000 >/dev/null 2>&1; then
        print_status "OK" "Gateway responding on localhost:8000"
        
        # Get gateway status
        GATEWAY_STATUS=$(kubectl exec -n $NAMESPACE $POD_NAME -c main -- curl -s http://localhost:8000 || echo "error")
        if echo "$GATEWAY_STATUS" | grep -q "gateway"; then
            print_status "OK" "Gateway endpoint returns valid response"
        fi
    else
        print_status "WARNING" "Gateway not responding on localhost:8000"
    fi
fi

# Check 6: Resource usage
echo -e "\n${COLOR_BLUE}6. Checking resource usage...${COLOR_RESET}"
if command -v kubectl-top &> /dev/null || kubectl top pods -n $NAMESPACE 2>/dev/null | grep -q "CPU"; then
    kubectl top pods -n $NAMESPACE 2>/dev/null | grep "$POD_NAME" || true
else
    print_status "INFO" "Metrics server not available"
fi

# Check 7: OpenClaw version
echo -e "\n${COLOR_BLUE}7. Checking OpenClaw version...${COLOR_RESET}"
if [ -n "$POD_NAME" ]; then
    VERSION=$(kubectl exec -n $NAMESPACE $POD_NAME -c main -- openclaw version 2>/dev/null | head -1 || echo "unknown")
    print_status "INFO" "OpenClaw version: $VERSION"
fi

# Check 8: Logs analysis (recent errors)
echo -e "\n${COLOR_BLUE}8. Checking recent logs for errors...${COLOR_RESET}"
if [ -n "$POD_NAME" ]; then
    ERROR_COUNT=$(kubectl logs -n $NAMESPACE $POD_NAME -c main --tail=100 2>/dev/null | grep -i -c "error\|fail\|panic\|exception" || true)
    
    if [ "$ERROR_COUNT" -eq 0 ]; then
        print_status "OK" "No recent errors in logs"
    elif [ "$ERROR_COUNT" -lt 5 ]; then
        print_status "WARNING" "Found $ERROR_COUNT error(s) in recent logs"
    else
        print_status "ERROR" "Found $ERROR_COUNT error(s) in recent logs"
    fi
    
    # Show last few log lines
    echo -e "\n${COLOR_BLUE}Recent log tail:${COLOR_RESET}"
    kubectl logs -n $NAMESPACE $POD_NAME -c main --tail=5 2>/dev/null | while read line; do
        echo "  $line"
    done || echo "  (cannot retrieve logs)"
fi

# Check 9: Skills directory
echo -e "\n${COLOR_BLUE}9. Checking skills installation...${COLOR_RESET}"
if [ -n "$POD_NAME" ]; then
    SKILLS_COUNT=$(kubectl exec -n $NAMESPACE $POD_NAME -c main -- bash -c 'ls -1 /home/node/.openclaw/workspace/skills/ 2>/dev/null | wc -l' || echo "0")
    
    if [ "$SKILLS_COUNT" -gt 0 ]; then
        print_status "OK" "Found $SKILLS_COUNT installed skill(s)"
        
        # List first few skills
        SKILLS_LIST=$(kubectl exec -n $NAMESPACE $POD_NAME -c main -- bash -c 'ls -1 /home/node/.openclaw/workspace/skills/ 2>/dev/null | head -5 | tr "\n" ","' || echo "")
        if [ -n "$SKILLS_LIST" ]; then
            echo "  Installed: ${SKILLS_LIST%,}"
        fi
    else
        print_status "WARNING" "No skills installed"
    fi
fi

# Summary
echo -e "\n${COLOR_BLUE}📋 Health Check Summary${COLOR_RESET}"
echo "================================"

# Count statuses
OK_COUNT=$(grep -c "✓" <<< "$(cat $0)" 2>/dev/null || echo "0")
WARNING_COUNT=$(grep -c "⚠" <<< "$(cat $0)" 2>/dev/null || echo "0")
ERROR_COUNT=$(grep -c "✗" <<< "$(cat $0)" 2>/dev/null || echo "0")

echo -e "${COLOR_GREEN}✓ Passed: $OK_COUNT${COLOR_RESET}"
echo -e "${COLOR_YELLOW}⚠ Warnings: $WARNING_COUNT${COLOR_RESET}"
echo -e "${COLOR_RED}✗ Errors: $ERROR_COUNT${COLOR_RESET}"

if [ "$ERROR_COUNT" -eq 0 ]; then
    if [ "$WARNING_COUNT" -eq 0 ]; then
        echo -e "\n${COLOR_GREEN}✅ All checks passed! OpenClaw is healthy.${COLOR_RESET}"
    else
        echo -e "\n${COLOR_YELLOW}⚠ OpenClaw is running with warnings.${COLOR_RESET}"
    fi
else
    echo -e "\n${COLOR_RED}❌ OpenClaw has issues that need attention.${COLOR_RESET}"
fi

echo -e "\n${COLOR_BLUE}🔧 Troubleshooting Commands:${COLOR_RESET}"
echo "  View logs: kubectl logs -n $NAMESPACE $POD_NAME -c main --tail=50"
echo "  Describe pod: kubectl describe pod -n $NAMESPACE $POD_NAME"
echo "  Restart deployment: kubectl rollout restart deployment -n $NAMESPACE"
echo "  Check events: kubectl get events -n $NAMESPACE --sort-by='.lastTimestamp'"

# Export for monitoring
if [ -n "$POD_NAME" ]; then
    echo -e "\n${COLOR_BLUE}📊 Quick Status Export:${COLOR_RESET}"
    cat > "/tmp/openclaw-health-$TIMESTAMP.json" << EOF
{
  "timestamp": "$(date -Iseconds)",
  "namespace": "$NAMESPACE",
  "pod": "$POD_NAME",
  "status": "$(if [ "$ERROR_COUNT" -eq 0 ]; then if [ "$WARNING_COUNT" -eq 0 ]; then "healthy"; else "warning"; fi else "unhealthy"; fi)",
  "checks": {
    "passed": $OK_COUNT,
    "warnings": $WARNING_COUNT,
    "errors": $ERROR_COUNT
  }
}
EOF
    echo "  Status exported to: /tmp/openclaw-health-$TIMESTAMP.json"
fi