# Gosift Backend

Backend API Gateway for Gosift built with Go and Gin. This repo is focused on the backend; the UI lives in a separate frontend repo.

## Prerequisites

- Go 1.25+
- Localstack
- awscli

## Run locally

Use your own '.env' file

```bash
cp .env.example .env
```

Then run the server:

```bash
docker compose up
go run ./cmd/server
```

## API surface (starter)

| Endpoint        | Description      |
| --------------- | ---------------- |
| `/api/v1`       |                  |
| `/api/v1/auth`  | register/login.  |
| `/api/v1/users` | CRUD operations. |

<!-- - `GET /ping`: health check returning `{"message":"pong"}` -->
