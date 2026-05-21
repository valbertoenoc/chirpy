# Chirpy

A chirpy social API built with Go — a guided project for learning backend development.

## Stack

- **Go 1.25** — standard library HTTP server (`net/http`)
- **PostgreSQL** — database with `sqlc` for type-safe query generation
- **Swaggo** — auto-generated Swagger/OpenAPI docs
- **JWT + Argon2** — authentication via access/refresh tokens
- **godotenv** — environment configuration

## Quick Start

1. Copy `.env.sample` to `.env` and fill in `DB_URL`, `SECRET_KEY`, and `POLKA_KEY`.
2. Run migrations:
   ```sh
   ./migrate_up.sh
   ```
3. Start the server:
   ```sh
   go run .
   ```

The API serves on `:8080`. Swagger UI is at `/docs/`.

## Features

### Users
- **POST** `/api/users` — create account (email + password, hashed with Argon2)
- **PUT** `/api/users` — update email/password (requires JWT)
- **POST** `/api/login` — authenticate, returns access + refresh tokens
- **POST** `/api/refresh` — rotate refresh token
- **POST** `/api/revoke` — invalidate refresh token

### Chirps
- **POST** `/api/chirps` — create a chirp (140 char max, profanity-redacted, requires JWT)
- **GET** `/api/chirps` — list chirps, optional `?author_id=` filter and `?sort=asc|desc`
- **GET** `/api/chirps/{id}` — get a single chirp
- **DELETE** `/api/chirps/{id}` — delete own chirp (requires JWT)

### Admin
- **GET** `/admin/metrics` — fileserver hit counter
- **POST** `/admin/reset` — reset database (dev only)

### Webhooks
- **POST** `/api/polka/webhooks` — upgrade user to Chirpy Red (requires API key)

## Project Layout

```
├── main.go                  # server setup, routes, swagger annotations
├── handler_*.go             # HTTP handlers
├── internal/
│   ├── auth/                # JWT, Argon2 password hashing, bearer token parsing
│   ├── database/            # sqlc-generated types and queries
│   └── utils/               # profanity filter
├── sql/
│   ├── schema/              # SQL migration files
│   ├── queries/             # Named SQL queries for sqlc
│   └── sqlc.yaml            # sqlc config
├── docs/                    # Auto-generated swagger docs
└── assets/                  # Static files
```

## Regenerating Docs

```sh
swag init -g main.go --output docs
```
