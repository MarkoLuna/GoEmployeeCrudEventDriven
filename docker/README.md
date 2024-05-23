# Keycloak Embedded in a Spring Boot Application
Keycloak is an open-source Identity and Access Management solution administered by RedHat and developed in Java by JBoss.


### Run authorization-server and Employee Service
```bash
docker compose up -d --build --remove-orphans
```

### Stop
```bash
docker compose down
```

### Generate tokens

```bash
curl --location 'http://localhost:8080/realms/master/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--header 'Cookie: AUTH_SESSION_ID_LEGACY=de120b71-c809-45e4-a67f-bbcf0a87bfdb' \
--data-urlencode 'client_id=master-client' \
--data-urlencode 'client_secret=z6vxpf3uzvJLlsErs9oufAyolCYFvEos' \
--data-urlencode 'username=marcosluna' \
--data-urlencode 'password=marco94' \
--data-urlencode 'grant_type=password'
```


```bash
curl --location 'http://localhost:8080/realms/dev/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_id=newClient' \
--data-urlencode 'client_secret=newClientSecret' \
--data-urlencode 'username=john@test.com' \
--data-urlencode 'password=123' \
--data-urlencode 'grant_type=password'
```

```bash
curl --location 'http://localhost:8080/realms/dev/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_id=newClient' \
--data-urlencode 'client_secret=newClientSecret' \
--data-urlencode 'username=mike@other.com' \
--data-urlencode 'password=pass' \
--data-urlencode 'grant_type=password'
```

## Users
- john@test.com / 123
- mike@other.com / pass
- marcosluna / marco94