# Auth Service

A comprehensive authentication and authorization service built with Go and Fiber framework, following the Business Requirements Document (BRD) specifications.

## 🚀 Features

- **Multiple Authentication Methods**:
  - Email-based OTP authentication (passwordless)
  - Phone-based OTP authentication (passwordless)
  - Google OIDC integration
  - Password-based authentication (for admin/partner users)

- **User Types**:
  - **Mobile Users**: Passwordless authentication via OTP or Google OIDC
  - **Super Admin**: Full system access with user management capabilities
  - **Admin Users**: Administrative access with limited permissions
  - **Partner Users**: Business partner operations access

- **Security Features**:
  - JWT-based authentication with 7-day access tokens
  - 90-day refresh token validity
  - Secure password hashing with bcrypt
  - Redis-based session management and OTP storage
  - Rate limiting and security middleware
  - Role-based access control (RBAC)

- **API Features**:
  - RESTful API design with Swagger documentation
  - Comprehensive error handling and validation
  - Health check endpoint (`/health`)
  - Version endpoint (`/version`) with build information
  - Standardized JSON responses

- **Technology Stack**:
  - Go 1.21
  - Fiber web framework
  - MySQL database with GORM
  - Redis for caching and session management
  - JWT for token-based authentication
  - Docker containerization
  - Kubernetes deployment ready

## 🏗️ Architecture

The service follows a clean architecture pattern with:

- **Repository Pattern**: Database abstraction layer
- **Service Layer**: Business logic implementation
- **Handler Layer**: HTTP request/response handling
- **Middleware**: Authentication and authorization
- **DTO**: Data transfer objects for API contracts

```
github.com/bbapp-org/auth-services/
├── app/
│   ├── config/          # Database and Redis configuration
│   ├── dto/             # Data transfer objects
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Authentication middleware
│   ├── models/          # Database models
│   ├── repo/            # Repository layer
│   ├── routes/          # Route definitions
│   ├── services/        # Business logic
│   └── utils/           # Utility functions
├── scripts/             # Migration and seed scripts
├── k8s/                 # Kubernetes manifests
├── .github/workflows/   # CI/CD pipeline
└── docs/                # API documentation
```

## 📡 API Endpoints

### Public Endpoints (No Authentication)
- `POST /v1/auth/register/email` - Register with email
- `POST /v1/auth/register/phone` - Register with phone
- `POST /v1/auth/register/google` - Register with Google OIDC
- `POST /v1/auth/login/email` - Login with email OTP
- `POST /v1/auth/login/phone` - Login with phone OTP
- `POST /v1/auth/login/google` - Login with Google OIDC
- `POST /v1/auth/login/password` - Login with password (admin/partner)
- `POST /v1/auth/verify-otp` - Verify OTP
- `POST /v1/auth/validate-token` - Validate JWT token (internal)
- `GET /v1/health` - Health check

### Protected Endpoints (JWT Required)
- `POST /v1/auth/refresh-token` - Refresh access token
- `GET /v1/auth/user-info` - Get user information
- `POST /v1/auth/change-password` - Change password
- `POST /v1/auth/logout` - Logout user

### Admin Endpoints (Super Admin Only)
- `POST /v1/auth/admin/create-user` - Create admin/partner user
- `POST /v1/auth/admin/reset-password` - Reset user password
- `GET /v1/auth/admin/users` - List users with pagination
- `GET /v1/auth/admin/users/:id` - Get user details
- `PUT /v1/auth/admin/users/:id` - Update user
- `DELETE /v1/auth/admin/users/:id` - Delete user
- `PUT /v1/auth/admin/users/:id/status` - Update user status
- `PUT /v1/auth/admin/users/:id/role` - Update user role

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- MySQL 8.0+
- Redis 7.0+

### Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd github.com/bbapp-org/auth-services
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your database and Redis configuration
   ```

3. **Run the setup script**:
   ```bash
   ./setup.sh
   ```

4. **Start the service**:
   ```bash
   go run main.go
   ```

### Manual Setup

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Set up database**:
   ```sql
   CREATE DATABASE auth_service;
   ```

3. **Run migrations**:
   ```bash
   go run scripts/migrate/main.go
   ```

4. **Seed initial data**:
   ```bash
   go run scripts/seed/main.go
   ```

5. **Start the service**:
   ```bash
   make run
   ```

## 🐳 Docker Deployment

### Build and Run with Docker

```bash
# Build the image
docker build -t github.com/bbapp-org/auth-service .

