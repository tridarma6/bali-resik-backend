# Rewards API

Base URL: `/api/v1/rewards`

Users earn points through waste-related activities (pickups, reports) and redeem them for rewards.

---

## GET /api/v1/rewards

List all available rewards.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Success Response (200)
```json
{
  "success": true,
  "data": [
    {
      "id": "d4e5f6a7-8901-2345-abcd-ef1234567890",
      "name": "Sembako Paket Sejahtera",
      "description": "Beras 5kg, Minyak Goreng 2L, Gula 1kg",
      "points_cost": 500,
      "stock": 50,
      "image_url": "https://storage.example.com/rewards/sembako.jpg",
      "is_active": true,
      "created_at": "2026-06-01T00:00:00Z"
    }
  ],
  "message": "Rewards retrieved"
}
```

---

## GET /api/v1/rewards/points

Get the authenticated user's total points.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "total_points": 1250
  },
  "message": "Points retrieved"
}
```

---

## GET /api/v1/rewards/transactions

Get the authenticated user's point transaction history.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Query Parameters
| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number (default: 1) |
| `per_page` | int | Items per page (default: 20, max: 100) |

### Success Response (200)
```json
{
  "success": true,
  "data": [
    {
      "id": "e5f6a7b8-9012-3456-abcd-ef1234567890",
      "points": 100,
      "type": "earn",
      "description": "Points earned from pickup request completion",
      "created_at": "2026-06-13T10:00:00Z"
    },
    {
      "id": "f6a7b8c9-0123-4567-abcd-ef1234567890",
      "points": -500,
      "type": "redeem",
      "description": "Redeemed Sembako Paket Sejahtera",
      "created_at": "2026-06-14T10:00:00Z"
    }
  ],
  "message": "Transaction history retrieved",
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 2,
    "total_page": 1
  }
}
```

**Transaction Types:** `earn` (points added), `redeem` (points deducted)

---

## POST /api/v1/rewards/redeem

Redeem points for a reward.

**Authentication:** Bearer Token Required

**Allowed Roles:** `citizen`

### Request
```json
{
  "reward_id": "d4e5f6a7-8901-2345-abcd-ef1234567890"
}
```

### Success Response (201)
```json
{
  "success": true,
  "data": {
    "id": "f6a7b8c9-0123-4567-abcd-ef1234567890",
    "points": -500,
    "type": "redeem",
    "description": "Redeemed Sembako Paket Sejahtera",
    "created_at": "2026-06-14T10:00:00Z"
  },
  "message": "Reward redeemed successfully"
}
```

### Error Response (400)
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Insufficient points"
  }
}
```

---

## POST /api/v1/rewards

Create a new reward (admin).

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`, `super_admin`

### Request
```json
{
  "name": "Sembako Paket Sejahtera",
  "description": "Beras 5kg, Minyak Goreng 2L, Gula 1kg",
  "points_cost": 500,
  "stock": 100,
  "image_url": "https://storage.example.com/rewards/sembako.jpg"
}
```

### Success Response (201)
```json
{
  "success": true,
  "data": {
    "id": "d4e5f6a7-8901-2345-abcd-ef1234567890",
    "name": "Sembako Paket Sejahtera",
    "description": "Beras 5kg, Minyak Goreng 2L, Gula 1kg",
    "points_cost": 500,
    "stock": 100,
    "image_url": "https://storage.example.com/rewards/sembako.jpg",
    "is_active": true,
    "created_at": "2026-06-01T00:00:00Z"
  },
  "message": "Reward created"
}
```

---

## PUT /api/v1/rewards/:id

Update an existing reward (admin).

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`, `super_admin`

### Request
```json
{
  "name": "Sembako Paket Sejahtera Plus",
  "points_cost": 600,
  "stock": 75,
  "is_active": true
}
```

All fields are optional — only send fields to update.

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "d4e5f6a7-8901-2345-abcd-ef1234567890",
    "name": "Sembako Paket Sejahtera Plus",
    "description": "Beras 5kg, Minyak Goreng 2L, Gula 1kg",
    "points_cost": 600,
    "stock": 75,
    "image_url": "https://storage.example.com/rewards/sembako.jpg",
    "is_active": true,
    "created_at": "2026-06-01T00:00:00Z"
  },
  "message": "Reward updated"
}
```

### Error Response (404)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Reward not found"
  }
}
```

---

## DELETE /api/v1/rewards/:id

Delete a reward (admin).

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`, `super_admin`

### Success Response (200)
```json
{
  "success": true,
  "data": null,
  "message": "Reward deleted"
}
```

### Error Response (404)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Reward not found"
  }
}
```
