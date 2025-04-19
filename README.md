# Flight Price Service

This service allows users to search for flights using multiple providers, aggregating the cheapest and fastest options available.

## ✅ Features Implemented

- 🔍 **Flight Search Aggregation** across multiple providers
- 📡 **Amadeus API Integration** (OAuth2 and flight offer endpoints)
- 🛡️ **JWT Authentication** support (with token generation endpoint)
- 🌐 **REST API** using `mux.Router`
- 💾 **Redis Cache Integration** to store recent search results (default TTL: 30s)
- ⚙️ **.env Configuration** support with `.env.example`
- 🧪 **Unit tests** and provider mocks
- 🧠 **Concurrency**: Provider calls are done concurrently for faster aggregation

## 🔧 Setup

### Prerequisites
- Go 1.23+
- Docker
- [golangci-lint](https://golangci-lint.run/) for linting

### Run Locally

```bash
git clone https://github.com/fehepe/flight-price-service.git
cd flight-price-service
cp .env.example .env
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
