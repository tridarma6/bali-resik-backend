# Pickups API

Base URL: `/api/v1/pickups`

---

## POST /api/v1/pickups

Create a new pickup request.

**Authentication:** Bearer Token Required

**Allowed Roles:** `citizen`

### Request
```json
{
  "waste_type": "organic",
  "latitude": -8.4095,
  "longitude": 115.1889,
  "address": "Jl. Raya Kuta No. 123, Badung",
  "scheduled_date": "2026-06-15T08:00:00Z",
  "notes": "Please pick up in the morning"
}
```

**Waste Types:** `organic`, `anorganic`, `mixed`, `electronic`, `hazardous`

### Success Response (201)
```json
{
  "success": true,
  "data": {
    "id": "a1b2c3d4-5e6f-7890-abcd-ef1234567890",
    "waste_type": "organic",
    "latitude": -8.4095,
    "longitude": 115.1889,
    "address": "Jl. Raya Kuta No. 123, Badung",
    "status": "pending",
    "scheduled_date": "2026-06-15T08:00:00Z",
    "notes": "Please pick up in the morning",
    "created_at": "2026-06-13T10:00:00Z",
    "updated_at": "2026-06-13T10:00:00Z",
    "user": {
      "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
      "full_name": "I Wayan Sudarma",
      "email": "user@gmail.com",
      "phone": "081234567890"
    },
    "collector": null
  },
  "message": "Pickup request created"
}
```

### Error Response (400)
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid waste type"
  }
}
```

---

## GET /api/v1/pickups/mine

List pickup requests for the authenticated citizen.

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
      "id": "a1b2c3d4-5e6f-7890-abcd-ef1234567890",
      "waste_type": "organic",
      "latitude": -8.4095,
      "longitude": 115.1889,
      "address": "Jl. Raya Kuta No. 123, Badung",
      "status": "pending",
      "scheduled_date": null,
      "notes": "",
      "created_at": "2026-06-13T10:00:00Z",
      "updated_at": "2026-06-13T10:00:00Z",
      "user": null,
      "collector": null
    }
  ],
  "message": "My pickup requests retrieved",
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_page": 1
  }
}
```

---

## GET /api/v1/pickups/assigned

List pickups assigned to the authenticated collector.

**Authentication:** Bearer Token Required

**Allowed Roles:** `collector`

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
      "id": "a1b2c3d4-5e6f-7890-abcd-ef1234567890",
      "waste_type": "organic",
      "latitude": -8.4095,
      "longitude": 115.1889,
      "address": "Jl. Raya Kuta No. 123, Badung",
      "status": "assigned",
      "scheduled_date": null,
      "notes": "",
      "created_at": "2026-06-13T10:00:00Z",
      "updated_at": "2026-06-13T10:00:00Z",
      "user": {
        "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
        "full_name": "I Wayan Sudarma",
        "email": "user@gmail.com",
        "phone": "081234567890"
      },
      "collector": null
    }
  ],
  "message": "Assigned pickups retrieved",
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_page": 1
  }
}
```

---

## GET /api/v1/pickups

List all pickup requests (admin view).

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
      "id": "a1b2c3d4-5e6f-7890-abcd-ef1234567890",
      "waste_type": "organic",
      "latitude": -8.4095,
      "longitude": 115.1889,
      "address": "Jl. Raya Kuta No. 123, Badung",
      "status": "pending",
      "scheduled_date": null,
      "notes": "",
      "created_at": "2026-06-13T10:00:00Z",
      "updated_at": "2026-06-13T10:00:00Z",
      "user": {
        "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
        "full_name": "I Wayan Sudarma",
        "email": "user@gmail.com",
        "phone": "081234567890"
      },
      "collector": null
    }
  ],
  "message": "Pickup requests retrieved",
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_page": 1
  }
}
```

---

## GET /api/v1/pickups/:id

Get a single pickup request by ID.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "a1b2c3d4-5e6f-7890-abcd-ef1234567890",
    "waste_type": "organic",
    "latitude": -8.4095,
    "longitude": 115.1889,
    "address": "Jl. Raya Kuta No. 123, Badung",
    "status": "assigned",
    "scheduled_date": null,
    "notes": "",
    "created_at": "2026-06-13T10:00:00Z",
    "updated_at": "2026-06-13T12:00:00Z",
    "user": {
      "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
      "full_name": "I Wayan Sudarma",
      "email": "user@gmail.com",
      "phone": "081234567890"
    },
    "collector": {
      "id": "4b5c6d7e-8f90-1234-abcd-ef1234567890",
      "full_name": "Ketut Pasek",
      "email": "collector@example.com",
      "phone": "089876543210"
    }
  },
  "message": "Pickup request retrieved"
}
```

### Error Response (404)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Pickup request not found"
  }
}
```

---

## PUT /api/v1/pickups/:id/assign

Assign a collector to a pickup request. Transitions status to `assigned`.

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`, `super_admin`

### Request
```json
{
  "collector_id": "4b5c6d7e-8f90-1234-abcd-ef1234567890"
}
```

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "a1b2c3d4-5e6f-7890-abcd-ef1234567890",
    "waste_type": "organic",
    "latitude": -8.4095,
    "longitude": 115.1889,
    "address": "Jl. Raya Kuta No. 123, Badung",
    "status": "assigned",
    "scheduled_date": null,
    "notes": "",
    "created_at": "2026-06-13T10:00:00Z",
    "updated_at": "2026-06-13T12:00:00Z",
    "user": null,
    "collector": {
      "id": "4b5c6d7e-8f90-1234-abcd-ef1234567890",
      "full_name": "Ketut Pasek",
      "email": "collector@example.com",
      "phone": "089876543210"
    }
  },
  "message": "Collector assigned"
}
```

### Error Response (404)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Collector not found"
  }
}
```

---

## PUT /api/v1/pickups/:id/status

Update the status of a pickup request.

**Authentication:** Bearer Token Required

**Allowed Roles:** `collector`, `admin_kabupaten`

### Request
```json
{
  "status": "in_progress"
}
```

**Valid Statuses:** `in_progress`, `completed`, `cancelled`

**Status Flow:** `pending` → `assigned` → `in_progress` → `completed`

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "a1b2c3d4-5e6f-7890-abcd-ef1234567890",
    "waste_type": "organic",
    "latitude": -8.4095,
    "longitude": 115.1889,
    "address": "Jl. Raya Kuta No. 123, Badung",
    "status": "in_progress",
    "scheduled_date": null,
    "notes": "",
    "created_at": "2026-06-13T10:00:00Z",
    "updated_at": "2026-06-13T13:00:00Z",
    "user": null,
    "collector": null
  },
  "message": "Pickup status updated"
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

## DELETE /api/v1/pickups/:id

Cancel a pickup request. Only the owner can cancel their own pickup.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user (owner)

### Success Response (200)
```json
{
  "success": true,
  "data": null,
  "message": "Pickup request cancelled"
}
```

### Error Response (403)
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "You are not allowed to cancel this pickup"
  }
}
```
