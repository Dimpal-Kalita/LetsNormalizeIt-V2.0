# Firebase ID Token Creator

A comprehensive tool for generating Firebase ID tokens directly for testing your backend APIs.

## 🚀 Features

- **Direct ID Token Generation**: Create Firebase ID tokens ready for immediate use
- **Admin Token Support**: Generate ID tokens with admin privileges
- **Web Interface**: Browser-based token generator with real-time authentication
- **Environment Configuration**: Secure configuration using environment variables
- **Multiple Token Types**: Support for regular users and admin users
- **Easy API Testing**: Generated ID tokens ready for immediate API endpoint testing

## 📋 Prerequisites

- Node.js 16+ installed
- Firebase project with Authentication enabled
- Firebase Admin SDK credentials file

## 🛠️ Setup

### 1. Install Dependencies

```bash
npm install
```

### 2. Configure Environment Variables

```bash
# Copy the example environment file
npm run setup
# or manually:
cp env.example .env
```

### 3. Edit `.env` file with your Firebase configuration

```env
# Get these from Firebase Console -> Project Settings -> General tab

FIREBASE_API_KEY=your_web_api_key_here
FIREBASE_AUTH_DOMAIN=your_project_id.firebaseapp.com
FIREBASE_PROJECT_ID=your_project_id
FIREBASE_STORAGE_BUCKET=your_project_id.appspot.com
FIREBASE_MESSAGING_SENDER_ID=123456789
FIREBASE_APP_ID=1:123456789:web:abcdef123456

# Path to your Firebase Admin SDK credentials file
FIREBASE_CREDENTIALS_FILE=../firebase-credentials.json
```

### 4. Get Firebase Configuration Values

1. Go to [Firebase Console](https://console.firebase.google.com)
2. Select your project
3. Click ⚙️ Settings → Project Settings
4. In the **General** tab, scroll to "Your apps"
5. Click the web app icon `</>` or create a new web app
6. Copy the configuration values to your `.env` file

### 5. Download Firebase Admin SDK Credentials

1. In Firebase Console → Project Settings
2. Click **Service accounts** tab
3. Click **Generate new private key**
4. Save the JSON file as `firebase-credentials.json` in your project root (or update the path in `.env`)

## 🎯 Usage

### Method 1: Command Line ID Token Generation

Generate ID tokens ready for immediate use:

```bash
# Generate ID tokens using Node.js
npm run generate-token
# or
npm start
```

This will output:
- A ready-to-use ID token for regular users
- A ready-to-use ID token for admin users  
- Example curl commands showing how to use the tokens

### Method 2: Web Interface

Start the web server to use the browser-based token generator:

```bash
# Start the web server
npm run serve
```

Then open `http://localhost:3000` in your browser:

1. **Sign Up**: Create a new test user
2. **Sign In**: Authenticate with existing credentials
3. **Copy Token**: Get the ID token for API testing
4. **Refresh**: Generate a new token when needed

## 🧪 Testing Your API

The generated ID tokens are ready for immediate use with your API endpoints:

```bash
# Set your token as an environment variable
export ID_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

# Test protected endpoints
curl -X GET "http://localhost:8080/api/v1/user/profile" \
  -H "Authorization: Bearer $ID_TOKEN" \
  -H "Content-Type: application/json"

# Test admin endpoints (with admin token)
curl -X GET "http://localhost:8080/api/v1/admin/users" \
  -H "Authorization: Bearer $ID_TOKEN" \
  -H "Content-Type: application/json"
```

## 📝 Scripts

| Command | Description |
|---------|-------------|
| `npm start` | Generate ID tokens via command line |
| `npm run generate-token` | Same as start |
| `npm run serve` | Start web server for browser interface |
| `npm run setup` | Copy environment template |

## 🔧 Configuration

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `FIREBASE_API_KEY` | Firebase Web API Key | `AIzaSyC...` |
| `FIREBASE_AUTH_DOMAIN` | Firebase Auth Domain | `myproject.firebaseapp.com` |
| `FIREBASE_PROJECT_ID` | Firebase Project ID | `myproject-12345` |
| `FIREBASE_STORAGE_BUCKET` | Firebase Storage Bucket | `myproject.appspot.com` |
| `FIREBASE_MESSAGING_SENDER_ID` | Firebase Messaging Sender ID | `123456789` |
| `FIREBASE_APP_ID` | Firebase App ID | `1:123:web:abc123` |
| `FIREBASE_CREDENTIALS_FILE` | Path to credentials JSON | `../firebase-credentials.json` |
| `PORT` | Web server port (optional) | `3000` |

### Custom Claims

You can modify the custom claims in `get-firebase-token.js`:

```javascript
const additionalClaims = {
  admin: true,           // Admin privileges
  email: 'user@test.com', // User email
  name: 'Test User',     // Display name
  role: 'editor',        // Custom role
  permissions: ['read', 'write'] // Custom permissions
};
```

## 🔒 Security Notes

- **Never commit** `.env` files or credentials to Git
- **Rotate credentials** regularly in production
- **Use different** Firebase projects for development and production
- **Limit token lifetime** by refreshing tokens frequently
- **Remove debug endpoints** (`/api/config`) in production

## 🐛 Troubleshooting

### Common Issues

**"Missing required environment variables"**
- Ensure `.env` file exists and contains all required variables
- Copy from `env.example` and fill in actual values

**"Firebase credentials file not found"**
- Check the path in `FIREBASE_CREDENTIALS_FILE`
- Ensure the JSON file exists and is readable

**"Error verifying ID token"**
- Token might be expired (tokens expire after 1 hour)
- Refresh the token using the web interface
- Ensure your server is using the same Firebase project

**"CORS errors in browser"**
- Ensure Firebase Authentication is enabled
- Check Firebase Console → Authentication → Settings → Authorized domains

### Debug Mode

Start the server and visit these endpoints for debugging:

- `http://localhost:3000/health` - Server health check
- `http://localhost:3000/api/config` - View Firebase configuration

## 📚 Resources

- [Firebase Authentication](https://firebase.google.com/docs/auth)
- [Firebase Admin SDK](https://firebase.google.com/docs/admin/setup)
- [Custom Token Creation](https://firebase.google.com/docs/auth/admin/create-custom-tokens)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License. 