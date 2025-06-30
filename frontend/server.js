const express = require('express');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 3000;
const GRPC_SERVER_URL = process.env.GRPC_SERVER_URL || 'localhost:50051';

// Load protobuf
const PROTO_PATH = path.join(__dirname, '../proto/user.proto');
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const userProto = grpc.loadPackageDefinition(packageDefinition).user;
const client = new userProto.UserService(GRPC_SERVER_URL, grpc.credentials.createInsecure());

app.use(express.json());
app.use(express.static('public'));

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

// Serve HTML
app.get('/', (req, res) => {
  res.send(`
    <!DOCTYPE html>
    <html lang="de">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>gRPC User Management</title>
        <style>
            body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
            .container { margin: 20px 0; }
            table { width: 100%; border-collapse: collapse; margin: 20px 0; }
            th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
            th { background-color: #f2f2f2; }
            input, button { margin: 5px; padding: 8px; }
            button { background-color: #4CAF50; color: white; border: none; cursor: pointer; }
            button:hover { background-color: #45a049; }
            .form-group { margin: 10px 0; }
            label { display: inline-block; width: 100px; }
        </style>
    </head>
    <body>
        <h1>gRPC User Management</h1>

        <div class="container">
            <h2>Neuen User erstellen</h2>
            <form id="userForm">
                <div class="form-group">
                    <label for="name">Name:</label>
                    <input type="text" id="name" required>
                </div>
                <div class="form-group">
                    <label for="email">Email:</label>
                    <input type="email" id="email" required>
                </div>
                <div class="form-group">
                    <label for="age">Alter:</label>
                    <input type="number" id="age" required>
                </div>
                <button type="submit">User erstellen</button>
            </form>
        </div>

        <div class="container">
            <h2>Alle Users</h2>
            <button onclick="loadUsers()">Users laden</button>
            <table id="usersTable">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>Email</th>
                        <th>Alter</th>
                    </tr>
                </thead>
                <tbody id="usersTableBody">
                </tbody>
            </table>
        </div>

        <script>
            async function loadUsers() {
                try {
                    const response = await fetch('/api/users');
                    const users = await response.json();

                    const tbody = document.getElementById('usersTableBody');
                    tbody.innerHTML = '';

                    users.forEach(user => {
                        const row = tbody.insertRow();
                        row.insertCell(0).textContent = user.id;
                        row.insertCell(1).textContent = user.name;
                        row.insertCell(2).textContent = user.email;
                        row.insertCell(3).textContent = user.age;
                    });
                } catch (error) {
                    console.error('Error loading users:', error);
                    alert('Fehler beim Laden der Users');
                }
            }

            document.getElementById('userForm').addEventListener('submit', async (e) => {
                e.preventDefault();

                const name = document.getElementById('name').value;
                const email = document.getElementById('email').value;
                const age = document.getElementById('age').value;

                try {
                    const response = await fetch('/api/users', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ name, email, age: parseInt(age) })
                    });

                    if (response.ok) {
                        document.getElementById('userForm').reset();
                        loadUsers();
                        alert('User erfolgreich erstellt!');
                    } else {
                        alert('Fehler beim Erstellen des Users');
                    }
                } catch (error) {
                    console.error('Error creating user:', error);
                    alert('Fehler beim Erstellen des Users');
                }
            });

            // Load users on page load
            loadUsers();
        </script>
    </body>
    </html>
  `);
});

app.listen(PORT, () => {
  console.log(`Frontend server running on port ${PORT}`);
  console.log(`Connecting to gRPC server at ${GRPC_SERVER_URL}`);
});
