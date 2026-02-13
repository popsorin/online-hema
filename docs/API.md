# API Documentation

## Base URL

`http://localhost:8080` (configurable via `SERVER_ADDR` environment variable)

## Endpoints

### Health Check

**GET /healthz**

Returns the health status of the API and database connection.

```bash
curl http://localhost:8080/healthz
```

Response:

```json
{
  "status": "ok",
  "timestamp": "2026-01-18T12:34:56Z"
}
```

In development mode, the status includes the environment:

```json
{
  "status": "ok (development)",
  "timestamp": "2026-01-18T12:34:56Z"
}
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Protected endpoints require an `Authorization` header with a Bearer token.

### Register

**POST /api/auth/register**

Create a new user account.

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123", "full_name": "John Doe"}'
```

Request Body:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | User's email address |
| `password` | string | Yes | Password (minimum 8 characters) |
| `full_name` | string | Yes | User's full name |

Response (201 Created):

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe"
  }
}
```

Error Responses:
- **400 Bad Request** - Missing fields or password too short
- **409 Conflict** - Email already registered

### Login

**POST /api/auth/login**

Authenticate with email and password.

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

Request Body:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | User's email address |
| `password` | string | Yes | User's password |

Response (200 OK):

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe"
  }
}
```

Error Responses:
- **400 Bad Request** - Missing email or password
- **401 Unauthorized** - Invalid email or password

### Refresh Token

**POST /api/auth/refresh**

Exchange a refresh token for a new token pair.

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "eyJhbGciOiJIUzI1NiIs..."}'
```

Request Body:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `refresh_token` | string | Yes | Valid refresh token |

Response (200 OK):

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe"
  }
}
```

Error Responses:
- **400 Bad Request** - Missing refresh token
- **401 Unauthorized** - Invalid or expired refresh token

### Logout

**POST /api/auth/logout**

Logout the current user (client should discard tokens).

```bash
curl -X POST http://localhost:8080/api/auth/logout
```

Response (200 OK):

```json
{
  "message": "logged out successfully"
}
```

### Get Current User

**GET /api/auth/me**

Get the currently authenticated user's profile. Requires authentication.

