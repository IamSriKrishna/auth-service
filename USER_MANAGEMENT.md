# User Management API

Complete API documentation for creating and retrieving users in the authentication service.

## Table of Contents

- [Create User](#create-user)
- [Get All Users](#get-all-users)
- [Get User by ID](#get-user-by-id)

---

## Create User

### Endpoint

```
POST /auth/admin/create-user
```

### Authentication

**Required**: Super Admin Bearer Token

```
Authorization: Bearer <super-admin-access-token>
```

### Request Body

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

#### Request Fields

| Field | Type | Required | Description | Constraints |
|-------|------|----------|-------------|-------------|
| `name` | string | Yes | User's full name | Min: 1, Max: 255 chars |
| `number` | string | No | User's phone number | Max: 20 chars |
| `email` | string | Yes | User's email address | Must be valid email, unique |
| `company_id` | number | Yes | ID of the company to associate user with | Must reference valid company |
| `password` | string | Yes | User's password | Min: 8 chars, must be secure |
| `user_type` | string | Yes | Type of user | Valid values: `"admin"`, `"partner"` |
| `role_name` | string | Yes | Role to assign to user | Must reference existing active role |

### Response (Success)

**Status Code**: `201 Created`

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
  },
  "meta": {}
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | number | Unique user identifier |
| `email` | string | User's email address |
| `phone` | string | User's phone number |
| `username` | string | User's name |
| `user_type` | string | Type of user (admin, partner, etc.) |
| `role` | string | Assigned role name |
| `status` | string | User status (active, inactive, pending) |
| `company_id` | number | Associated company ID |
| `company` | object | Company details |
| `company.id` | number | Company ID |
| `company.company_name` | string | Company name |
| `company.business_type_id` | number | Business type ID |
| `company.gst_number` | string | GST registration number |
| `company.pan_number` | string | PAN (optional) |
| `company.created_at` | string | Company creation timestamp |
| `company.updated_at` | string | Company last update timestamp |
| `created_at` | string | User creation timestamp |

### Response (Error - User Already Exists)

**Status Code**: `400 Bad Request`

```json
{
  "error": "user already exists with this email"
}
```

### Response (Error - Invalid Role)

**Status Code**: `400 Bad Request`

```json
{
  "error": "invalid role name"
}
```

### Response (Error - Company Not Found)

**Status Code**: `400 Bad Request`

```json
{
  "error": "company not found"
}
```

### Response (Error - Invalid User Type)

**Status Code**: `400 Bad Request`

```json
{
  "error": "invalid user type"
}
```

### cURL Example

```bash
curl -X POST http://localhost:8088/auth/admin/create-user \
  -H "Authorization: Bearer YOUR_SUPERADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "number": "9876543210",
    "email": "john.doe@company.com",
    "company_id": 1,
    "password": "SecurePassword123!",
    "user_type": "admin",
    "role_name": "admin"
  }'
```

---

## Get All Users

### Endpoint

```
GET /auth/admin/users
```

### Authentication

**Required**: Admin or Super Admin Bearer Token

```
Authorization: Bearer <admin-access-token>
```

### Query Parameters

| Parameter | Type | Required | Description | Default | Constraints |
|-----------|------|----------|-------------|---------|-------------|
| `page` | number | No | Page number for pagination | 1 | Min: 1 |
| `limit` | number | No | Items per page | 10 | Min: 1, Max: 100 |
| `search` | string | No | Search by username, email, or phone | - | Max: 255 chars |
| `role` | string | No | Filter users by role name | - | Must be exact role name |

### Request Examples

#### Basic Request (Default Pagination)

```
GET /auth/admin/users
Authorization: Bearer <admin-access-token>
```

#### With Custom Pagination

```
GET /auth/admin/users?page=2&limit=5
Authorization: Bearer <admin-access-token>
```

#### With Search

```
GET /auth/admin/users?search=john&limit=20
Authorization: Bearer <admin-access-token>
```

#### With Role Filter

```
GET /auth/admin/users?role=admin&page=1&limit=10
Authorization: Bearer <admin-access-token>
```

#### Combined Filters

```
GET /auth/admin/users?search=example.com&role=admin&page=1&limit=15
Authorization: Bearer <admin-access-token>
```

### Response (Success)

**Status Code**: `200 OK`

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
      "id": 3,
      "email": "admin@example.com",
      "username": "superadmin",
      "user_type": "superadmin",
      "role": "superadmin",
      "phone": "+1234567890",
      "status": "active",
      "created_at": "2026-01-21T23:25:36.697+05:30",
      "last_login_at": "2026-03-03T23:46:05.38+05:30"
    },
    {
      "id": 7,
      "email": "superadmin@example.com",
      "username": "superadmin1",
      "user_type": "superadmin",
      "role": "superadmin",
      "phone": "1234567890",
      "status": "active",
      "created_at": "2026-02-18T16:00:44.388+05:30",
      "last_login_at": "2026-03-02T23:01:08.084+05:30"
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
    "total": 4,
    "total_pages": 1
  }
}
```

### Response Fields

#### User List Item

| Field | Type | Description |
|-------|------|-------------|
| `id` | number | Unique user identifier |
| `email` | string | User's email address |
| `username` | string | User's name |
| `user_type` | string | Type of user (superadmin, admin, partner, mobile_user) |
| `role` | string | Assigned role name |
| `phone` | string | User's phone number (optional) |
| `status` | string | User status (active, inactive, pending) |
| `company_name` | string | Associated company name (optional, when user has company) |
| `created_at` | string | User creation timestamp |
| `created_by` | number | ID of user who created this user (optional) |
| `last_login_at` | string | Last login timestamp (optional) |

#### Meta Data

| Field | Type | Description |
|-------|------|-------------|
| `current_page` | number | Current page number |
| `per_page` | number | Items per page |
| `total` | number | Total number of users |
| `total_pages` | number | Total number of pages |

### Response (Empty Results)

**Status Code**: `200 OK`

```json
{
  "success": true,
  "data": [],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 0,
    "total_pages": 0
  }
}
```

### cURL Examples

#### Get All Users

```bash
curl -X GET http://localhost:8088/auth/admin/users \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

