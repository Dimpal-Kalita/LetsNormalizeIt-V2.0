require('dotenv').config();
const admin = require('firebase-admin');
const fs = require('fs');
const path = require('path');
const fetch = require('node-fetch');

// Validate required environment variables
const requiredEnvVars = [
  'FIREBASE_PROJECT_ID',
  'FIREBASE_API_KEY',
  'FIREBASE_CREDENTIALS_FILE'
];

const missingVars = requiredEnvVars.filter(varName => !process.env[varName]);
if (missingVars.length > 0) {
  console.error('❌ Missing required environment variables:');
  missingVars.forEach(varName => console.error(`   - ${varName}`));
  console.error('\n💡 Please copy env.example to .env and fill in your Firebase configuration');
  process.exit(1);
}

// Initialize Firebase Admin
try {
  const credentialsPath = path.resolve(__dirname, process.env.FIREBASE_CREDENTIALS_FILE);
  
  if (!fs.existsSync(credentialsPath)) {
    console.error(`❌ Firebase credentials file not found: ${credentialsPath}`);
    console.error('💡 Please ensure your firebase-credentials.json file is in the correct location');
    process.exit(1);
  }

  const serviceAccount = require(credentialsPath);

  admin.initializeApp({
    credential: admin.credential.cert(serviceAccount),
    projectId: process.env.FIREBASE_PROJECT_ID
  });

  console.log('✅ Firebase Admin SDK initialized successfully');
} catch (error) {
  console.error('❌ Failed to initialize Firebase Admin SDK:', error.message);
  process.exit(1);
}

async function exchangeCustomTokenForIdToken(customToken) {
  try {
    const API_KEY = process.env.FIREBASE_API_KEY;
    const response = await fetch(`https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=${API_KEY}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        token: customToken,
        returnSecureToken: true
      })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data.idToken;
  } catch (error) {
    console.error('❌ Error exchanging custom token for ID token:', error);
    return null;
  }
}

async function createIdToken() {
  try {
    // Create a custom token first
    const uid = 'test-user-' + Date.now();
    
    // You can add custom claims here
    const additionalClaims = {
      admin: false, // Set to true for admin testing
      email: 'test@example.com',
      name: 'Test User'
    };

    // Create custom token using the correct Firebase Admin SDK method
    const customToken = await admin.auth().createCustomToken(uid, additionalClaims);
    
    // Exchange custom token for ID token
    const idToken = await exchangeCustomTokenForIdToken(customToken);
    
    if (idToken) {
      console.log('\n🔑 ID Token Generated:');
      console.log(idToken);
      console.log('\n📝 User ID:', uid);
      console.log('📋 Claims:', JSON.stringify(additionalClaims, null, 2));
    }
    
    return idToken;
  } catch (error) {
    console.error('❌ Error creating ID token:', error);
    return null;
  }
}

// Create admin ID token
async function createAdminIdToken() {
  try {
    const uid = 'admin-user-' + Date.now();
    
    const additionalClaims = {
      admin: true, // Admin claim
      email: 'admin@example.com',
      name: 'Admin User'
    };

    // Create custom token using the correct Firebase Admin SDK method
    const customToken = await admin.auth().createCustomToken(uid, additionalClaims);
    
    // Exchange custom token for ID token
    const idToken = await exchangeCustomTokenForIdToken(customToken);
    
    if (idToken) {
      console.log('\n👑 Admin ID Token Generated:');
      console.log(idToken);
      console.log('\n📝 User ID:', uid);
      console.log('📋 Claims:', JSON.stringify(additionalClaims, null, 2));
    }
    
    return idToken;
  } catch (error) {
    console.error('❌ Error creating admin ID token:', error);
    return null;
  }
}

async function main() {
  console.log('🔥 Firebase ID Token Generator');
  console.log('══════════════════════════════════════════════════════════════');
  console.log(`📦 Project ID: ${process.env.FIREBASE_PROJECT_ID}`);
  console.log(`🔑 Credentials: ${process.env.FIREBASE_CREDENTIALS_FILE}`);
  console.log('══════════════════════════════════════════════════════════════');
  
  console.log('\n1️⃣  Creating regular user ID token...');
  const regularIdToken = await createIdToken();
  
  console.log('\n2️⃣  Creating admin user ID token...');
  const adminIdToken = await createAdminIdToken();
  
  console.log('\n✨ ID token generation completed!');
  console.log('\n📚 Next steps:');
  console.log('   1. Use the ID tokens above to test your authentication endpoints');
  console.log('   2. Run "npm run serve" to start the HTML token generator');
  
  if (regularIdToken) {
    console.log('\n🔧 Example usage with curl:');
    console.log('─'.repeat(80));
    console.log(`curl -X GET "http://your-api-endpoint" \\`);
    console.log(`  -H "Authorization: Bearer ${regularIdToken.substring(0, 50)}..."`);
    console.log('─'.repeat(80));
  }
}

// Handle graceful shutdown
process.on('SIGINT', () => {
  console.log('\n👋 Shutting down gracefully...');
  process.exit(0);
});

main().catch(console.error); 