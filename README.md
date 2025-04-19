# Flight Price Service

This service allows users to search for flights using multiple providers, aggregating the cheapest and fastest options available.

## ✅ Features Implemented

- 🔍 **Flight Search Aggregation** across multiple providers
- 📡 **Amadeus API Integration** (OAuth2 and flight offer endpoints)
- 🛡️ **JWT Authentication** support (with token generation endpoint)
- 🌐 **REST API** using `mux.Router`
- ⚙️ **.env Configuration** support with `.env.example`
- 🧪 **Unit tests** and provider mocks
- 🧠 **Concurrency**: Provider calls are done concurrently for faster aggregation
- 📊 **Smart Response**: 
  - Cheapest flight
  - Fastest flight
  - Grouped provider offers

## 🔧 Setup

### Prerequisites
- Go 1.23+
- Docker (optional for containerized development)
- [golangci-lint](https://golangci-lint.run/) for linting

### Clone and Run
```bash
git clone https://github.com/fehepe/flight-price-service.git
cd flight-price-service
cp .env.example .env
```

Edit `.env` and set values for:

```env
PORT=3000
JWT_SECRET=your_secret

MAX_FLIGHT_RESULTS_PER_CLIENT=10

AMADEUS_API_KEY=your_key
AMADEUS_API_SECRET=your_secret
AMADEUS_BASE_URL=https://test.api.amadeus.com
```

## 🚀 Run Locally

```bash
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
| `adults`        | int     | ❌       | Number of adult travelers (default: `1`)     |
| `non_stop`      | bool    | ❌       | Filter only non-stop flights (default: `false`) |

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
├── cmd/flight-service          # Main entrypoint
├── internal/
│   ├── config/                 
│   ├── middleware/       
│   ├── handlers/
│   ├── providers/
│   └── services/flight/
├── pkg
├── .env.example
├── Dockerfile
└── README.md
```

---
## 📦 Build for Production
```bash
docker build -t flight-service .
```
