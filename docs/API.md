# API Documentation

## Base URL

`http://localhost:8080` (configurable via `SERVER_ADDR` environment variable)

---

## Endpoints

### Health Check

**GET /healthz**

Returns the health status of the API.

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

In development mode, the status includes the environment name:

```json
{
  "status": "ok (development)",
  "timestamp": "2026-01-18T12:34:56Z"
}
```

---

## Resources

A **resource** represents any learning material — a historical manuscript, a modern course, or any other structured content.

### List Resources

**GET /api/resources**

Returns a paginated list of resources, ordered alphabetically by title. Each item includes the author name when one is associated.

Query Parameters:

| Parameter   | Type | Default | Description                         |
|-------------|------|---------|-------------------------------------|
| `page`      | int  | 1       | Page number (minimum 1)             |
| `page_size` | int  | 20      | Items per page (minimum 1, max 100) |

```bash
curl "http://localhost:8080/api/resources?page=1&page_size=2"
```

Response:

```json
{
  "data": [
    {
      "id": 4,
      "author_id": 4,
      "title": "De Arte Gladiatoria Dimicandi",
      "description": "On the Art of Swordsmanship...",
      "publication_year": 1482,
      "cover_image_url": "/assets/books/de-arte-gladiatoria-dimicandi/cover/cover.jpg",
      "author_name": "Filippo Vadi"
    },
    {
      "id": 3,
      "author_id": 3,
      "title": "Fechtbuch",
      "description": "A detailed commentary on Liechtenauer's teachings with practical applications",
      "publication_year": 1440,
      "cover_image_url": "/assets/books/fechtbuch/cover/cover.jpg",
      "author_name": "Hans Talhoffer"
    }
  ],
  "page": 1,
  "page_size": 2,
  "total_count": 4,
  "total_pages": 2
}
```

Fields:

| Field               | Type    | Description                                             |
|---------------------|---------|---------------------------------------------------------|
| `id`                | int     | Resource ID                                             |
| `author_id`         | int     | ID of the associated author (omitted if none)           |
| `title`             | string  | Resource title                                          |
| `description`       | string  | Short description                                       |
| `publication_year`  | int     | Year of publication/creation (omitted if unknown)       |
| `cover_image_url`   | string  | URL to the cover image (omitted if none)                |
| `author_name`       | string  | Author's name (omitted if no author is linked)          |

Error Responses:

- **500 Internal Server Error** — server error

---

### Get Resource by ID

**GET /api/resources/{id}**

Returns a single resource by its ID, including the author name.

Path Parameters:

| Parameter | Type | Description                    |
|-----------|------|--------------------------------|
| `id`      | int  | Resource ID (must be > 0)      |

```bash
curl http://localhost:8080/api/resources/2
```

Response:

```json
{
  "id": 2,
  "author_id": 2,
  "title": "Fior di Battaglia",
  "description": "The Flower of Battle - a comprehensive medieval combat manual covering armed and unarmed combat",
  "publication_year": 1409,
  "cover_image_url": "/assets/books/fior-di-battaglia/cover/cover.jpg",
  "author_name": "Fiore dei Liberi"
}
```

Error Responses:

- **400 Bad Request** — invalid ID format or ID <= 0
- **404 Not Found** — resource with the given ID does not exist

---

### List Root Sections by Resource

**GET /api/resources/{id}/sections**

Returns the top-level sections of a resource, ordered by `position`. Only root sections are returned (sections with no parent). Use `GET /api/sections/{id}/sections` to retrieve nested sections.

Path Parameters:

| Parameter | Type | Description                    |
|-----------|------|--------------------------------|
| `id`      | int  | Resource ID (must be > 0)      |

```bash
curl http://localhost:8080/api/resources/2/sections
```

Response:

```json
[
  {
    "id": 1,
    "resource_id": 2,
    "kind": "chapter",
    "title": "Abrazare (Wrestling)",
    "description": "Unarmed combat and grappling techniques forming the foundation of Fiore's system",
    "position": 1
  },
  {
    "id": 2,
    "resource_id": 2,
    "kind": "chapter",
    "title": "Dagger (Daga)",
    "description": "Fighting with and against the dagger at close range, including defenses against common attacks",
    "position": 2
  }
]
```

Empty response (when resource has no sections):

```json
[]
```

Fields:

