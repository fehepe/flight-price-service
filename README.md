# Flight Price Service

This service allows users to search for flights using multiple providers, aggregating the cheapest and fastest options available.

## ✅ Features Implemented

- 🔍 **Flight Search Aggregation** across multiple providers
- 📡 **AmadeusAPI, SerAPI and PriceLine Integrations** (OAuth2 and flight offer endpoints)
- 🛡️ **JWT Authentication** support (with token generation endpoint)
- 🌐 **REST API** using `mux.Router`
- 💾 **Redis Cache Integration** to store recent search results (default TTL: 30s)
- 🧪 **Unit tests** and provider mocks
- 🧠 **Concurrency**: Provider calls are done concurrently for faster aggregation
- 🔒 **Encrypted** credentials via git-crypt

## 🔧 Setup

### Prerequisites
- Go 1.23+
- Docker
- **git‑crypt** (for decrypting `credentials.json`)

## Installing git‑crypt

### macOS (Homebrew)
```bash
brew install git-crypt gnupg
```

### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install -y git-crypt gnupg
```

### Fedora/CentOS
```bash
sudo yum install -y git-crypt gnupg2
```

## Run Locally

1. **Clone the repo**
   ```bash
   git clone https://github.com/fehepe/flight-price-service.git
   
   cd flight-price-service
   ```
2. **Copy environment template**
   ```bash
   cp .env.example .env
   ```
3. **Get added as a git‑crypt recipient**  
   Contact the maintainer (`fehepe11@gmail.com`) so I can send you the gpg file.
4. **Unlock encrypted credentials**
   ```bash
   git-crypt unlock
   ```
5. **Start the service**
   ```bash
   docker-compose up --build
   ```

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

## 📘 API Endpoints

### Health Check
```http
GET /health
```

### Search Flights
```http
GET /flights/search?origin=SYD&destination=BKK&departure_date=2025-05-02
Authorization: Bearer <your_token>
```
**Query Parameters:**

| Name            | Type    | Required | Description                                 |
|-----------------|---------|----------|---------------------------------------------|
| `origin`        | string  | ✅       | IATA code of departure airport (e.g., `SYD`) |
| `destination`   | string  | ✅       | IATA code of arrival airport (e.g., `BKK`)   |
| `departure_date`| string  | ✅       | Departure date in `YYYY-MM-DD` format        |

Returns:
```json
{
  "cheapest": { ... },
  "fastest": { ... },
  "providers": {
    "ProviderA": [ ... ],
    "ProviderB": [ ... ],
    "ProviderC": [ ... ]
  }
}
```

## 📂 Structure

```
flight-price-service/
├── cmd/flight-service
├── internal/
│   ├── cache/  
│   ├── config/                 
│   ├── middleware/       
│   ├── handlers/
│   ├── providers/
│   └── services/
├── pkg
├── .env.example
├── Dockerfile
└── README.md
```
