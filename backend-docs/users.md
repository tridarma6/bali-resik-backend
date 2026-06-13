# Users API

Base URL: `/api/v1/users`

---

## GET /api/v1/users/me

Get the authenticated user's profile.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
    "email": "user@gmail.com",
    "full_name": "I Wayan Sudarma",
    "phone": "081234567890",
    "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "roles": ["citizen"],
    "is_active": true,
    "created_at": "2026-06-13T10:00:00Z"
  },
  "message": "Profile retrieved"
}
```

---

## PUT /api/v1/users/me

Update the authenticated user's profile.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Request
```json
{
  "full_name": "I Wayan Sudarma Putra",
  "phone": "081234567891"
}
```

All fields are optional — only send fields to update.

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
    "email": "user@gmail.com",
    "full_name": "I Wayan Sudarma Putra",
    "phone": "081234567891",
    "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "roles": ["citizen"],
    "is_active": true,
    "created_at": "2026-06-13T10:00:00Z"
  },
  "message": "Profile updated"
}
```

---

## PUT /api/v1/users/me/password

Change the authenticated user's password.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Request
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

### Success Response (200)
```json
{
  "success": true,
  "data": null,
  "message": "Password changed"
}
```

### Error Response (401)
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Current password is incorrect"
  }
}
```

---

## POST /api/v1/users/me/avatar

Upload an avatar image for the authenticated user.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Request (multipart/form-data)
| Field | Type | Description |
|-------|------|-------------|
| `avatar` | file | Image file (JPG, PNG, WebP, max 5MB) |

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
    "email": "user@gmail.com",
    "full_name": "I Wayan Sudarma",
    "phone": "081234567890",
    "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "roles": ["citizen"],
    "is_active": true,
    "created_at": "2026-06-13T10:00:00Z"
  },
  "message": "Avatar updated"
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

---

## User Response Object

The `UserResponse` is used across the API for user profile data.

### Fields
| Field | Type | Description |
|-------|------|-------------|
| `id` | uuid | User ID |
| `email` | string | Email address |
| `full_name` | string | Full name |
| `phone` | string | Phone number (optional) |
| `tenant_id` | uuid | Tenant/region ID |
| `roles` | string[] | Role names (e.g. `["citizen"]`) |
| `is_active` | bool | Account active status |
| `created_at` | datetime | Account creation time |

### User Reference (`UserRef`)

Used in pickup, report, and education responses.
```json
{
  "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
  "full_name": "I Wayan Sudarma",
  "email": "user@gmail.com",
  "phone": "081234567890"
}
```

### User Brief (`UserBriefResponse`)

Used in collector application responses.
```json
{
  "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
  "email": "user@gmail.com",
  "full_name": "I Wayan Sudarma",
  "phone": "081234567890"
}
```

---

## Roles

| Role | Description |
|------|-------------|
| `citizen` | Regular user who can create pickups and reports |
| `collector` | Waste collector assigned to pickups |
| `admin_kabupaten` | Regional admin who manages their tenant region |
| `super_admin` | Global admin who manages tenants |

Accounts are created via:
- `POST /api/v1/auth/register` — creates a user with `citizen` role
- `POST /api/v1/admin/admins` — creates a user with `admin_kabupaten` role