| Field         | Type   | Description                                                |
|---------------|--------|------------------------------------------------------------|
| `id`          | int    | Section ID                                                 |
| `resource_id` | int    | ID of the parent resource                                  |
| `kind`        | string | Section type label (e.g. `"chapter"`, `"sub-chapter"`)     |
| `title`       | string | Section title                                              |
| `description` | string | Section description                                        |
| `position`    | int    | Ordering within the parent (1-based)                       |

Error Responses:

- **400 Bad Request** — invalid ID format or ID <= 0
- **404 Not Found** — resource with the given ID does not exist

---

## Sections

### Get Section by ID

**GET /api/sections/{id}**

Returns a single section by its ID.

Path Parameters:

| Parameter | Type | Description                   |
|-----------|------|-------------------------------|
| `id`      | int  | Section ID (must be > 0)      |

```bash
curl http://localhost:8080/api/sections/1
```

Response:

```json
{
  "id": 1,
  "resource_id": 2,
  "kind": "chapter",
  "title": "Abrazare (Wrestling)",
  "description": "Unarmed combat and grappling techniques forming the foundation of Fiore's system",
  "position": 1
}
```

For nested sections, `parent_id` is also present:

```json
{
  "id": 17,
  "resource_id": 2,
  "parent_id": 4,
  "kind": "sub-chapter",
  "title": "Zornhau Plays",
  "description": "Counter-plays from the Wrath Strike",
  "position": 1
}
```

Error Responses:

- **400 Bad Request** — invalid ID format or ID <= 0
- **404 Not Found** — section with the given ID does not exist

---

### List Child Sections

**GET /api/sections/{id}/sections**

Returns the direct child sections of a section, ordered by `position`.

Path Parameters:

| Parameter | Type | Description                   |
|-----------|------|-------------------------------|
| `id`      | int  | Section ID (must be > 0)      |

```bash
curl http://localhost:8080/api/sections/4/sections
```

Response:

```json
[
  {
    "id": 17,
    "resource_id": 2,
    "parent_id": 4,
    "kind": "sub-chapter",
    "title": "Zornhau Plays",
    "description": "Counter-plays from the Wrath Strike",
    "position": 1
  }
]
```

Empty response (when section has no children):

```json
[]
```

Error Responses:

- **400 Bad Request** — invalid ID format or ID <= 0
- **404 Not Found** — section with the given ID does not exist

---

### List Items by Section

**GET /api/sections/{id}/items**

Returns the items belonging to a section, ordered by `position`.

Path Parameters:

| Parameter | Type | Description                   |
|-----------|------|-------------------------------|
| `id`      | int  | Section ID (must be > 0)      |

```bash
curl http://localhost:8080/api/sections/1/items
```

Response:

```json
[
  {
    "id": 1,
    "section_id": 1,
    "kind": "technique",
    "title": "First Remedy Master of Abrazare",
    "description": "The foundational wrestling master position",
    "position": 1,
    "attributes": {
      "instructions": "As the opponent reaches to grab you, step offline...",
      "historical_image_url": "/assets/books/fior-di-battaglia/techniques/first-remedy-master-of-abrazare/historical.jpg"
    }
  },
  {
    "id": 2,
    "section_id": 1,
    "kind": "technique",
    "title": "Ligadura Soprana (Upper Lock)",
    "description": "An arm lock that forces the opponent's arm upward",
    "position": 2,
    "attributes": {
      "instructions": "From a grip on the opponent's right arm, thread your right arm under their elbow...",
      "historical_image_url": "/assets/books/fior-di-battaglia/techniques/ligadura-soprana-upper-lock/historical.jpg"
    }
  }
]
```

Empty response (when section has no items):

```json
[]
```

Fields:

| Field         | Type   | Description                                               |
|---------------|--------|-----------------------------------------------------------|
| `id`          | int    | Item ID                                                   |
| `section_id`  | int    | ID of the parent section                                  |
| `kind`        | string | Item type label (e.g. `"technique"`, `"drill"`)           |
| `title`       | string | Item title                                                |
| `description` | string | Short description                                         |
| `position`    | int    | Ordering within the section (1-based)                     |
| `attributes`  | object | Free-form key/value map; shape varies by `kind` (omitted if empty) |

Error Responses:

- **400 Bad Request** — invalid ID format or ID <= 0
- **404 Not Found** — section with the given ID does not exist

---

## Running Tests

Tests use an in-memory store and require no external services. Run them with Docker:

```bash
docker run --rm \
  -v "$(pwd)":/app \
  -w /app \
  golang:1.22-alpine \
  sh -c "go test -v ./..."
```

Or directly if Go is installed:

```bash
go test ./...
```
