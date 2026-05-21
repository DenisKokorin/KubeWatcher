# Backend E2E test infrastructure

This folder contains helpers and end-to-end tests that interact with a running backend and its dependencies.

Overview
- `docker-compose.test.yml` — starts test services: Postgres, Redis, MinIO, Wiremock.
- `env/.env.test` — environment template for backend when running in test mode.
- `e2e` folder — Go tests and helpers. Tests expect `BACKEND_BASE_URL` (default http://localhost:7979) and a running backend configured with `env/.env.test`.

Quickstart
1. Start test services:
```bash
cd backend
docker-compose -f docker-compose.test.yml up -d
```

2. Create test database schema and run migrations (example using psql):
```bash
# adjust host/port if running elsewhere
psql "postgres://test:testpass@localhost:5432/kubewatcher_test" -f migrations/001_init_schema.sql
psql "postgres://test:testpass@localhost:5432/kubewatcher_test" -f migrations/002_add_role_column.sql
```

3. Create MinIO bucket (if not created):
```bash
mc alias set local http://localhost:9000 minioadmin minioadmin
mc mb local/kubewatcher-test
```

4. Start backend configured with `env/.env.test`:
```bash
cp env/.env.test env/.env
go run ./cmd -kubeconfig="" 
# or run binary with same env file
```

5. Run E2E tests (from backend folder):
```bash
export BACKEND_BASE_URL=http://localhost:7979
go test ./e2e -v
```

Notes
- Tests will create and modify data in the test database; use isolated DB.
- `ResetDatabase()` helper truncates public tables — call it between tests if you run tests individually.
- Wiremock runs inside the test compose network on port `8080`. The compose file maps container port `8080` to host port `8089`.
	- If you run the backend inside the same Docker Compose network (for example, by running the backend in a container), set `EXTERNAL_API_URL` to `http://wiremock:8080`.
	- If you run the backend on the host (common during local dev), set `EXTERNAL_API_URL` to `http://localhost:8089` (this is the default in `env/.env.test`).
	- Ensure you have Wiremock mappings for the endpoints used by tests. You can mount a `mappings` folder into the Wiremock container or use the Wiremock admin API to create mappings.
