# Staging Overlay

## Prerequisites

1. **Install Sealed Secrets Controller** (one-time per cluster):
   ```bash
   kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.34.0/controller.yaml
   ```

2. **Install kubeseal CLI**:
   ```bash
   brew install kubeseal
   ```

## Creating Sealed Secrets

Replace `sealed-secrets.yaml` with encrypted secrets:

```bash
# Get values from .env file or environment
kubectl create secret generic app-secrets \
  --namespace=staging \
  --from-literal=MONGODB_URI="${MONGODB_URI}" \
  --from-literal=MONGODB_DATABASE="${MONGODB_DATABASE}" \
  --from-literal=USER_SERVICE_DEFAULT_ADMIN_USERNAME="${USER_SERVICE_DEFAULT_ADMIN_USERNAME}" \
  --from-literal=USER_SERVICE_DEFAULT_ADMIN_PASSWORD="${USER_SERVICE_DEFAULT_ADMIN_PASSWORD}" \
  --from-literal=JWT_SECRET="${JWT_SECRET}" \
  --from-literal=JWT_EXPIRY="24h" \
  --from-literal=AXIOM_API_TOKEN="${AXIOM_API_TOKEN}" \
  --from-literal=AXIOM_ENDPOINT="us-east-1.aws.edge.axiom.co" \
  --from-literal=AXIOM_DATASET="traces" \
  --from-literal=AXIOM_METRICS_DATASET="metrics" \
  --from-literal=ENVIRONMENT="staging" \
  --dry-run=client -o yaml | \
kubeseal --format yaml > sealed-secrets.yaml
```

Commit the encrypted `sealed-secrets.yaml`.

## Deploying to Staging

```bash
# From repository root
kubectl apply -k k8s/overlays/staging

# Check if running
kubectl get pods -n staging
kubectl get ingress -n staging
kubectl get sealedSecret -n staging
```

## Updating Image Tags

```bash
cd k8s/overlays/staging
kustomize edit set image ghcr.io/ORG/user-service:NEW_TAG
kubectl apply -k .
```

## Ingress Configuration

Update `staging.yourdomain.com` in `ingress.yaml` to your actual domain.