# Run the container
docker run -p 3001:3001 \
  -e DB_HOST=your-mysql-host \
  -e DB_PASSWORD=your-password \
  -e REDIS_ADDR=your-redis-host:6379 \
  github.com/bbapp-org/auth-service
```

### Docker Compose (for development)

```yaml
version: '3.8'
services:
  github.com/bbapp-org/auth-service:
    build: .
    ports:
      - "3001:3001"
    environment:
      - DB_HOST=mysql
      - REDIS_ADDR=redis:6379
    depends_on:
      - mysql
      - redis
  
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: auth_service
    
  redis:
    image: redis:7-alpine
```

## ☸️ Kubernetes Deployment

### Deploy to Kubernetes

1. **Apply the manifests**:
   ```bash
   kubectl apply -f k8s/
   ```

2. **Check deployment status**:
   ```bash
   kubectl get pods -l app=github.com/bbapp-org/auth-service
   ```

3. **Access the service**:
   ```bash
   kubectl port-forward svc/github.com/bbapp-org/auth-service 3001:80
   ```

### Kubernetes Resources

- **Deployment**: Scalable pod management
- **Service**: Internal load balancing
- **Ingress**: External access via Traefik
- **ConfigMap**: Configuration management
- **Secret**: Sensitive data management
- **HPA**: Horizontal Pod Autoscaler

## 🔐 Authentication Flow

### Mobile Users (Passwordless)

1. User registers/logs in with email or phone
2. System generates and sends 6-digit OTP
3. User verifies OTP within 5 minutes
4. System issues JWT access token (7 days) and refresh token (90 days)

### Admin/Partner Users

1. User logs in with email/username and password
2. System validates credentials against bcrypt hash
3. System issues JWT tokens with appropriate role claims

### Google OIDC

1. User authenticates with Google
2. System validates Google token
3. System creates/updates user record
4. System issues JWT tokens

## 🛡️ Security Features

- **Password Security**: bcrypt hashing with salt
- **JWT Security**: HS256 signing with secret key
- **Rate Limiting**: Protection against brute force attacks
- **Session Management**: Redis-based session storage
- **RBAC**: Role-based access control with fine-grained permissions
- **Input Validation**: Request validation and sanitization
- **HTTPS Only**: Secure communication in production

## 📊 Database Schema

### Users Table
- Supports multiple identity types (email, phone, Google ID)
- Enforces constraints based on user type
- Tracks verification status and login history

### Roles Table
- Defines user roles and permissions
- JSON-based permission storage
- Supports hierarchical access control

### Tokens Table
- Manages refresh tokens
- Tracks token expiration and revocation
- Supports user session management

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test -v ./...

# Run with coverage
make test-coverage

# Run specific test
go test -v ./main_test.go
```

### Test Coverage
- Unit tests for utilities and services
- Integration tests for API endpoints
- Database transaction tests
- Authentication flow tests

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | MySQL host | `localhost` |
| `DB_PORT` | MySQL port | `3306` |
| `DB_USER` | MySQL user | `root` |
| `DB_PASSWORD` | MySQL password | `` |
| `DB_NAME` | Database name | `auth_service` |
| `REDIS_ADDR` | Redis address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password | `` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key...` |
| `PORT` | Service port | `3001` |

### Default Credentials

**Super Admin User**:
- Email: `admin@example.com`
- Password: `admin123`

*Change these credentials immediately in production!*

## 🚦 Health Checks

The service provides health check endpoints:

- `GET /health` - Basic health check
- `GET /v1/health` - Detailed health check with version info

## 📈 Monitoring

- **Metrics**: Service metrics endpoint for Prometheus
- **Logging**: Structured logging with different levels
- **Tracing**: Request tracing for debugging
- **Health Checks**: Kubernetes-ready health probes

## 🔄 CI/CD Pipeline

The service includes a comprehensive CI/CD pipeline:

- **Testing**: Unit and integration tests
- **Linting**: Code quality checks
- **Security**: Security vulnerability scanning
- **Building**: Docker image building
- **Deployment**: Automated Kubernetes deployment

## 📚 API Documentation

- **Swagger UI**: Available at `/swagger/` when running
- **OpenAPI**: Full API specification included
- **Postman**: Collection available for testing

## 🛠️ Development

### Local Development

```bash
# Install air for hot reloading
go install github.com/cosmtrek/air@latest