```bash
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

Response (200 OK):

```json
{
  "id": 1,
  "email": "user@example.com",
  "full_name": "John Doe"
}
```

Error Responses:
- **401 Unauthorized** - Missing or invalid token

## Content Endpoints

### List Fighting Books

**GET /api/fighting-books**

Returns a paginated list of fighting books with their sword master names.

Query Parameters:
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number (minimum 1) |
| `page_size` | int | 20 | Items per page (1-100) |

```bash
curl "http://localhost:8080/api/fighting-books?page=1&page_size=2"
```

Response:

```json
{
  "data": [
    {
      "id": 1,
      "sword_master_id": 2,
      "title": "Fior di Battaglia",
      "description": "The Flower of Battle - a comprehensive medieval combat manual covering armed and unarmed combat",
      "publication_year": 1409,
      "cover_image_url": null,
      "created_at": "2026-01-18T10:00:00Z",
      "updated_at": "2026-01-18T10:00:00Z",
      "sword_master_name": "Fiore dei Liberi"
    },
    {
      "id": 2,
      "sword_master_id": 3,
      "title": "Fechtbuch",
      "description": "A detailed commentary on Liechtenauer's teachings with practical applications",
      "publication_year": 1440,
      "cover_image_url": null,
      "created_at": "2026-01-18T10:00:00Z",
      "updated_at": "2026-01-18T10:00:00Z",
      "sword_master_name": "Sigmund Ringeck"
    }
  ],
  "page": 1,
  "page_size": 2,
  "total_count": 3,
  "total_pages": 2
}
```

### Get Fighting Book by ID

**GET /api/fighting-books/{id}**

Returns a single fighting book by its ID.

Path Parameters:
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | int | The fighting book ID (must be > 0) |

```bash
curl http://localhost:8080/api/fighting-books/1
```

Response:

```json
{
  "id": 1,
  "sword_master_id": 2,
  "title": "Fior di Battaglia",
  "description": "The Flower of Battle - a comprehensive medieval combat manual covering armed and unarmed combat",
  "publication_year": 1409,
  "cover_image_url": null,
  "created_at": "2026-01-18T10:00:00Z",
  "updated_at": "2026-01-18T10:00:00Z",
  "sword_master_name": "Fiore dei Liberi"
}
```

Error Responses:

- **400 Bad Request** - Invalid ID format or ID <= 0
- **404 Not Found** - Fighting book with the given ID does not exist
- **500 Internal Server Error** - Database or server error

### List Chapters by Fighting Book

**GET /api/fighting-books/{id}/chapters**

Returns all chapters for a specific fighting book, ordered by chapter number.

Path Parameters:
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | int | The fighting book ID (must be > 0) |

```bash
curl http://localhost:8080/api/fighting-books/1/chapters
```

Response:

```json
[
  {
    "id": 1,
    "fighting_book_id": 1,
    "chapter_number": 1,
    "title": "Wrestling",
    "description": "Techniques for unarmed combat and grappling",
    "created_at": "2026-01-18T10:00:00Z",
    "updated_at": "2026-01-18T10:00:00Z"
  },
  {
    "id": 2,
    "fighting_book_id": 1,
    "chapter_number": 2,
    "title": "Dagger Combat",
    "description": "Fighting with the dagger in various situations",
    "created_at": "2026-01-18T10:00:00Z",
    "updated_at": "2026-01-18T10:00:00Z"
  },
  {
    "id": 3,
    "fighting_book_id": 1,
    "chapter_number": 3,
    "title": "Longsword",
    "description": "The art of fighting with the longsword",
    "created_at": "2026-01-18T10:00:00Z",
    "updated_at": "2026-01-18T10:00:00Z"
  },
  {
    "id": 4,
    "fighting_book_id": 1,
    "chapter_number": 4,
    "title": "Poleaxe",
    "description": "Combat techniques with the poleaxe",
    "created_at": "2026-01-18T10:00:00Z",
    "updated_at": "2026-01-18T10:00:00Z"
  }
]
```

Empty Response (when fighting book has no chapters):

```json
[]
```

Error Responses:

- **400 Bad Request** - Invalid ID format or ID <= 0
- **404 Not Found** - Fighting book with the given ID does not exist
- **500 Internal Server Error** - Database or server error

### List Techniques by Chapter

**GET /api/chapters/{id}/techniques**

Returns all techniques for a specific chapter, ordered by `order_in_chapter`.

Path Parameters:
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | int | The chapter ID (must be > 0) |

```bash
curl http://localhost:8080/api/chapters/3/techniques
```

Response:

```json
[
  {
    "id": 1,
    "chapter_id": 3,
    "name": "First Guard - Posta di Donna",
    "description": "The Woman's Guard - a high guard position",
    "instructions": "Hold the sword with the hilt near your right shoulder, point aimed at the opponent's face",
    "video_url": null,
    "thumbnail_url": null,
    "order_in_chapter": 1,
    "created_at": "2026-01-18T10:00:00Z",
    "updated_at": "2026-01-18T10:00:00Z"
  },
  {
    "id": 2,
    "chapter_id": 3,
    "name": "Zornhau",
    "description": "The Wrath Strike - a powerful diagonal cut",
    "instructions": "Strike diagonally from your right shoulder to the opponent's left side with full commitment",
    "video_url": null,
    "thumbnail_url": null,
    "order_in_chapter": 2,
    "created_at": "2026-01-18T10:00:00Z",
    "updated_at": "2026-01-18T10:00:00Z"
  },
  {
    "id": 3,
    "chapter_id": 3,
    "name": "Krumphau",
    "description": "The Crooked Strike - an off-line attack",
    "instructions": "Step offline and strike with the false edge, hands crossed",
    "video_url": null,
    "thumbnail_url": null,
    "order_in_chapter": 3,
    "created_at": "2026-01-18T10:00:00Z",
    "updated_at": "2026-01-18T10:00:00Z"
  }
]
```

Empty Response (when chapter has no techniques):

```json
[]
```

Error Responses:

- **400 Bad Request** - Invalid ID format or ID <= 0
- **404 Not Found** - Chapter with the given ID does not exist
- **500 Internal Server Error** - Database or server error

## Running tests

Use the Docker test runner script (requires the `hema-lessons_default` network and test DB env vars):

```
export TEST_DATABASE_HOST=localhost
export TEST_DATABASE_PORT=5432
export TEST_DATABASE_USER=postgres
export TEST_DATABASE_PASSWORD=postgres
export TEST_DATABASE_DBNAME=hema_lessons_test
export TEST_DATABASE_SSLMODE=disable
export TEST_APP_ENVIRONMENT=testing
./scripts/run-tests.sh
```
