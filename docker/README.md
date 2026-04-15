# Docker Infrastructure & Authentication Stack

This directory contains the Docker Compose configurations to spin up the complete development environment, including microservices, databases, messaging brokers, and the Identity and Access Management (IAM) server.

## 🚀 Quick Start

Ensure you have [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed.

```bash
# Start the full stack in detached mode
docker compose up -d --build --remove-orphans

# Stop the stack
docker compose down
```

> [!TIP]
> The setup uses a `.env` file in this directory to configure versions and default credentials. By default, it uses `keycloak-compose.yml`.

---

## 🛠 Service Map

The stack exposes several services to your host machine for development and monitoring:

| Service | Container Name | Port (Host) | Description |
|---|---|---|---|
| **Employee Service** | `employee-service` | `8080` | Go CRUD API (REST) |
| **Employee Consumer** | `employee-consumer` | `8081` | Go Kafka Consumer / Internal API |
| **Keycloak** | `ws-keycloak-x` | `8082` | IAM Server (OIDC/SAML) |
| **Kafka** | `kafka` | `9092` | Message Broker |
| **PostgreSQL (App)** | `postgres_db` | `5433` | Database for Employee Services |
| **PostgreSQL (IAM)** | `ws-keycloak-x-pg` | `5432` | Database for Keycloak |
| **Zookeeper** | `zookeeper` | `2181` | Kafka Orchestration |

---

## 🔑 Authentication (Keycloak)

The environment comes with a pre-configured Keycloak instance.

### Admin Console
- **URL**: [http://localhost:8082](http://localhost:8082)
- **User**: `dev`
- **Password**: `dev`

### Generating JWT Tokens

You can use the following `curl` commands to obtain access tokens for testing the protected endpoints in the Employee Service.

#### 🌍 Master Realm (Admin)
```bash
curl --location 'http://localhost:8082/realms/master/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_id=master-client' \
--data-urlencode 'client_secret=z6vxpf3uzvJLlsErs9oufAyolCYFvEos' \
--data-urlencode 'username=marcosluna' \
--data-urlencode 'password=marco94' \
--data-urlencode 'grant_type=password'
```

#### 🛠 Dev Realm (Standard Users)
The `dev` realm is imported automatically on startup from `./keycloak/import/dev-realm`.

**User: John Doe**
```bash
curl --location 'http://localhost:8082/realms/dev/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_id=newClient' \
--data-urlencode 'client_secret=newClientSecret' \
--data-urlencode 'username=john@test.com' \
--data-urlencode 'password=123' \
--data-urlencode 'grant_type=password'
```

**User: Mike Smith**
```bash
curl --location 'http://localhost:8082/realms/dev/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_id=newClient' \
--data-urlencode 'client_secret=newClientSecret' \
--data-urlencode 'username=mike@other.com' \
--data-urlencode 'password=pass' \
--data-urlencode 'grant_type=password'
```

---

## 🔍 Operations & Troubleshooting

### Viewing Logs
```bash
# View all logs
docker compose logs -f

# View logs for a specific service
docker compose logs -f employee-service
```

### Resetting the Environment
If you need to wipe all data and start fresh (including databases):
```bash
docker compose down -v
docker compose up -d --build
```

### Common Issues
- **Port 5432 Conflict**: Ensure no local PostgreSQL is running on port 5432, as it's used by the Keycloak database.
- **Kafka Connectivity**: If services fail to connect to Kafka, ensure the `kafka` container is healthy (`docker ps`).