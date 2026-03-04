# API Quick Reference

## Create User

**POST** `/auth/admin/create-user`

### Request
```json
{
  "name": "John Doe",
  "number": "9876543210",
  "email": "john.doe@company.com",
  "company_id": 1,
  "password": "SecurePassword123!",
  "user_type": "admin",
  "role_name": "admin"
}
```

### Response (201 Created)
```json
{
  "success": true,
  "data": {
    "id": 8,
    "email": "john.doe@company.com",
    "phone": "9876543210",
    "username": "John Doe",
    "user_type": "admin",
    "role": "admin",
    "status": "active",
    "company_id": 1,
    "company": {
      "id": 1,
      "company_name": "ABC Enterprises",
      "business_type_id": 2,
      "gst_number": "18AABCU1234H1Z5",
      "pan_number": "AABCU1234H",
      "created_at": "2026-01-15T10:30:00+05:30",
      "updated_at": "2026-02-20T14:45:00+05:30"
    },
    "created_at": "2026-03-04T01:15:22.456+05:30"
  }
}
```

---

## Get All Users

**GET** `/auth/admin/users?page=1&limit=10&search=&role=`

### Request
```
Headers:
Authorization: Bearer <token>

Query Parameters:
- page (optional, default: 1)
- limit (optional, default: 10, max: 100)
- search (optional)
- role (optional)
```

### Response (200 OK)
```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "email": "admin@bbcloud.app",
      "username": "admin",
      "user_type": "superadmin",
      "role": "admin",
      "status": "active",
      "created_at": "2026-01-21T23:05:35+05:30"
    },
    {
      "id": 8,
      "email": "john.doe@company.com",
      "username": "John Doe",
      "user_type": "admin",
      "role": "admin",
      "phone": "9876543210",
      "status": "active",
      "company_name": "ABC Enterprises",
      "created_at": "2026-03-04T00:58:07.317+05:30",
      "created_by": 3
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 2,
    "total_pages": 1
  }
}
```

---

## Get User by ID

**GET** `/auth/admin/users/:id`

### Request
```
Headers:
Authorization: Bearer <token>

URL Parameter:
- id (required)
```

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "id": 8,
    "email": "john.doe@company.com",
    "phone": "9876543210",
    "username": "John Doe",
    "user_type": "admin",
    "role": "admin",
    "status": "active",
    "company_id": 1,
    "company": {
      "id": 1,
      "company_name": "ABC Enterprises",
      "business_type_id": 2,
      "gst_number": "18AABCU1234H1Z5",
      "pan_number": "AABCU1234H",
      "created_at": "2026-01-15T10:30:00+05:30",
      "updated_at": "2026-02-20T14:45:00+05:30"
    },
    "created_at": "2026-03-04T00:58:07.317+05:30",
    "created_by": 3,
    "last_login_at": "2026-03-04T10:15:22.456+05:30"
  }
}
```

---

## Update User

**PUT** `/auth/admin/users/:id`

### Request
```json
{
  "username": "John Doe Updated",
  "phone": "9876543210",
  "email": "john.updated@company.com"
}
```

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "id": 8,
    "email": "john.updated@company.com",
    "phone": "9876543210",
    "username": "John Doe Updated",
    "user_type": "admin",
    "role": "admin",
    "status": "active",
    "company_id": 1,
    "created_at": "2026-03-04T00:58:07.317+05:30",
    "updated_at": "2026-03-04T10:20:00.123+05:30"
  }
}
```

---

## Delete User

**DELETE** `/auth/admin/users/:id`

### Request
```
Headers:
Authorization: Bearer <token>

URL Parameter:
- id (required)
```

### Response (200 OK)
```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

---

## Update User Status

**PUT** `/auth/admin/users/:id/status`

### Request
```json
{
  "status": "inactive"
}
```

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "id": 8,
    "email": "john.doe@company.com",
    "username": "John Doe",
    "user_type": "admin",
    "role": "admin",
    "status": "inactive",
    "created_at": "2026-03-04T00:58:07.317+05:30"
  }
}
```

---

## Update User Role

**PUT** `/auth/admin/users/:id/role`

### Request
```json
{
  "role_name": "partner"
}
```

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "id": 8,
    "email": "john.doe@company.com",
    "username": "John Doe",
    "user_type": "admin",
    "role": "partner",
    "status": "active",
    "created_at": "2026-03-04T00:58:07.317+05:30"
  }
}
```

---

## Create Role

**POST** `/roles`

### Request
```json
{
  "role_name": "manager",
  "permissions": ["users.read", "users.update", "reports.read"],
  "description": "Manager role with read and update permissions",
  "is_active": true
}
```

### Response (201 Created)
```json
{
  "success": true,
  "data": {
    "id": 3,
    "role_name": "manager",
    "permissions": ["users.read", "users.update", "reports.read"],
    "description": "Manager role with read and update permissions",
    "is_active": true,
    "created_at": "2026-03-04T10:30:00.123+05:30",
    "updated_at": "2026-03-04T10:30:00.123+05:30"
  }
}
```

---

## Get Role by ID

**GET** `/roles/:id`

### Request
```
Headers:
Authorization: Bearer <token>

URL Parameter:
- id (required)
```

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "id": 1,
    "role_name": "admin",
    "permissions": ["users.create", "users.read", "users.update", "users.delete"],
    "description": "Administrator role",
    "is_active": true,
    "created_at": "2026-01-15T10:00:00+05:30",
    "updated_at": "2026-01-15T10:00:00+05:30"
  }
}
```

---

## Get All Roles

**GET** `/roles`

### Request
```
Headers:
Authorization: Bearer <token>
```

### Response (200 OK)
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "role_name": "admin",
      "permissions": ["users.create", "users.read", "users.update"],
      "description": "Administrator role",
      "is_active": true,
      "created_at": "2026-01-15T10:00:00+05:30",
      "updated_at": "2026-01-15T10:00:00+05:30"
    },
    {
      "id": 2,
      "role_name": "superadmin",
      "permissions": ["*"],
      "description": "Super Administrator role",
      "is_active": true,
      "created_at": "2026-01-15T10:00:00+05:30",
      "updated_at": "2026-01-15T10:00:00+05:30"
    }
  ]
}
```
