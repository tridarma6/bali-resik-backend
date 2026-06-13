# Educational Content API

Base URL: `/api/v1/education`

---

## GET /api/v1/education

List educational content.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Query Parameters
| Param | Type | Description |
|-------|------|-------------|
| `category` | string | Filter by category (e.g. "pemilahan", "daur-ulang") |
| `published` | bool | Show only published content (default: false) |

### Success Response (200)
```json
{
  "success": true,
  "data": [
    {
      "id": "a7b8c9d0-1234-5678-abcd-ef1234567890",
      "title": "Cara Memilah Sampah Rumah Tangga",
      "content": "Pemilahan sampah adalah langkah awal yang penting dalam pengelolaan waste yang baik...",
      "content_type": "article",
      "category": "pemilahan",
      "image_url": "https://storage.example.com/education/pilah-sampah.jpg",
      "is_published": true,
      "created_at": "2026-06-01T00:00:00Z",
      "updated_at": "2026-06-10T00:00:00Z",
      "author": {
        "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
        "full_name": "I Wayan Sudarma",
        "email": "user@gmail.com",
        "phone": "081234567890"
      }
    }
  ],
  "message": "Content list retrieved"
}
```

**Content Types:** `article`, `video`, `infographic`

---

## GET /api/v1/education/:id

Get a single educational content item by ID.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "a7b8c9d0-1234-5678-abcd-ef1234567890",
    "title": "Cara Memilah Sampah Rumah Tangga",
    "content": "Pemilahan sampah adalah langkah awal yang penting dalam pengelolaan waste yang baik...",
    "content_type": "article",
    "category": "pemilahan",
    "image_url": "https://storage.example.com/education/pilah-sampah.jpg",
    "is_published": true,
    "created_at": "2026-06-01T00:00:00Z",
    "updated_at": "2026-06-10T00:00:00Z",
    "author": {
      "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
      "full_name": "I Wayan Sudarma",
      "email": "user@gmail.com",
      "phone": "081234567890"
    }
  },
  "message": "Content retrieved"
}
```

### Error Response (404)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Content not found"
  }
}
```

---

## POST /api/v1/education

Create educational content (admin).

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`, `super_admin`

### Request
```json
{
  "title": "Cara Memilah Sampah Rumah Tangga",
  "content": "Pemilahan sampah adalah langkah awal yang penting dalam pengelolaan waste yang baik. Pisahkan sampah organik, anorganik, dan limbah berbahaya...",
  "content_type": "article",
  "category": "pemilahan",
  "image_url": "https://storage.example.com/education/pilah-sampah.jpg"
}
```

### Success Response (201)
```json
{
  "success": true,
  "data": {
    "id": "a7b8c9d0-1234-5678-abcd-ef1234567890",
    "title": "Cara Memilah Sampah Rumah Tangga",
    "content": "Pemilahan sampah adalah langkah awal yang penting dalam pengelolaan waste yang baik. Pisahkan sampah organik, anorganik, dan limbah berbahaya...",
    "content_type": "article",
    "category": "pemilahan",
    "image_url": "https://storage.example.com/education/pilah-sampah.jpg",
    "is_published": false,
    "created_at": "2026-06-13T10:00:00Z",
    "updated_at": "2026-06-13T10:00:00Z",
    "author": {
      "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
      "full_name": "Admin User",
      "email": "admin@badung.go.id",
      "phone": "081234567890"
    }
  },
  "message": "Content created"
}
```

---

## PUT /api/v1/education/:id

Update educational content (admin).

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`, `super_admin`

### Request
```json
{
  "title": "Updated Title",
  "is_published": true
}
```

All fields are optional — only send fields to update.

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "a7b8c9d0-1234-5678-abcd-ef1234567890",
    "title": "Updated Title",
    "content": "Pemilahan sampah adalah langkah awal yang penting dalam pengelolaan waste yang baik...",
    "content_type": "article",
    "category": "pemilahan",
    "image_url": "https://storage.example.com/education/pilah-sampah.jpg",
    "is_published": true,
    "created_at": "2026-06-01T00:00:00Z",
    "updated_at": "2026-06-13T12:00:00Z",
    "author": null
  },
  "message": "Content updated"
}
```

---

## DELETE /api/v1/education/:id

Delete educational content (admin).

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`, `super_admin`

### Success Response (200)
```json
{
  "success": true,
  "data": null,
  "message": "Content deleted"
}
```

### Error Response (404)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Content not found"
  }
}
```
