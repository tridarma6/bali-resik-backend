# Waste Reports API

Base URL: `/api/v1/reports`

---

## POST /api/v1/reports

Create a new waste report (illegal dumping, littering, etc.).

**Authentication:** Bearer Token Required

**Allowed Roles:** `citizen`

### Request
```json
{
  "latitude": -8.4095,
  "longitude": 115.1889,
  "description": "Illegal waste dumping near the riverbank at Jl. Tukad Badung",
  "severity": "high"
}
```

**Severity Levels:** `low`, `medium`, `high`, `critical`

### Success Response (201)
```json
{
  "success": true,
  "data": {
    "id": "b2c3d4e5-6f78-9012-abcd-ef1234567890",
    "latitude": -8.4095,
    "longitude": 115.1889,
    "description": "Illegal waste dumping near the riverbank at Jl. Tukad Badung",
    "status": "reported",
    "severity": "high",
    "created_at": "2026-06-13T10:00:00Z",
    "updated_at": "2026-06-13T10:00:00Z",
    "user": {
      "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
      "full_name": "I Wayan Sudarma",
      "email": "user@gmail.com",
      "phone": "081234567890"
    },
    "images": []
  },
  "message": "Waste report created"
}
```

### Error Response (400)
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Description must be at least 10 characters"
  }
}
```

---

## GET /api/v1/reports/mine

List waste reports submitted by the authenticated citizen.

**Authentication:** Bearer Token Required

**Allowed Roles:** `citizen`

### Query Parameters
| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number (default: 1) |
| `per_page` | int | Items per page (default: 20, max: 100) |
| `status` | string | Filter by status |

### Success Response (200)
```json
{
  "success": true,
  "data": [
    {
      "id": "b2c3d4e5-6f78-9012-abcd-ef1234567890",
      "latitude": -8.4095,
      "longitude": 115.1889,
      "description": "Illegal waste dumping near the riverbank at Jl. Tukad Badung",
      "status": "reported",
      "severity": "high",
      "created_at": "2026-06-13T10:00:00Z",
      "updated_at": "2026-06-13T10:00:00Z",
      "user": null,
      "images": []
    }
  ],
  "message": "My waste reports retrieved",
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_page": 1
  }
}
```

---

## GET /api/v1/reports/nearby

Find waste reports near a geographic location.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Query Parameters
| Param | Type | Description |
|-------|------|-------------|
| `lat` | float | Latitude (required) |
| `lng` | float | Longitude (required) |
| `radius` | float | Search radius in km (0.1 - 50) |

### Request
```
GET /api/v1/reports/nearby?lat=-8.4095&lng=115.1889&radius=5
```

### Success Response (200)
```json
{
  "success": true,
  "data": [
    {
      "id": "b2c3d4e5-6f78-9012-abcd-ef1234567890",
      "latitude": -8.4095,
      "longitude": 115.1889,
      "description": "Illegal waste dumping near the riverbank at Jl. Tukad Badung",
      "status": "reported",
      "severity": "high",
      "created_at": "2026-06-13T10:00:00Z",
      "updated_at": "2026-06-13T10:00:00Z",
      "user": {
        "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
        "full_name": "I Wayan Sudarma",
        "email": "user@gmail.com",
        "phone": "081234567890"
      },
      "images": []
    }
  ],
  "message": "Nearby waste reports retrieved"
}
```

---

## GET /api/v1/reports

List all waste reports (admin view).

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`, `super_admin`

### Query Parameters
| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number (default: 1) |
| `per_page` | int | Items per page (default: 20, max: 100) |
| `status` | string | Filter by status |

### Success Response (200)
```json
{
  "success": true,
  "data": [
    {
      "id": "b2c3d4e5-6f78-9012-abcd-ef1234567890",
      "latitude": -8.4095,
      "longitude": 115.1889,
      "description": "Illegal waste dumping near the riverbank at Jl. Tukad Badung",
      "status": "reported",
      "severity": "high",
      "created_at": "2026-06-13T10:00:00Z",
      "updated_at": "2026-06-13T10:00:00Z",
      "user": {
        "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
        "full_name": "I Wayan Sudarma",
        "email": "user@gmail.com",
        "phone": "081234567890"
      },
      "images": []
    }
  ],
  "message": "Waste reports retrieved",
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_page": 1
  }
}
```

---

## GET /api/v1/reports/:id

Get a single waste report by ID.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "b2c3d4e5-6f78-9012-abcd-ef1234567890",
    "latitude": -8.4095,
    "longitude": 115.1889,
    "description": "Illegal waste dumping near the riverbank at Jl. Tukad Badung",
    "status": "verified",
    "severity": "high",
    "created_at": "2026-06-13T10:00:00Z",
    "updated_at": "2026-06-13T14:00:00Z",
    "user": {
      "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
      "full_name": "I Wayan Sudarma",
      "email": "user@gmail.com",
      "phone": "081234567890"
    },
    "images": [
      {
        "id": "c3d4e5f6-7890-1234-abcd-ef1234567890",
        "image_url": "/uploads/reports/b2c3d4e5_1718251200000000000.jpg"
      }
    ]
  },
  "message": "Waste report retrieved"
}
```

### Error Response (404)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Waste report not found"
  }
}
```

---

## PUT /api/v1/reports/:id/status

Update the status of a waste report.

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`

### Request
```json
{
  "status": "verified"
}
```

**Valid Statuses:** `verified`, `cleaning`, `resolved`, `rejected`

**Status Flow:** `reported` → `verified` → `cleaning` → `resolved`

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "b2c3d4e5-6f78-9012-abcd-ef1234567890",
    "latitude": -8.4095,
    "longitude": 115.1889,
    "description": "Illegal waste dumping near the riverbank at Jl. Tukad Badung",
    "status": "verified",
    "severity": "high",
    "created_at": "2026-06-13T10:00:00Z",
    "updated_at": "2026-06-13T14:00:00Z",
    "user": null,
    "images": []
  },
  "message": "Report status updated"
}
```

### Error Response (400)
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid status transition"
  }
}
```

---

## POST /api/v1/reports/:id/images

Upload an image to a waste report.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Request (multipart/form-data)
| Field | Type | Description |
|-------|------|-------------|
| `image` | file | Image file (JPG, PNG, WebP, max 10MB) |

### Success Response (201)
```json
{
  "success": true,
  "data": {
    "image_url": "/uploads/reports/b2c3d4e5_1718251200000000000.jpg"
  },
  "message": "Image uploaded successfully"
}
```

### Error Response (400)
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Only JPG, PNG, and WebP images are allowed"
  }
}
```
