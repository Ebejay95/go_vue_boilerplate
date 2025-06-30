const express = require('express');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

const app = express();
const PORT = process.env.PORT || process.env.FRONTEND_PORT || 3000;
const GRPC_SERVER_URL = process.env.GRPC_SERVER_URL || 'localhost:50051';

console.log(`Frontend server starting on port ${PORT}`);
console.log(`Connecting to gRPC server at: ${GRPC_SERVER_URL}`);

// Load protobuf
const PROTO_PATH = path.join(__dirname, '../proto/user.proto');
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

// WICHTIG: user.UserService statt nur userProto.UserService
const userProto = grpc.loadPackageDefinition(packageDefinition);
const client = new userProto.user.UserService(GRPC_SERVER_URL, grpc.credentials.createInsecure());

app.use(express.json());

// CORS middleware for Vue development
app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept, Authorization');

  if (req.method === 'OPTIONS') {
    res.sendStatus(200);
  } else {
    next();
  }
});

// Serve Vue.js build files in production
if (process.env.NODE_ENV === 'production') {
  app.use(express.static(path.join(__dirname, 'dist')));
}

// Test der gRPC-Verbindung beim Start
client.ListUsers({}, (error, response) => {
  if (error) {
    console.error('❌ gRPC connection failed:', error.message);
  } else {
    console.log('✅ gRPC connection successful');
    console.log(`Found ${response.users ? response.users.length : 0} users`);
  }
});

// API Routes
app.get('/api/users', (req, res) => {
  client.ListUsers({}, (error, response) => {
    if (error) {
      console.error('Error:', error);
      res.status(500).json({ error: error.message });
      return;
    }
    res.json(response.users || []);
  });
});

app.get('/api/users/:id', (req, res) => {
  const id = parseInt(req.params.id);
  client.GetUser({ id }, (error, response) => {
    if (error) {
      console.error('Error:', error);
      res.status(404).json({ error: error.message });
      return;
    }
    res.json(response.user);
  });
});

app.post('/api/users', (req, res) => {
  const { name, email, age } = req.body;
  client.CreateUser({ name, email, age: parseInt(age) }, (error, response) => {
    if (error) {
      console.error('Error:', error);
      res.status(500).json({ error: error.message });
      return;
    }
    res.json(response.user);
  });
});

// Health check endpoint
app.get('/api/health', (req, res) => {
  res.json({
    status: 'healthy',
    grpcServer: GRPC_SERVER_URL,
    port: PORT
  });
});

// Serve Vue app for all other routes in production
if (process.env.NODE_ENV === 'production') {
  app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, 'dist/index.html'));
  });
}

app.listen(PORT, '0.0.0.0', () => {
  console.log(`✅ Frontend server listening on port ${PORT}`);
  if (process.env.NODE_ENV !== 'production') {
    console.log(`Vue development server should be running on port 8080`);
  }
});
