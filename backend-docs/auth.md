# Authentication API

Base URL: `/api/v1/auth`

---

## POST /api/v1/auth/register

Register a new citizen account.

**Authentication:** Public

**Allowed Roles:** None

### Request
```json
{
  "email": "user@gmail.com",
  "password": "password123",
  "full_name": "I Wayan Sudarma",
  "phone": "081234567890",
  "tenant_slug": "badung"
}
```

### Success Response (201)
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJl...",
    "expires_in": 3600,
    "user": {
      "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
      "email": "user@gmail.com",
      "full_name": "I Wayan Sudarma",
      "phone": "081234567890",
      "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      "roles": ["citizen"],
      "is_active": true,
      "created_at": "2026-06-13T10:00:00Z"
    }
  },
  "message": "Registration successful"
}
```

### Error Response (409)
```json
{
  "success": false,
  "error": {
    "code": "CONFLICT",
    "message": "Email already registered"
  }
}
```

---

## POST /api/v1/auth/login

Authenticate with email and password.

**Authentication:** Public

**Allowed Roles:** None

### Request
```json
{
  "email": "user@gmail.com",
  "password": "password123"
}
```

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJl...",
    "expires_in": 3600,
    "user": {
      "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
      "email": "user@gmail.com",
      "full_name": "I Wayan Sudarma",
      "phone": "081234567890",
      "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      "roles": ["citizen"],
      "is_active": true,
      "created_at": "2026-06-13T10:00:00Z"
    }
  },
  "message": "Login successful"
}
```

### Error Response (401)
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid credentials"
  }
}
```

---

## POST /api/v1/auth/refresh

Obtain a new access token using a refresh token.

**Authentication:** Public

**Allowed Roles:** None

### Request
```json
{
  "refresh_token": "dGhpcyBpcyBhIHJlZnJl..."
}
```

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "bmV3IHJlZnJlc2ggdG9r...",
    "expires_in": 3600,
    "user": {
      "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
      "email": "user@gmail.com",
      "full_name": "I Wayan Sudarma",
      "phone": "081234567890",
      "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      "roles": ["citizen"],
      "is_active": true,
      "created_at": "2026-06-13T10:00:00Z"
    }
  },
  "message": "Token refreshed"
}
```

### Error Response (401)
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or expired refresh token"
  }
}
```

---

## POST /api/v1/auth/logout

Invalidate the current refresh token.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
X-Refresh-Token: dGhpcyBpcyBhIHJlZnJl...
```

### Success Response (200)
```json
{
  "success": true,
  "data": null,
  "message": "Logout successful"
}
```

### Error Response (401)
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "User not authenticated"
  }
}
```
