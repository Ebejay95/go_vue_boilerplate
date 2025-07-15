const express = require('express');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

const app = express();
const PORT = process.env.PORT || process.env.FRONTEND_PORT;
const GRPC_SERVER_URL = process.env.GRPC_SERVER_URL;

console.log(`Frontend server starting on port ${PORT}`);
console.log(`Connecting to gRPC server at: ${GRPC_SERVER_URL}`);

const PROTO_PATH = path.join(__dirname, '../proto/user.proto');
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
	keepCase: true,
	longs: String,
	enums: String,
	defaults: true,
	oneofs: true,
});

const userProto = grpc.loadPackageDefinition(packageDefinition);
const client = new userProto.user.UserService(GRPC_SERVER_URL, grpc.credentials.createInsecure());

app.use(express.json());

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

if (process.env.NODE_ENV === 'production') {
	app.use(express.static(path.join(__dirname, 'dist')));
}

client.ListUsers({}, (error, response) => {
	if (error) {
		console.error('❌ gRPC connection failed:', error.message);
	} else {
		console.log('✅ gRPC connection successful');
		console.log(`Found ${response.users ? response.users.length : 0} users`);
	}
});

app.listen(PORT, '0.0.0.0', () => {
	console.log(`✅ Frontend server listening on port ${PORT}`);
});
