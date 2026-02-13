# Environment Setup Guide

This guide explains how to set up environment variables for the HEMA Lessons API.

## Quick Start

1. **Copy the example environment file:**
   ```bash
   cp env.example .env
   ```

2. **Edit the `.env` file with your values:**
   ```bash
   nano .env  # or your preferred editor
   ```

3. **Never commit the `.env` file** (it's already in `.gitignore`)

4. **For Docker development, make sure environment variables are available:**
   ```bash
   # Either set them in your shell:
   export DATABASE_PASSWORD=your_secure_password
   export TEST_DATABASE_PASSWORD=your_secure_password

   # Or they will be read from your .env file automatically
   ```

## Environment Variables

### Server Configuration
- `SERVER_ADDR`: Server address and port (default: `:8080`)
- `SERVER_READ_HEADER_TIMEOUT`: HTTP read header timeout in seconds (default: `5`)

### Database Configuration
- `DATABASE_HOST`: PostgreSQL host (default: `localhost`)
- `DATABASE_PORT`: PostgreSQL port (default: `5432`)
- `DATABASE_USER`: PostgreSQL username
- `DATABASE_PASSWORD`: PostgreSQL password (**REQUIRED - never commit this!**)
- `DATABASE_DBNAME`: PostgreSQL database name (default: `hema_lessons`)
- `DATABASE_SSLMODE`: SSL mode for database connection (default: `disable`)

### Redis Configuration
- `REDIS_URL`: Redis connection URL (default: `redis://localhost:6379/0`)

### Application Configuration
- `APP_ENVIRONMENT`: Environment type (`development`, `production`, `staging`)

## Docker Development

For local Docker development, environment variables are set in `docker-compose.yml`:

```yaml
environment:
  SERVER_ADDR: ":8080"
  DATABASE_HOST: "db"
  DATABASE_PORT: "5432"
  DATABASE_USER: "hema"
  DATABASE_PASSWORD: "hema"
  DATABASE_DBNAME: "hema_lessons"
  DATABASE_SSLMODE: "disable"
  REDIS_URL: "redis://cache:6379/0"
  APP_ENVIRONMENT: "development"
```

## Production Deployment

For production, set environment variables in your deployment platform:

### Railway
```bash
railway variables set SERVER_ADDR=":8080"
railway variables set DATABASE_HOST="${{RAILWAY_POSTGRESQL_HOST}}"
railway variables set DATABASE_PORT="${{RAILWAY_POSTGRESQL_PORT}}"
railway variables set DATABASE_USER="${{RAILWAY_POSTGRESQL_USER}}"
railway variables set DATABASE_PASSWORD="${{RAILWAY_POSTGRESQL_PASSWORD}}"
railway variables set DATABASE_DBNAME="${{RAILWAY_POSTGRESQL_DATABASE}}"
railway variables set DATABASE_SSLMODE="require"
railway variables set REDIS_URL="${{RAILWAY_REDIS_URL}}"
railway variables set APP_ENVIRONMENT="production"
```

### Heroku
```bash
heroku config:set SERVER_ADDR=":8080"
heroku config:set DATABASE_URL="your_postgres_url"
heroku config:set REDIS_URL="your_redis_url"
heroku config:set APP_ENVIRONMENT="production"
```

### AWS/GCP/Azure
Use their respective secret management and environment variable services.

## Security Notes

- **Never commit sensitive data** to version control
- Use strong, unique passwords for database access
- Rotate credentials regularly in production
- Use SSL/TLS in production (`DATABASE_SSLMODE=require`)
- Consider using managed database services for production

## Testing

For running tests, you **must** set test-specific environment variables. The test utilities will fail if these are not provided (no default credentials for security):

```bash
# Required environment variables for tests
export TEST_DATABASE_HOST=db
export TEST_DATABASE_PORT=5432
export TEST_DATABASE_USER=hema
export TEST_DATABASE_PASSWORD=your_secure_test_password
export TEST_DATABASE_DBNAME=hema_lessons
export TEST_DATABASE_SSLMODE=disable
export TEST_APP_ENVIRONMENT=testing

# Then run tests
./scripts/run-tests.sh
```

**Security Note:** Test credentials should be different from production credentials, and the test script will validate that all required variables are set before running tests.