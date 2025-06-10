require('dotenv').config();
const express = require('express');
const fs = require('fs');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 3000;

// Validate required environment variables
const requiredEnvVars = [
  'FIREBASE_API_KEY',
  'FIREBASE_AUTH_DOMAIN',
  'FIREBASE_PROJECT_ID',
  'FIREBASE_STORAGE_BUCKET',
  'FIREBASE_MESSAGING_SENDER_ID',
  'FIREBASE_APP_ID'
];

const missingVars = requiredEnvVars.filter(varName => !process.env[varName]);
if (missingVars.length > 0) {
  console.error('❌ Missing required environment variables:');
  missingVars.forEach(varName => console.error(`   - ${varName}`));
  console.error('\n💡 Please copy env.example to .env and fill in your Firebase configuration');
  process.exit(1);
}

// Firebase configuration from environment variables
const firebaseConfig = {
  apiKey: process.env.FIREBASE_API_KEY,
  authDomain: process.env.FIREBASE_AUTH_DOMAIN,
  projectId: process.env.FIREBASE_PROJECT_ID,
  storageBucket: process.env.FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.FIREBASE_APP_ID
};

// Serve the HTML file with injected Firebase config
app.get('/', (req, res) => {
  try {
    // Read the HTML template
    const htmlPath = path.join(__dirname, 'test-auth.html');
    let html = fs.readFileSync(htmlPath, 'utf8');
    
    // Replace the placeholder Firebase config with actual config
    const configString = JSON.stringify(firebaseConfig, null, 12);
    html = html.replace(
      /const firebaseConfig = \{[\s\S]*?\};/,
      `const firebaseConfig = ${configString};`
    );
    
    res.send(html);
  } catch (error) {
    console.error('❌ Error serving HTML file:', error);
    res.status(500).send('Internal Server Error');
  }
});

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({ 
    status: 'OK', 
    timestamp: new Date().toISOString(),
    firebase: {
      projectId: firebaseConfig.projectId,
      authDomain: firebaseConfig.authDomain
    }
  });
});

// API endpoint to get Firebase config (for debugging)
app.get('/api/config', (req, res) => {
  res.json({
    firebase: firebaseConfig,
    note: 'This endpoint is for debugging. Remove in production.'
  });
});

// Static files
app.use(express.static(__dirname));

// Start server
app.listen(PORT, () => {
  console.log('🔥 Firebase Token Creator Server');
  console.log('═══════════════════════════════════════════════════');
  console.log(`🌐 Server running at: http://localhost:${PORT}`);
  console.log(`📦 Firebase Project: ${firebaseConfig.projectId}`);
  console.log(`🏠 Auth Domain: ${firebaseConfig.authDomain}`);
  console.log('═══════════════════════════════════════════════════');
  console.log('\n📚 Available endpoints:');
  console.log(`   • http://localhost:${PORT}/          - Token generator UI`);
  console.log(`   • http://localhost:${PORT}/health    - Health check`);
  console.log(`   • http://localhost:${PORT}/api/config - Firebase config (debug)`);
  console.log('\n💡 Open the first URL in your browser to generate tokens');
});

// Handle graceful shutdown
process.on('SIGINT', () => {
  console.log('\n👋 Shutting down server...');
  process.exit(0);
}); 