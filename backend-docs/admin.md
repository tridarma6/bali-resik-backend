# Admin API

Admin endpoints for tenant management, collector application review, and analytics.

---

## POST /api/v1/admin/tenants

Create a new tenant (region).

**Authentication:** Bearer Token Required

**Allowed Roles:** `super_admin`

### Request
```json
{
  "name": "Kabupaten Badung",
  "slug": "badung",
  "region_type": "kabupaten"
}
```

**Region Types:** `kota`, `kabupaten`

### Success Response (201)
```json
{
  "success": true,
  "data": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "name": "Kabupaten Badung",
    "slug": "badung",
    "region_type": "kabupaten",
    "is_active": true,
    "created_at": "2026-06-13T10:00:00Z"
  },
  "message": "Tenant created successfully"
}
```

### Error Response (409)
```json
{
  "success": false,
  "error": {
    "code": "CONFLICT",
    "message": "Tenant slug already exists"
  }
}
```

---

## GET /api/v1/admin/tenants

List all tenants.

**Authentication:** Bearer Token Required

**Allowed Roles:** `super_admin`

### Success Response (200)
```json
{
  "success": true,
  "data": [
    {
      "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      "name": "Kabupaten Badung",
      "slug": "badung",
      "region_type": "kabupaten",
      "is_active": true,
      "created_at": "2026-06-01T00:00:00Z"
    }
  ],
  "message": "Tenants retrieved"
}
```

---

## POST /api/v1/admin/admins

Create a new admin user for a tenant.

**Authentication:** Bearer Token Required

**Allowed Roles:** `super_admin`

### Request
```json
{
  "email": "admin@badung.go.id",
  "password": "adminpassword123",
  "full_name": "Admin Badung",
  "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
}
```

### Success Response (201)
```json
{
  "success": true,
  "data": {
    "id": "9a8b7c6d-5e4f-3210-abcd-ef1234567890",
    "email": "admin@badung.go.id",
    "full_name": "Admin Badung",
    "phone": "",
    "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "roles": ["admin_kabupaten"],
    "is_active": true,
    "created_at": "2026-06-13T10:00:00Z"
  },
  "message": "Admin user created successfully"
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

## Collector Applications

### POST /api/v1/collector-applications

Submit a collector application (citizen applies to become a collector).

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

No request body required.

### Success Response (201)
```json
{
  "success": true,
  "data": {
    "id": "c9d0e1f2-3456-7890-abcd-ef1234567890",
    "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "user_id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
    "status": "pending",
    "admin_notes": "",
    "reviewed_by": null,
    "reviewed_at": null,
    "created_at": "2026-06-13T10:00:00Z",
    "user": null
  },
  "message": "Collector application submitted"
}
```

### Error Response (409)
```json
{
  "success": false,
  "error": {
    "code": "CONFLICT",
    "message": "You already have a pending application"
  }
}
```

---

### GET /api/v1/collector-applications/mine

List the authenticated user's collector applications.

**Authentication:** Bearer Token Required

**Allowed Roles:** Any authenticated user

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
      "id": "c9d0e1f2-3456-7890-abcd-ef1234567890",
      "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      "user_id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
      "status": "pending",
      "admin_notes": "",
      "reviewed_by": null,
      "reviewed_at": null,
      "created_at": "2026-06-13T10:00:00Z",
      "user": null
    }
  ],
  "message": "Applications retrieved",
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_page": 1
  }
}
```

---

### GET /api/v1/admin/collector-applications

