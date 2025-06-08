# LetsNormalizeIt - Blogging Platform Backend

A scalable, maintainable blogging platform backend built with Go, MongoDB, Redis, and Firebase Auth.

## Features

- Firebase Authentication
- Blog post management
- Likes and bookmarks
- User profiles
- Admin functionality
- Rate limiting
- Caching with Redis

## Architecture

This application follows a clean architecture approach with the following components:

- **Handlers**: HTTP request handlers
- **Services**: Business logic
- **Repositories**: Data access
- **Models**: Data structures
- **Middleware**: Request processing middleware

## Prerequisites

- Go 1.18+
- MongoDB
- Redis
- Firebase account with Firebase Auth enabled

## Setup

1. Clone the repository
2. Place your Firebase credentials file in the root directory as `firebase-credentials.json`
3. Configure the application in `configs/config.yaml`
4. Run the application:

```bash
go run cmd/server/main.go
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/signup`: Register a new user
- `POST /api/v1/auth/signin`: Sign in a user (Note: actual auth is done via Firebase SDK)

### Public Routes

- `GET /api/v1/blogs`: Get a list of blogs
- `GET /api/v1/blogs/:id`: Get a specific blog
- `GET /api/v1/blogs/:id/comments`: Get comments for a specific blog

### Protected Routes (require authentication)

- `POST /api/v1/blogs`: Create a new blog
- `POST /api/v1/blogs/:id/like`: Like a blog
- `POST /api/v1/blogs/:id/bookmark`: Bookmark a blog
- `POST /api/v1/comments`: Add a comment to a blog
- `GET /api/v1/user/bookmarks`: Get user bookmarks
- `GET /api/v1/user/profile`: Get user profile
- `PUT /api/v1/user/profile`: Update user profile

### Admin Routes

- `GET /api/v1/admin/users`: Get a list of users (admin only)
- `POST /api/v1/admin/users/:id/set-admin`: Set admin privileges for a user (admin only)

## Authentication Flow

### Client-Side Authentication

This backend is designed to work with Firebase Authentication. The client application should:

1. Initialize Firebase Auth in the client app
2. Use Firebase Auth UI or custom UI to sign up/sign in users
3. Get the ID token from Firebase Auth
4. Include the ID token in the Authorization header for requests:
   ```
   Authorization: Bearer <firebase-id-token>
   ```

### Backend Authentication

The backend will:

1. Verify the Firebase ID token
2. Extract the user ID from the token
3. Create/fetch the user profile from MongoDB
4. Allow or deny access to protected resources

## Firebase Admin SDK

This application uses the Firebase Admin SDK to:

- Verify ID tokens
- Access user information
- Create new users
- Update user profiles

## Development

### Running in Development Mode

```bash
go run cmd/server/main.go
```

### Building for Production

```bash
go build -o server cmd/server/main.go
```

### Testing

```bash
go test ./...
```