#### Get Users with Search

```bash
curl -X GET "http://localhost:8088/auth/admin/users?search=john&limit=20" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

#### Get Users by Role

```bash
curl -X GET "http://localhost:8088/auth/admin/users?role=admin&page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

#### Get Specific Page

```bash
curl -X GET "http://localhost:8088/auth/admin/users?page=2&limit=15" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

## Get User by ID

### Endpoint

```
GET /auth/admin/users/:id
```

### Authentication

**Required**: Admin or Super Admin Bearer Token

```
Authorization: Bearer <admin-access-token>
```

### URL Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | number | Yes | User ID |

### Request Example

```
GET /auth/admin/users/8
Authorization: Bearer <admin-access-token>
```

### Response (Success)

**Status Code**: `200 OK`

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

### Response (Error - User Not Found)

**Status Code**: `404 Not Found`

```json
{
  "error": "user not found"
}
```

### cURL Example

```bash
curl -X GET http://localhost:8088/auth/admin/users/8 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

## Filtering by Role - Detailed Example

When you need to get all users with a specific role, use the `role` query parameter.

### Available Roles

First, retrieve available roles:

```
GET /roles
Authorization: Bearer <superadmin-access-token>
```

Response:

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

### Filter Examples

**Get all admin users:**

```bash
curl -X GET "http://localhost:8088/auth/admin/users?role=admin" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Get all superadmin users with pagination:**

```bash
curl -X GET "http://localhost:8088/auth/admin/users?role=superadmin&page=1&limit=20" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Get admin users with specific search:**

```bash
curl -X GET "http://localhost:8088/auth/admin/users?role=admin&search=john" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

## Common Workflows

### 1. Create Multiple Users for a Company

**Step 1**: Get the company ID

```bash
curl -X GET http://localhost:8088/companies \
  -H "Authorization: Bearer YOUR_SUPERADMIN_TOKEN"
```

**Step 2**: Create first user

```bash
curl -X POST http://localhost:8088/auth/admin/create-user \
  -H "Authorization: Bearer YOUR_SUPERADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "number": "9876543210",
    "email": "john@company.com",
    "company_id": 1,
    "password": "SecurePass123!",
    "user_type": "admin",
    "role_name": "admin"
  }'
```

**Step 3**: Create second user

```bash
curl -X POST http://localhost:8088/auth/admin/create-user \
  -H "Authorization: Bearer YOUR_SUPERADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "number": "9876543211",
    "email": "jane@company.com",
    "company_id": 1,
    "password": "SecurePass456!",
    "user_type": "partner",
    "role_name": "partner"
  }'
```

### 2. Search and Filter Users

**Get all admin users from company:**

```bash
curl -X GET "http://localhost:8088/auth/admin/users?role=admin&search=company.com&limit=25" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Get only active superadmin users:**

```bash
curl -X GET "http://localhost:8088/auth/admin/users?role=superadmin&search=&page=1&limit=50" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 3. Pagination Through Large User Lists

**Get first 20 users:**

```bash
curl -X GET "http://localhost:8088/auth/admin/users?page=1&limit=20" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Get next 20 users:**

```bash
curl -X GET "http://localhost:8088/auth/admin/users?page=2&limit=20" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

## Error Handling Guide

| Error | Status Code | Cause | Solution |
|-------|-------------|-------|----------|
| Bad Request | 400 | Invalid request body or parameters | Check request format and required fields |
| User already exists | 400 | Email is already registered | Use unique email address |
| Invalid role name | 400 | Role doesn't exist or is inactive | Check role name against available roles |
| Company not found | 400 | Company ID doesn't reference valid company | Create company first or use valid company ID |
| Invalid user type | 400 | user_type not "admin" or "partner" | Use valid user_type value |
| Unauthorized | 401 | Missing or invalid authorization token | Provide valid bearer token |
| Forbidden | 403 | User doesn't have permission | Use appropriate admin level token |
| Not Found | 404 | User ID doesn't exist | Verify user ID and try again |
| Internal Server Error | 500 | Server error | Check server logs and retry |

---

## Notes

- All timestamps are in ISO 8601 format with timezone information
- Passwords must be at least 8 characters long and contain mixed case and special characters
- Email addresses are case-insensitive but stored as provided
- User status values: `active`, `inactive`, `pending`
- User type values: `mobile_user`, `superadmin`, `admin`, `partner`
- Pagination defaults: page 1, limit 10, maximum limit 100
- Company association is required when creating admin/partner users
- Only super admin can create users with super admin type
- Company name in user list response only shows if company exists in database