List all collector applications (admin view).

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
      "id": "c9d0e1f2-3456-7890-abcd-ef1234567890",
      "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      "user_id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
      "status": "pending",
      "admin_notes": "",
      "reviewed_by": null,
      "reviewed_at": null,
      "created_at": "2026-06-13T10:00:00Z",
      "user": {
        "id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
        "email": "user@gmail.com",
        "full_name": "I Wayan Sudarma",
        "phone": "081234567890"
      }
    }
  ],
  "message": "Applications retrieved",
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_page": 1
  }
}
```

**Application Statuses:** `pending`, `approved`, `rejected`

---

### PUT /api/v1/admin/collector-applications/:id/approve

Approve a collector application. The applicant's role will be updated to `collector`.

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`, `super_admin`

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "c9d0e1f2-3456-7890-abcd-ef1234567890",
    "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "user_id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
    "status": "approved",
    "admin_notes": "",
    "reviewed_by": "9a8b7c6d-5e4f-3210-abcd-ef1234567890",
    "reviewed_at": "2026-06-13T14:00:00Z",
    "created_at": "2026-06-13T10:00:00Z",
    "user": null
  },
  "message": "Application approved"
}
```

### Error Response (400)
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Application is not in pending status"
  }
}
```

---

### PUT /api/v1/admin/collector-applications/:id/reject

Reject a collector application.

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`, `super_admin`

### Request
```json
{
  "admin_notes": "Mohon lengkapi persyaratan dokumen KTP dan SKCK"
}
```

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "id": "c9d0e1f2-3456-7890-abcd-ef1234567890",
    "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "user_id": "3a2b1c4d-5e6f-7890-abcd-ef1234567890",
    "status": "rejected",
    "admin_notes": "Mohon lengkapi persyaratan dokumen KTP dan SKCK",
    "reviewed_by": "9a8b7c6d-5e4f-3210-abcd-ef1234567890",
    "reviewed_at": "2026-06-13T14:00:00Z",
    "created_at": "2026-06-13T10:00:00Z",
    "user": null
  },
  "message": "Application rejected"
}
```

---

## Analytics

### GET /api/v1/analytics/dashboard

Get dashboard overview for the admin's tenant.

**Authentication:** Bearer Token Required

**Allowed Roles:** `admin_kabupaten`, `super_admin`

### Success Response (200)
```json
{
  "success": true,
  "data": {
    "overview": {
      "total_pickups": 150,
      "pickups_by_status": {
        "pending": 30,
        "assigned": 20,
        "in_progress": 15,
        "completed": 80,
        "cancelled": 5
      },
      "total_reports": 75,
      "reports_by_status": {
        "reported": 20,
        "verified": 15,
        "cleaning": 10,
        "resolved": 25,
        "rejected": 5
      },
      "total_citizens": 500,
      "total_collectors": 25,
      "pickup_completion_rate": 53.33
    },
    "pickup_trends": [
      { "year": 2026, "month": 5, "count": 45 },
      { "year": 2026, "month": 6, "count": 55 }
    ],
    "report_trends": [
      { "year": 2026, "month": 5, "count": 20 },
      { "year": 2026, "month": 6, "count": 30 }
    ],
    "waste_type_distribution": [
      { "waste_type": "organic", "count": 60 },
      { "waste_type": "anorganic", "count": 40 },
      { "waste_type": "mixed", "count": 25 },
      { "waste_type": "electronic", "count": 15 },
      { "waste_type": "hazardous", "count": 10 }
    ],
    "severity_distribution": [
      { "severity": "low", "count": 10 },
      { "severity": "medium", "count": 30 },
      { "severity": "high", "count": 25 },
      { "severity": "critical", "count": 10 }
    ],
    "new_users": [
      { "year": 2026, "month": 5, "count": 100 },
      { "year": 2026, "month": 6, "count": 120 }
    ]
  },
  "message": "Dashboard data retrieved"
}
```

---

### GET /api/v1/admin/regional-stats

Get regional statistics across all tenants (super admin only).

**Authentication:** Bearer Token Required

**Allowed Roles:** `super_admin`

### Success Response (200)
```json
{
  "success": true,
  "data": [
    {
      "tenant_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      "tenant_name": "Kabupaten Badung",
      "total_pickups": 150,
      "total_reports": 75,
      "total_users": 500
    }
  ],
  "message": "Regional statistics retrieved"
}
```
