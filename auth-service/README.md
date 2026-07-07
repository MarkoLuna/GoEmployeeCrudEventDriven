# Auth Service

Centralized authentication and user management service for the GoEmployeeCrudEventDriven platform.

| Property | Value |
|---|---|
| Port | `8082` |
| Module | `github.com/MarkoLuna/AuthService` |
| Role | Token generation, user CRUD, credential validation |

## Architecture

```
POST /oauth/token  ──> validates client + user ──> returns JWT
GET  /oauth/userinfo ──> extracts JWT claims
/api/user/*        ──> CRUD operations on users
```

Auth-service delegates to **Keycloak** when `OAUTH_PROVIDER=keycloak`, or uses its own local OAuth (HS512 JWT) with PostgreSQL-backed users when `OAUTH_PROVIDER=local`.

## API Endpoints

### Auth

| Method | Path | Description |
|---|---|---|
| `POST` | `/oauth/token` | Generate JWT (Basic auth client credentials + form user/password) |
| `GET` | `/oauth/userinfo` | Get claims from bearer token |

### User Management

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/user` | Create user |
| `GET` | `/api/user` | List all users |
| `GET` | `/api/user/:id` | Get user by ID |
| `PUT` | `/api/user/:id` | Update user |
| `DELETE` | `/api/user/:id` | Soft-delete user |

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | `8082` | HTTP listen port |
| `SERVER_HOST` | `0.0.0.0` | HTTP listen host |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_NAME` | `auth_db` | Database name |
| `DB_USER` | `auth_user` | Database user |
| `DB_PASSWORD` | `authpw` | Database password |
| `OAUTH_ENABLED` | `false` | Enable JWT middleware |
| `OAUTH_PROVIDER` | `local` | `local` or `keycloak` |
| `OAUTH_SIGNING_KEY` | `00000000` | HMAC signing key (local mode) |
| `KEYCLOAK_AUTH_SERVER_URL` | `http://localhost:8084` | Keycloak server URL (keycloak mode) |
| `KEYCLOAK_REALM` | `dev` | Keycloak realm (keycloak mode) |

## Running Locally

```bash
cd auth-service
go run cmd/api-server/main.go
```

## Example Curls

```bash
# Login
curl -X POST http://localhost:8082/oauth/token \
  -H "Authorization: Basic $(echo -n 'client:password' | base64)" \
  -d "username=user&password=secret"

# Get user info
curl http://localhost:8082/oauth/userinfo \
  -H "Authorization: Bearer <token>"

# Create user
curl -X POST http://localhost:8082/api/user/ \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","password":"secret","email":"user@example.com","firstName":"New","lastName":"User"}'

# List users
curl http://localhost:8082/api/user/
```

## Database

Initialization script: `resources/init.sql`

Creates `users` and `clients` tables. Passwords are bcrypt-hashed.
