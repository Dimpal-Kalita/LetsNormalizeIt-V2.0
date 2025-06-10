# Development Guide

## Hot Reload Development

This project uses [Air](https://github.com/air-verse/air) for hot reload functionality during development.

### Prerequisites

1. Ensure Air is installed:
   ```bash
   go install github.com/air-verse/air@latest
   ```

2. Make sure your environment variables are set up:
   ```bash
   make setup  # This will create .env from .env.example
   ```

### Running in Development Mode

Start the development server with hot reload:

```bash
make dev
```

This will:
- Create a `tmp/` directory for build artifacts
- Start the server with Air
- Automatically restart the server when you make changes to `.go`, `.json`, `.yaml`, or `.yml` files
- Display colored output for better readability

### Manual Commands

If you prefer to run Air directly:

```bash
air
```

Or to run without hot reload:

```bash
make build-run
```

## Logging

The application now uses structured logging with contextual information throughout all API handlers.

### Logger Features

- **Structured Logging**: All logs are in JSON format with consistent fields
- **Contextual Information**: Each API request includes operation, path, method, client IP, and user ID
- **Different Log Levels**: Debug, Info, Warn, Error, Fatal
- **Request Tracking**: Each operation is logged with relevant context

### Log Levels

- **Debug**: Detailed information for debugging (request body parsing, etc.)
- **Info**: General operational information (successful operations, user actions)
- **Warn**: Warning conditions (user not found, invalid requests)
- **Error**: Error conditions that need attention (database errors, auth failures)
- **Fatal**: Critical errors that cause the application to exit

### Example Log Entries

```json
{
  "level": "INFO",
  "timestamp": "2023-10-01T12:00:00Z",
  "caller": "user/handler.go:45",
  "message": "Processing user login request",
  "operation": "Login",
  "path": "/api/v1/user/login",
  "method": "POST",
  "clientIP": "192.168.1.1"
}
```

```json
{
  "level": "INFO", 
  "timestamp": "2023-10-01T12:00:05Z",
  "caller": "user/handler.go:67",
  "message": "User login successful",
  "operation": "Login",
  "path": "/api/v1/user/login", 
  "method": "POST",
  "clientIP": "192.168.1.1",
  "userID": "firebase_user_id_123",
  "email": "user@example.com"
}
```

## File Structure

```
.
├── .air.toml              # Air configuration for hot reload
├── tmp/                   # Temporary build files (auto-created)
├── build-errors.log       # Build error logs from Air
├── cmd/server/main.go     # Main server entry point
├── internal/
│   ├── user/handler.go    # User API handlers with logging
│   ├── utils/logger.go    # Logger utilities
│   └── ...
└── Makefile              # Build commands
```

## API Endpoints

### User Authentication

All user endpoints require Firebase ID token authentication via `Authorization: Bearer <token>` header.

- `POST /api/v1/user/login` - Login existing user
- `POST /api/v1/user/login-or-register` - Login or auto-register user  
- `POST /api/v1/user/register` - Register new user

Each endpoint logs the complete request flow with appropriate context.

## Development Tips

1. **Monitor Logs**: Watch the console output to see structured logs for each request
2. **Hot Reload**: Save any `.go` file to trigger automatic server restart
3. **Debug Mode**: Set log level to "debug" in configuration to see detailed request information
4. **Build Errors**: Check `build-errors.log` if the server fails to restart 