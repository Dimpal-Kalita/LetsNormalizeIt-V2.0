# LetsNormalizeIt-V2.0

A scalable blogging platform backend built with Go.

## Quick Start with Docker

### Development Environment (with Redis caching)
```bash
cd deploy/docker
docker-compose up -d
```

### Production Environment (minimal dependencies)
```bash
cd deploy/docker
docker-compose -f docker-compose.prod.yml up -d
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LNI_SERVER_PORT` | `8080` | Server port |
| `LNI_MONGODB_URI` | `mongodb://localhost:27017` | MongoDB connection URI |
| `LNI_MONGODB_DATABASE` | `letsnormalizeit` | MongoDB database name |
| `LNI_REDIS_ADDRESS` | `localhost:6379` | Redis address (optional) |
| `LNI_AUTH_TYPE` | `jwt` | Auth type: `jwt` or `firebase` |
| `LNI_JWT_SECRET` | - | JWT signing secret (required for JWT auth) |
| `LNI_JWT_EXPIRY` | `24h` | JWT token expiry |
| `LNI_CACHE_ENABLED` | `true` | Enable Redis caching |
| `LNI_LOGGER_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |

## Without Docker

### Prerequisites
- Go 1.24+
- MongoDB 7+
- Redis (optional, for caching)

### Running
```bash
# Copy environment file
cp .env.example .env

# Edit .env with your configuration
# For JWT auth: set LNI_AUTH_TYPE=jwt and LNI_JWT_SECRET=your-secret

# Install dependencies
go mod tidy

# Run the server
go run cmd/server/main.go
```

## API Endpoints

### Public
- `GET /api/v1/blogs` - List blogs
- `GET /api/v1/blogs/:id` - Get blog
- `GET /api/v1/blogs/:id/comments` - Get comments

### Authenticated
- `POST /api/v1/user/register` - Register user
- `POST /api/v1/user/login` - Login
- `POST /api/v1/blogs` - Create blog
- `POST /api/v1/blogs/:id/like` - Toggle like
- `POST /api/v1/blogs/:id/bookmark` - Toggle bookmark
- `GET /api/v1/user/profile` - Get profile
- `PUT /api/v1/user/profile` - Update profile

### Admin
- `GET /api/v1/admin/users` - List users
- `POST /api/v1/admin/users/:id/set-admin` - Set admin
