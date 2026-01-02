#!/bin/bash

# Deploy Production Service to Kubernetes
# Usage: ./deploy-k8s.sh [namespace]

NAMESPACE=${1:-golunch}

echo "🏭 Deploying Production Service to namespace: ${NAMESPACE}"
echo "💰 Cost: $0 (using PostgreSQL StatefulSet)"

# Create namespace if it doesn't exist
kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -

echo "🗄️ Deploying PostgreSQL..."
kubectl apply -f k8s/postgres-statefulset.yaml -n ${NAMESPACE}

echo "⏳ Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres-production -n ${NAMESPACE} --timeout=300s

echo "📦 Applying ConfigMap..."
kubectl apply -f k8s/operation-service-configmap.yaml -n ${NAMESPACE}

echo "🔐 Applying Secrets..."
kubectl apply -f k8s/operation-service-secrets.yaml -n ${NAMESPACE}

echo "🚀 Applying Deployment..."
kubectl apply -f k8s/operation-service-deployment.yaml -n ${NAMESPACE}

echo "🌐 Applying Service..."
kubectl apply -f k8s/operation-service-service.yaml -n ${NAMESPACE}

echo "📈 Applying HPA..."
kubectl apply -f k8s/operation-service-hpa.yaml -n ${NAMESPACE}

# Wait for deployment to be ready
echo "⏳ Waiting for Production Service to be ready..."
kubectl rollout status deployment/operation-service -n ${NAMESPACE} --timeout=300s

# Show deployment status
echo ""
echo "✅ Production Service Deployment Status:"
kubectl get pods -l app=operation-service -n ${NAMESPACE}
kubectl get pods -l app=postgres-production -n ${NAMESPACE}
kubectl get svc -n ${NAMESPACE} | grep production

echo ""
echo "🎉 Production Service deployed successfully!"
echo ""
echo "📊 Next Steps:"
echo "  • Test: kubectl port-forward svc/operation-service 8083:8083 -n ${NAMESPACE}"
echo "  • Check: curl http://localhost:8083/ping"
echo "  • Logs: kubectl logs -f deployment/operation-service -n ${NAMESPACE}"
echo "  • DB Access: kubectl port-forward svc/postgres-production 5432:5432 -n ${NAMESPACE}"