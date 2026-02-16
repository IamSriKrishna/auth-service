# Auth Service Kubernetes Deployment

ThisAll images are stored in Google Container Registry:
- Registry: `asia-south1-docker.pkg.dev/bb-app-461714/bbapp-images/github.com/bbapp-org/auth-service`
- Development tags: `<commit-sha>`
- Staging tags: `<calver-tag>` (e.g., `v25.01.1`) and `latest`ectory contains Kubernetes manifests for deploying the github.com/bbapp-org/auth-service using Kustomize with overlay pattern for dual deployment strategy.

## Directory Structure

```
k8s/
├── base/                           # Base configurations
│   ├── kustomization.yaml
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── configmap.yaml
│   └── hpa.yaml
├── overlays/
│   ├── dev/                        # Development environment (bbapp-dev namespace)
│   │   ├── kustomization.yaml
│   │   └── configmap-patch.yaml
│   └── stg/                        # Staging environment (bbapp-stg namespace)
│       ├── kustomization.yaml
│       └── configmap-patch.yaml
└── README.md
```

## Deployment Strategy

### Dual Deployment Pattern

- **Development Environment (bbapp-dev)**:
  - Triggered by pushes to `develop` branch
  - Uses image tag: `<commit-sha>`
  - Deployed to `bbapp-dev` namespace
  - Environment-specific configurations for development

- **Staging Environment (bbapp-stg)**:
  - Triggered by pushes to `main` branch
  - Uses image tag: `<calver-tag>` (e.g., `v25.01.1`) and `latest`
  - Deployed to `bbapp-stg` namespace
  - Environment-specific configurations for staging

### Image Registry

All images are stored in Google Container Registry:
- Registry: `asia-south1-docker.pkg.dev/bb-app-461714/bbapp-images/github.com/bbapp-org/auth-service`
- Development tags: `develop-<commit-sha>`
- Staging tags: `<version>` (CalVer) and `latest`

## Manual Deployment

### Prerequisites

1. Install `kubectl` and `kustomize`
2. Configure `kubectl` to connect to your GKE cluster
3. Ensure you have access to the target namespace

### Deploy to Development

```bash
# Deploy to development environment
kubectl apply -k k8s/overlays/dev

# Check deployment status
kubectl get deployments -n bbapp-dev
kubectl get pods -n bbapp-dev
```

### Deploy to Staging

```bash
# Deploy to staging environment
kubectl apply -k k8s/overlays/stg

# Check deployment status
kubectl get deployments -n bbapp-stg
kubectl get pods -n bbapp-stg
```

### Customize Image Tag

To deploy with a specific image tag:

```bash
# For development
cd k8s/overlays/dev
kustomize edit set image asia-south1-docker.pkg.dev/bb-app-461714/bbapp-images/github.com/bbapp-org/auth-service:abc123def456
kubectl apply -k .

# For staging
cd k8s/overlays/stg
kustomize edit set image asia-south1-docker.pkg.dev/bb-app-461714/bbapp-images/github.com/bbapp-org/auth-service:v25.01.1
kubectl apply -k .
```

## CI/CD Integration

### GitHub Actions Workflows

1. **Test Workflow** (`test.yml`):
   - Runs on pull requests
   - Executes unit tests and security scans

2. **Development Workflow** (`development.yml`):
   - Triggered by pushes to `develop` branch
   - Builds and pushes development image
   - Deploys to `bbapp-dev` namespace

3. **Release Workflow** (`release.yml`):
   - Triggered by pushes to `main` branch
   - Builds and pushes staging image with CalVer versioning
   - Deploys to `bbapp-stg` namespace
   - Creates GitHub release

### Required GitHub Secrets

- `GCP_WORKLOAD_IDENTITY_PROVIDER`: Google Cloud Workload Identity Provider
- `GCP_SERVICE_ACCOUNT`: Service account for GKE access
- `GITHUB_TOKEN`: For creating releases (automatically provided)

### Required GitHub Variables

- `GKE_CLUSTER_NAME`: Name of the GKE cluster
- `GCP_ZONE`: GCP zone of the cluster
- `GCP_PROJECT_ID`: Google Cloud Project ID

## Environment Configuration

### Development Environment

- Database: `auth_service_dev`
- App Environment: `development`
- Namespace: `bbapp-dev`
- Image Tag: `<commit-sha>`

### Staging Environment

- Database: `auth_service_stg`
- App Environment: `staging`
- Namespace: `bbapp-stg`
- Image Tag: `<calver-tag>` (e.g., `v25.01.1`)

## Monitoring and Troubleshooting

### Check Deployment Status

```bash
# Check deployment rollout
kubectl rollout status deployment/github.com/bbapp-org/auth-service -n bbapp-dev
kubectl rollout status deployment/github.com/bbapp-org/auth-service -n bbapp-stg

# Check pod logs
kubectl logs -f deployment/github.com/bbapp-org/auth-service -n bbapp-dev
kubectl logs -f deployment/github.com/bbapp-org/auth-service -n bbapp-stg

# Check service endpoints
kubectl get svc -n bbapp-dev
kubectl get svc -n bbapp-stg
```

### Health Checks

The service includes readiness and liveness probes:
- Health endpoint: `/v1/health`
- Readiness probe: Initial delay 5s, period 5s
- Liveness probe: Initial delay 30s, period 10s

### Horizontal Pod Autoscaling

HPA is configured with:
- Min replicas: 3
- Max replicas: 10
- CPU target: 70%
- Memory target: 80%

## Security

- All images are scanned with Trivy and Gosec
- Service uses non-root user
- Secrets are stored in Kubernetes secrets
- Network policies can be applied for additional security

## Support

For issues or questions regarding deployment, please check:
1. GitHub Actions logs for build/deploy issues
2. Kubernetes pod logs for runtime issues
3. Service health endpoints for application status