# Run in development mode
make dev
```

### Available Make Commands

```bash
make build          # Build the application
make run            # Run the application
make test           # Run tests
make test-coverage  # Run tests with coverage
make migrate        # Run database migrations
make seed           # Seed database with initial data
make docker-build   # Build Docker image
make clean          # Clean build artifacts
make lint           # Run linter
make fmt            # Format code
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📋 Roadmap

- [ ] Multi-factor authentication (MFA)
- [ ] OAuth2 provider support (GitHub, Microsoft)
- [ ] User profile management endpoints
- [ ] Advanced audit logging
- [ ] Rate limiting per user
- [ ] Geo-location based access control
- [ ] API versioning
- [ ] GraphQL API support

## 🚀 Deployment

### Dual Deployment Strategy

The service follows a dual deployment pattern:

- **Development Environment**: Push to `develop` branch deploys to `bbapp-dev` namespace
- **Staging Environment**: Push to `main` branch deploys to `bbapp-stg` namespace

### Kubernetes Deployment

The service is deployed using Kustomize with overlay pattern:

```
k8s/
├── base/                    # Base configurations
├── overlays/
│   ├── dev/                # Development environment
│   └── stg/                # Staging environment
```

#### Prerequisites

1. Access to GKE cluster
2. `kubectl` configured
3. Required GitHub secrets and variables set

#### Deployment Process

1. **Development Deployment**:
   ```bash
   # Automatic on push to develop branch
   # Or manual deployment:
   kubectl apply -k k8s/overlays/dev
   ```

2. **Staging Deployment**:
   ```bash
   # Automatic on push to main branch
   # Or manual deployment:
   kubectl apply -k k8s/overlays/stg
   ```

#### CI/CD Workflows

- **Test Workflow**: Runs on pull requests
- **Development Workflow**: Builds and deploys to dev on push to develop
- **Release Workflow**: Builds, versions, and deploys to staging on push to main

#### Required GitHub Configuration

**Secrets**:
- `GCP_WORKLOAD_IDENTITY_PROVIDER`: Google Cloud Workload Identity Provider
- `GCP_SERVICE_ACCOUNT`: Service account for GKE access

**Variables**:
- `GKE_CLUSTER_NAME`: Name of the GKE cluster
- `GCP_ZONE`: GCP zone of the cluster
- `GCP_PROJECT_ID`: Google Cloud Project ID

### Docker Images

Images are stored in Google Container Registry:
- Registry: `asia-south1-docker.pkg.dev/bb-app-461714/bbapp-images/github.com/bbapp-org/auth-service`
- Development tags: `<commit-sha>`
- Staging tags: `<calver-tag>` (e.g., `v25.01.1`) and `latest`

For detailed deployment instructions, see [k8s/README.md](k8s/README.md).

## 🔧 Development

### Local Development

1. **Set up dependencies**:
   ```bash
   make setup
   ```

2. **Run tests**:
   ```bash
   make test
   ```

3. **Start development server**:
   ```bash
   make run
   ```

4. **Build binary**:
   ```bash
   make build
   ```

### Docker Development

1. **Build Docker image**:
   ```bash
   docker build -t github.com/bbapp-org/auth-service .
   ```

2. **Run with Docker Compose**:
   ```bash
   docker-compose up -d
   ```

## 🆘 Troubleshooting

### Common Issues

1. **Database Connection Failed**:
   - Check MySQL service is running
   - Verify database credentials in `.env`
   - Ensure database exists

2. **Redis Connection Failed**:
   - Check Redis service is running
   - Verify Redis address in `.env`
   - Check Redis authentication

3. **JWT Token Invalid**:
   - Verify JWT secret configuration
   - Check token expiration
   - Ensure proper token format

4. **OTP Not Received**:
   - Check notification service integration
   - Verify Redis OTP storage
   - Check email/SMS configuration

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the troubleshooting guide

---

**Built with ❤️ for secure authentication**
# auth-service
