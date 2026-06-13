# Notifications API

Base URL: `/api/v1/notifications`

---

## GET /api/v1/notifications

List notifications for the authenticated user.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Query Parameters
| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number (default: 1) |
| `per_page` | int | Items per page (default: 20, max: 100) |
| `read` | string | Filter: `"true"`, `"false"`, or empty (all) |

### Success Response (200)
```json
{
  "success": true,
  "data": [
    {
      "id": "b8c9d0e1-2345-6789-abcd-ef1234567890",
      "title": "Pickup Completed",
      "message": "Your pickup request for organic waste has been completed.",
      "type": "pickup",
      "reference_id": "a1b2c3d4-5e6f-7890-abcd-ef1234567890",
      "is_read": false,
      "created_at": "2026-06-13T10:00:00Z"
    }
  ],
  "message": "Notifications retrieved",
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_page": 1
  }
}
```

**Notification Types:** `pickup`, `report`, `reward`, `system`

---

## GET /api/v1/notifications/unread-count

Get the count of unread notifications.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "count": 3
  },
  "message": "Unread count retrieved"
}
```

---

## PUT /api/v1/notifications/:id/read

Mark a single notification as read.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Success Response (200)
```json
{
  "success": true,
  "data": null,
  "message": "Notification marked as read"
}
```

### Error Response (404)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Notification not found"
  }
}
```

---

## PUT /api/v1/notifications/read-all

Mark all notifications as read for the authenticated user.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Success Response (200)
```json
{
  "success": true,
  "data": null,
  "message": "All notifications marked as read"
}
```

---

## DELETE /api/v1/notifications/:id

Delete a notification.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

### Success Response (200)
```json
{
  "success": true,
  "data": null,
  "message": "Notification deleted"
}
```

### Error Response (404)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Notification not found"
  }
}
```
