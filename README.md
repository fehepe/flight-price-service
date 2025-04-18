# Flight Price Service

This Go microservice fetches flight prices from multiple providers concurrently and returns structured comparison data. It includes authentication, caching, and containerization for easy deployment.

---

## 📁 Project Structure

```
flight-price-service/
├── cmd/                  # Entry point (main.go)
│   └── flight-service/
├── internal/             # Application code
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # Auth and logging middleware
│   ├── server/           # Server setup and router
│   └── config/           # Environment loading
├── pkg/                  # Shared types and utilities
│   ├── models/           # API models
│   └── utils/            # JSON response helpers
├── .env.example          # Sample environment file
├── .env                  # Local environment config
├── .gitignore
├── Dockerfile
├── docker-compose.yml (optional)
├── go.mod / go.sum
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites
- Go 1.23+
- Docker (optional for containerized development)
- [golangci-lint](https://golangci-lint.run/) for linting

### Clone and Run
```bash
git clone https://github.com/fehepe/flight-price-service.git
cd flight-price-service

# Set up env
cp .env.example .env

# Run locally
go run ./cmd/flight-service
```

### Using Docker
```bash
docker build -t flight-service .
docker run --rm -p 3000:3000 --env-file .env flight-service
```

Or with Compose:
```bash
docker-compose up --build
```

---

## 🔐 JWT Authentication

### Token Endpoint
```
POST /auth/token
Content-Type: application/json
```
```json
{
  "username": "user",
  "password": "pass"
}
```
Returns a JWT token to use in authenticated endpoints.

### Protected Endpoint Example
```
GET /flights/search
Authorization: Bearer <your_token>
```

---

## ⚙️ Environment Variables

Here's a sample `.env.example`:
```env
PORT=3000
LOG_LEVEL=info
READ_TIMEOUT=5
WRITE_TIMEOUT=10
IDLE_TIMEOUT=120

AUTH_USERNAME=user
AUTH_PASSWORD=pass

JWT_SECRET=your_super_secret_key_here
JWT_EXPIRY_HOURS=24
JWT_ISSUER=flight-service
```

---

## 🧪 Testing

To test manually, you can:
```bash
curl -X POST http://localhost:3000/auth/token \
     -H "Content-Type: application/json" \
     -d '{"username":"user","password":"pass"}'

# Then use the token in /flights/search
```

---

## 🧼 Linting

```bash
golangci-lint run ./...
```

---

## 📦 Build for Production
```bash
docker build -t flight-service .
```
